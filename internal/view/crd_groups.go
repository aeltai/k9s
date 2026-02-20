// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package view

import (
	"fmt"
	"strings"
)

// crdGroups defines related CRD "tab groups" for ecosystem navigation.
// Use ← / → arrows to cycle between related CRDs within the same group.
// Entries use the alias format "resource.group" for CRDs, or "v1/resource" for core K8s resources.
var crdGroups = [][]string{
	// F7: Longhorn
	{
		"volumes.longhorn.io",
		"replicas.longhorn.io",
		"engines.longhorn.io",
		"nodes.longhorn.io",
		"backupvolumes.longhorn.io",
	},
	// F6: Fleet
	{
		"gitrepos.fleet.cattle.io",
		"bundledeployments.fleet.cattle.io",
		"bundles.fleet.cattle.io",
		"clustergroups.fleet.cattle.io",
		"clusters.fleet.cattle.io",
	},
	// F2: Rancher
	{
		"clusters.management.cattle.io",
		"projects.management.cattle.io",
		"users.management.cattle.io",
		"settings.management.cattle.io",
		"clusterrepos.catalog.cattle.io",
	},
	// F8: KubeVirt (CDI datavolumes skipped — requires separate CDI operator install)
	{
		"virtualmachines.kubevirt.io",
		"virtualmachineinstances.kubevirt.io",
	},
	// F3: Distro (RKE2/K3s) — HelmCharts → HelmChartConfigs → UpgradePlans → K3s Addons
	{
		"helmcharts.helm.cattle.io",
		"helmchartconfigs.helm.cattle.io",
		"plans.upgrade.cattle.io",
		"addons.k3s.cattle.io",
	},
	// F4: etcd — control-planes first (etcd health visible immediately), snapshots second
	{
		"v1/nodes|node-role.kubernetes.io/control-plane",
		"etcdsnapshots.rke.cattle.io",
	},
	// F5: Nodes ecosystem — Nodes → NodePools → Machines → MachineDeployments
	{
		"v1/nodes",
		"nodepools.management.cattle.io",
		"machines.cluster.x-k8s.io",
		"machinedeployments.cluster.x-k8s.io",
	},
	// Kubewarden
	{
		"clusteradmissionpolicies.policies.kubewarden.io",
		"admissionpolicies.policies.kubewarden.io",
		"policyservers.policies.kubewarden.io",
	},
}

// crdGroupIndex maps each CRD alias key to its group and position for O(1) lookup.
var crdGroupIndex map[string]struct {
	group int
	pos   int
}

func init() {
	crdGroupIndex = make(map[string]struct {
		group int
		pos   int
	})
	for gi, grp := range crdGroups {
		for pi, crd := range grp {
			crdGroupIndex[crd] = struct {
				group int
				pos   int
			}{gi, pi}
		}
	}
}

// parseCRDEntry splits a crdGroups entry that may have an optional label selector
// appended with "|", e.g. "v1/nodes|node-role.kubernetes.io/control-plane".
// Returns (navCmd, labelSelector, hasFilter).
func parseCRDEntry(entry string) (string, string, bool) {
	if idx := strings.Index(entry, "|"); idx >= 0 {
		return entry[:idx], entry[idx+1:], true
	}
	return entry, "", false
}

// gvrToAliasKey converts a full GVR path (group/version/resource) to the
// alias key used in crdGroupIndex.
//
//   - Exact match wins (e.g. an already-alias or a "v1/nodes|label" key)
//   - group/version/resource → resource.group  (CRD path → alias name)
//   - v1/resource             → look for "v1/resource|..." prefix entry
//
// Example: "fleet.cattle.io/v1alpha1/gitrepos" → "gitrepos.fleet.cattle.io"
// Example: "v1/nodes" (when etcd group has "v1/nodes|...") → "v1/nodes|..."
func gvrToAliasKey(gvrStr string) string {
	if _, ok := crdGroupIndex[gvrStr]; ok {
		return gvrStr
	}
	parts := strings.Split(gvrStr, "/")
	if len(parts) == 3 {
		alias := parts[2] + "." + parts[0]
		if _, ok := crdGroupIndex[alias]; ok {
			return alias
		}
		return alias
	}
	// For "v1/resource" core resources, check for a "v1/resource|label" entry.
	if len(parts) == 2 {
		prefix := gvrStr + "|"
		for k := range crdGroupIndex {
			if strings.HasPrefix(k, prefix) {
				return k
			}
		}
	}
	return gvrStr
}

