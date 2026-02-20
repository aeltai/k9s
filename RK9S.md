# rk9s – SUSE/Rancher Kubernetes TUI

rk9s is a **SUSE/Rancher-focused** fork of k9s. It adds plugins and aliases for Rancher, Fleet, RKE2, K3s, Longhorn, Harvester, Kubewarden, and Traefik/NGINX Ingress.

**Press `?` in rk9s to open the Help view.** The **RK9S** section in the legend lists all our shortcuts. They also appear in **RESOURCE** when viewing the matching resource (e.g. nodes, volumes).

---

## Introduction: How rk9s works

rk9s extends k9s with **plugins** – key bindings that run shell commands. Each plugin is scoped to a resource type (e.g. nodes, pods, Longhorn volumes). When you press a shortcut in the right view, rk9s passes context via env vars (`$NAME`, `$NAMESPACE`, `$CONTEXT`, `$KUBECONFIG`) and runs the command.

**Example:** On a node, press **Shift-C** → rk9s runs `kubectl debug node/$NAME ...` to show RKE2/K3s config.

---

## How CLIs are integrated

rk9s calls external CLIs when available; otherwise it falls back to kubectl or shows a hint.

| Integration | Primary CLI | Fallback | Auth |
|-------------|-------------|----------|------|
| **Rancher** | rancher CLI | kubectl (projects, clusters CRs) | kubeconfig + API token |
| **Longhorn** | longhornctl | kubectl (Volume CRs) | KUBECONFIG |
| **Harvester** | virtctl, harvester-cli | kubectl | KUBECONFIG |
| **Fleet** | fleet | kubectl (GitRepo, Bundle CRs) | KUBECONFIG |
| **Kubewarden** | kwctl | kubectl (Policy CRs) | KUBECONFIG |
| **RKE2/K3s** | kubectl debug | – | KUBECONFIG |

