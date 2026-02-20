// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of rk9s

package dao

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/k9s/internal/slogs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

const mcMaxParallel = 10

// ContextObject pairs a runtime.Object with the context it was fetched from.
type ContextObject struct {
	Context string
	Object  runtime.Object
}

var (
	dynClientCache sync.Map // map[string]dynamic.Interface
)

// ResetDynClientCache clears cached per-context dynamic clients.
func ResetDynClientCache() {
	dynClientCache = sync.Map{}
}

func dynClientFor(rawConfig api.Config, ctxName string) (dynamic.Interface, error) {
	if c, ok := dynClientCache.Load(ctxName); ok {
		return c.(dynamic.Interface), nil
	}

	overrides := &clientcmd.ConfigOverrides{CurrentContext: ctxName}
	cc := clientcmd.NewDefaultClientConfig(rawConfig, overrides)
	restCfg, err := cc.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("rest config for context %q: %w", ctxName, err)
	}
	restCfg.QPS = 50
	restCfg.Burst = 100

	dc, err := dynamic.NewForConfig(restCfg)
	if err != nil {
		return nil, fmt.Errorf("dynamic client for context %q: %w", ctxName, err)
	}
	dynClientCache.Store(ctxName, dc)

	return dc, nil
}

// MultiContextList fetches resources across multiple contexts in parallel
// using Go dynamic clients. Unreachable contexts are logged and skipped.
func MultiContextList(
	rawConfig api.Config,
	contexts []string,
	gvr schema.GroupVersionResource,
	ns string,
	labelSel string,
) ([]ContextObject, error) {
	type result struct {
		ctx     string
		objects []*unstructured.Unstructured
		err     error
	}

	ch := make(chan result, len(contexts))
	sem := make(chan struct{}, mcMaxParallel)

	for _, ctxName := range contexts {
		sem <- struct{}{}
		go func(ctx string) {
			defer func() { <-sem }()

			dc, err := dynClientFor(rawConfig, ctx)
			if err != nil {
				ch <- result{ctx: ctx, err: err}
				return
			}

			opts := metav1.ListOptions{}
			if labelSel != "" {
				opts.LabelSelector = labelSel
			}

			var list *unstructured.UnstructuredList
			if ns == "" || ns == client.ClusterScope || ns == client.NamespaceAll {
				list, err = dc.Resource(gvr).List(context.Background(), opts)
			} else {
				list, err = dc.Resource(gvr).Namespace(ns).List(context.Background(), opts)
			}
			if err != nil {
				ch <- result{ctx: ctx, err: err}
				return
			}

			objs := make([]*unstructured.Unstructured, len(list.Items))
			for i := range list.Items {
				objs[i] = &list.Items[i]
			}
			ch <- result{ctx: ctx, objects: objs}
		}(ctxName)
	}

	var out []ContextObject
	for range contexts {
		r := <-ch
		if r.err != nil {
			slog.Warn("Multi-context list skipped context",
				slogs.Subsys, "mc",
				"context", r.ctx,
				slogs.Error, r.err,
			)
			continue
		}
		for _, o := range r.objects {
			out = append(out, ContextObject{Context: r.ctx, Object: o})
		}
	}

	return out, nil
}

// MultiContextServerVersions queries the /version endpoint for each context
// and returns a map of context-name -> K8s version string.
func MultiContextServerVersions(rawConfig api.Config, contexts []string) map[string]string {
	type verResult struct {
		ctx string
		ver string
	}

	ch := make(chan verResult, len(contexts))
	for _, ctxName := range contexts {
		go func(ctx string) {
			overrides := &clientcmd.ConfigOverrides{CurrentContext: ctx}
			cc := clientcmd.NewDefaultClientConfig(rawConfig, overrides)
			restCfg, err := cc.ClientConfig()
			if err != nil {
				ch <- verResult{ctx: ctx, ver: client.NA}
				return
			}
			restCfg.Timeout = 5 * time.Second

			dc, err := discovery.NewDiscoveryClientForConfig(restCfg)
			if err != nil {
				ch <- verResult{ctx: ctx, ver: client.NA}
				return
			}
			info, err := dc.ServerVersion()
			if err != nil {
				ch <- verResult{ctx: ctx, ver: client.NA}
				return
			}
			ch <- verResult{ctx: ctx, ver: info.GitVersion}
		}(ctxName)
	}

	out := make(map[string]string, len(contexts))
	for range contexts {
		r := <-ch
		out[r.ctx] = r.ver
	}
	return out
}