// crdDisplayNames maps alias keys to short human-readable tab labels.
var crdDisplayNames = map[string]string{
	// Longhorn
	"volumes.longhorn.io":      "Volumes",
	"replicas.longhorn.io":     "Replicas",
	"engines.longhorn.io":      "Engines",
	"nodes.longhorn.io":        "LH-Nodes",
	"backupvolumes.longhorn.io": "Backups",
	// Fleet
	"gitrepos.fleet.cattle.io":           "GitRepos",
	"bundledeployments.fleet.cattle.io":  "BundleDeploys",
	"bundles.fleet.cattle.io":            "Bundles",
	"clustergroups.fleet.cattle.io":      "ClusterGroups",
	"clusters.fleet.cattle.io":           "Clusters",
	// Rancher
	"clusters.management.cattle.io":   "Clusters",
	"projects.management.cattle.io":   "Projects",
	"users.management.cattle.io":      "Users",
	"settings.management.cattle.io":   "Settings",
	"clusterrepos.catalog.cattle.io":  "Repos",
	// KubeVirt
	"virtualmachines.kubevirt.io":         "VMs",
	"virtualmachineinstances.kubevirt.io": "VMIs",
	// Distro / F3
	"helmcharts.helm.cattle.io":       "HelmCharts",
	"helmchartconfigs.helm.cattle.io": "HelmConfigs",
	"plans.upgrade.cattle.io":         "UpgradePlans",
	"addons.k3s.cattle.io":            "Addons",
	// etcd / F4
	"etcdsnapshots.rke.cattle.io":                    "Snapshots",
	"v1/nodes|node-role.kubernetes.io/control-plane": "ControlPlanes",
	// Nodes ecosystem / F5
	"v1/nodes":                             "Nodes",
	"nodepools.management.cattle.io":       "NodePools",
	"machines.cluster.x-k8s.io":           "Machines",
	"machinedeployments.cluster.x-k8s.io": "MachineDeployments",
	// Kubewarden
	"clusteradmissionpolicies.policies.kubewarden.io": "ClusterPolicies",
	"admissionpolicies.policies.kubewarden.io":        "Policies",
	"policyservers.policies.kubewarden.io":            "PolicyServers",
}

// crdTabHint builds a coloured tab bar string for the table title.
// The current CRD is highlighted in green/bold; others are dimmed.
// Returns "" when the GVR is not part of any group.
func crdTabHint(gvrStr string) string {
	aliasKey := gvrToAliasKey(gvrStr)
	entry, ok := crdGroupIndex[aliasKey]
	if !ok {
		return ""
	}
	grp := crdGroups[entry.group]
	if len(grp) < 2 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("  [gray::-]‹[-]")
	for i, crd := range grp {
		// Use the full entry (with optional |label) as display name key first,
		// then fall back to just the nav command part.
		label := crd
		if l, ok := crdDisplayNames[crd]; ok {
			label = l
		} else if navCmd, _, hasFilter := parseCRDEntry(crd); hasFilter {
			if l, ok := crdDisplayNames[navCmd]; ok {
				label = l
			}
		}
		if i == entry.pos {
			sb.WriteString(fmt.Sprintf("[green::b] %s [-]", label))
		} else {
			sb.WriteString(fmt.Sprintf("[gray::-] %s [-]", label))
		}
		if i < len(grp)-1 {
			sb.WriteString("[gray::-]·[-]")
		}
	}
	sb.WriteString("[gray::-]›[-]")
	return sb.String()
}

// nextCRDInGroup returns the next CRD entry in the group (navCmd, labelSelector, found).
func nextCRDInGroup(aliasKey string) (string, string, bool) {
	entry, ok := crdGroupIndex[aliasKey]
	if !ok {
		return "", "", false
	}
	grp := crdGroups[entry.group]
	if len(grp) < 2 {
		return "", "", false
	}
	next := (entry.pos + 1) % len(grp)
	navCmd, label, _ := parseCRDEntry(grp[next])
	return navCmd, label, true
}

// prevCRDInGroup returns the previous CRD entry in the group (navCmd, labelSelector, found).
func prevCRDInGroup(aliasKey string) (string, string, bool) {
	entry, ok := crdGroupIndex[aliasKey]
	if !ok {
		return "", "", false
	}
	grp := crdGroups[entry.group]
	if len(grp) < 2 {
		return "", "", false
	}
	prev := (entry.pos - 1 + len(grp)) % len(grp)
	navCmd, label, _ := parseCRDEntry(grp[prev])
	return navCmd, label, true
}