- **rancher CLI** – `rancher login`, `rancher context switch`, `rancher cluster ls`. Needs Rancher API token in `~/.rancher/cli2.json` or kubeconfig.
- **longhornctl** – [Releases](https://github.com/longhorn/cli/releases). Uses `--kubeconfig $KUBECONFIG`.
- **virtctl** – KubeVirt CLI for VM console/VNC. Same kubeconfig as the cluster.
- **kwctl** – Kubewarden policy tool. Inspects policies from registry.
- **fleet** – Rancher Fleet CLI. Uses kubeconfig for GitOps operations.

Plugins are in `~/.local/share/rk9s/plugins/` (auto-installed on first run) or `~/.config/rk9s/plugins.yaml`.

---

## Quick Install (rk9s + bundled CLIs)

```bash
curl -sSL https://raw.githubusercontent.com/aeltai/k9s/master/scripts/install.sh | bash
```

This installs rk9s plus optional CLIs (longhornctl, kwctl, virtctl, fleet) and preconfigured plugins to `~/.local/bin`.

**Options:**
- `--rk9s-only` – Only rk9s
- `--clis-only` – Only CLIs
- `--plugins-only` – Only copy plugins

**Add to PATH:**
```bash
export PATH="$HOME/.local/bin:$PATH"
```

---

## Download rk9s

| Method | Command |
|--------|---------|
| **One-liner** | `curl -sSL https://raw.githubusercontent.com/aeltai/k9s/master/scripts/install.sh \| bash` |
| **GitHub Releases** | [Releases](https://github.com/aeltai/k9s/releases) – download `rk9s_<OS>_<arch>.tar.gz` |
| **From source** | `go install github.com/aeltai/k9s@latest` |
| **Build yourself** | `git clone https://github.com/aeltai/k9s && cd k9s && make build` |

---

## Config

- **Config:** `~/.config/rk9s/`
- **State:** `~/.local/state/rk9s/`
- **Plugins:** `~/.config/rk9s/plugins.yaml` or `~/.local/share/rk9s/plugins/*.yaml`

---

## Run

```bash
rk9s
```

---

## Legend – rk9s operations

Press **?** in rk9s to see the Help view. The **RK9S** section lists our shortcuts. The **RESOURCE** section shows view-specific bindings.

### Multi-context selection (contexts view)

| Shortcut | Action |
|----------|--------|
| **Space** | Toggle current context selection |
| **Ctrl-A** | Select all contexts |
| **Ctrl-Space** | Clear selection |
| **Shift-M** | Run across selected contexts (nodes) |

Select one or more contexts; `$CONTEXTS` is then available to plugins (comma-separated) for multi-cluster operations.

### All views
| Shortcut | Action |
|----------|--------|
| **Shift-O** | Open Rancher dashboard |
| **Shift-F** | Open Fleet UI |
| **Shift-K** | Sync kubeconfig (rancher context switch) |
| **Shift-J** | List Rancher projects (management API) |
| **Shift-U** | List Rancher clusters (rancher CLI / kubectl) |
| **Shift-L** | Open Longhorn UI (port-forward + browser) |
| **Shift-H** | Open Harvester UI |

### Nodes (RKE2/K3s)
| Shortcut | Action | CLI |
|----------|--------|-----|
| **Shift-C** | View RKE2 or K3s config | kubectl debug |
| **Shift-D** | Node services (systemctl) | kubectl debug |
| **Shift-P** | crictl ps on node | kubectl debug |
| **Shift-E** | etcdctl endpoint health | kubectl debug |
| **Shift-N** | Longhorn node info | longhornctl / kubectl |

### Fleet
| Shortcut | Action | CLI |
|----------|--------|-----|
| **Shift-G** | GitRepo status | kubectl |
| **Shift-R** | Force reconcile GitRepo | kubectl annotate |
| **Shift-T** | Bundle target (which clusters) | fleet / kubectl |

### Longhorn
| Shortcut | Action | CLI |
|----------|--------|-----|
| **Shift-R** | Volume replica info | kubectl |
| **Shift-S** | Volume YAML | kubectl |
| **Shift-T** | Trim volume | longhornctl |
| **Shift-B** | Volume snapshots | longhornctl / kubectl |
| **Shift-N** | Longhorn node info | longhornctl / kubectl |

### Harvester / KubeVirt
| Shortcut | Action | CLI |
|----------|--------|-----|
| **Shift-S** | VM shell | virtctl / harvester |
| **Shift-V** | VM VNC console | virtctl |

### Kubewarden
| Shortcut | Action | CLI |
|----------|--------|-----|
| **Shift-K** | List policies | kwctl |
| **Shift-I** | Inspect policy | kwctl |
| **Shift-S** | Scaffold manifest | kwctl |

### Ingress (Traefik / NGINX)
| Shortcut | Action |
|----------|--------|
| **Shift-I** | Describe IngressRoute / VirtualServer |

---

## How-tos

### How to: Use Rancher management cluster

1. Add a context pointing at your Rancher server (API token in kubeconfig).
2. Switch to that context (`:contexts` → Enter on the context).
3. **Shift-O** opens the Rancher UI, **Shift-J** lists projects, **Shift-U** lists clusters.

### How to: Run commands across multiple clusters

1. `:contexts` to open the context list.
2. **Space** to select contexts, **Ctrl-A** to select all.
3. **Shift-M** runs `kubectl get nodes` on each selected context.

### How to: Trim a Longhorn volume

1. Go to Longhorn volumes (`:vol` or find volumes.longhorn.io).
2. Select a volume.
3. **Shift-T** → confirm → runs `longhornctl trim volume $NAME`.

### How to: Open a VM console in Harvester

1. Go to VirtualMachines (`:vm` or virtualmachines.kubevirt.io).
2. Select a VM.
3. **Shift-S** for shell (virtctl console) or **Shift-V** for VNC.

### How to: Inspect a Kubewarden policy

1. Go to ClusterAdmissionPolicy or AdmissionPolicy.
2. Select a policy.
3. **Shift-I** runs `kwctl inspect` on the policy module.

### How to: Add custom plugins

Create `~/.config/rk9s/plugins.yaml` or add YAML in `~/.local/share/rk9s/plugins/`:

```yaml
plugins:
  my-plugin:
    shortCut: Shift-X
    description: My custom action
    scopes: [pods]
    command: echo
    args: [$NAME $NAMESPACE]
```

---

## API / token usage

rk9s uses your **kubeconfig** and active context. No extra config for:

- **Downstream clusters** – kubectl, longhornctl, virtctl use KUBECONFIG
- **Rancher management API** – use a context pointing at the Rancher server with an API token ([Rancher API quickstart](https://ranchermanager.docs.rancher.com/api/quickstart))
- **Longhorn** – cluster API (Volume, Setting CRs)
- **Harvester** – KubeVirt API (VirtualMachine CRs) and virtctl

Optional env vars:
- `KUBECONFIG` – set by rk9s
- `RANCHER_SERVER` – Rancher server URL (rancher CLI)
- `HARVESTER_CONFIG` – Harvester kubeconfig (harvester CLI)

See the main [README](README.md) for k9s usage.
