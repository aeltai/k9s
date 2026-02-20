// Package mc provides parallel multi-context kubectl execution.
// Inspired by https://github.com/jonnylangefeld/kubectl-mc (MIT).
// Runs kubectl commands across multiple Kubernetes contexts concurrently
// using goroutines with a configurable concurrency limit.
package mc

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

const defaultMaxProc = 10

// Result holds the output from a single context execution.
type Result struct {
	Context string
	Output  string
	Err     error
}

// RunParallel executes kubectl with the given args across all contexts in parallel.
// maxProc limits concurrent goroutines (0 = default 10).
// Returns results in the same order as the input contexts.
func RunParallel(contexts []string, args []string, maxProc int) []Result {
	if maxProc <= 0 {
		maxProc = defaultMaxProc
	}

	results := make([]Result, len(contexts))
	sem := make(chan struct{}, maxProc)
	var wg sync.WaitGroup

	for i, ctx := range contexts {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, context string) {
			defer wg.Done()
			defer func() { <-sem }()

			localArgs := injectContext(args, context)
			cmd := exec.Command("kubectl", localArgs...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			r := Result{Context: context}
			if err != nil {
				errMsg := strings.TrimSpace(stderr.String())
				if errMsg == "" {
					errMsg = err.Error()
				}
				r.Err = fmt.Errorf("%s", errMsg)
				r.Output = errMsg
			} else {
				r.Output = stdout.String()
			}
			results[idx] = r
		}(i, ctx)
	}
	wg.Wait()

	return results
}

// FormatResults produces a human-readable multi-context output like kubectl-mc.
func FormatResults(results []Result) string {
	var b strings.Builder
	for _, r := range results {
		header := r.Context
		b.WriteString(fmt.Sprintf("\n%s\n%s\n", header, strings.Repeat("-", len(header))))
		if r.Err != nil {
			b.WriteString(fmt.Sprintf("  (error) %s\n", r.Output))
		} else {
			b.WriteString(r.Output)
		}
	}
	return b.String()
}

// injectContext adds --context <ctx> to the kubectl args.
// If args contain "--", inject before it.
func injectContext(args []string, context string) []string {
	local := make([]string, 0, len(args)+2)
	injected := false
	for _, a := range args {
		if a == "--" && !injected {
			local = append(local, "--context", context)
			injected = true
		}
		local = append(local, a)
	}
	if !injected {
		local = append(local, "--context", context)
	}
	return local
}
