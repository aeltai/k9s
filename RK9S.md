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
| **Harvester / KubeVirt** | virtctl, harvester-cli | kubectl patch VM | KUBECONFIG |
| **Fleet** | fleet | kubectl (GitRepo, Bundle CRs) | KUBECONFIG |
| **Kubewarden** | kwctl | kubectl (Policy CRs) | KUBECONFIG |
| **RKE2/K3s** | kubectl debug | – | KUBECONFIG |

- **rancher CLI** – `rancher login`, `rancher context switch`, `rancher cluster ls`. Needs Rancher API token in `~/.rancher/cli2.json` or kubeconfig.
- **longhornctl** – [Releases](https://github.com/longhorn/cli/releases). Uses `--kubeconfig $KUBECONFIG`.
- **virtctl** – KubeVirt CLI for full VM lifecycle: console, VNC, start, stop, restart, pause, unpause, SSH, migrate, guest agent queries. Same kubeconfig as the cluster. Install via `https://kubevirt.io/user-guide/user_workloads/virtctl_client_tool/`.
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

On startup, rk9s displays a flash banner showing which CLIs are detected (rancher, virtctl, longhornctl, kwctl, fleet) and how many contexts are selected.

Type `:rk9s` or `:status` to see a full status page: CLI versions, selected contexts, nodes across those contexts, and installed plugins.

---

## Legend – rk9s operations

Press **?** in rk9s to see the Help view. The **RK9S** section lists our shortcuts. The **RESOURCE** section shows view-specific bindings.

### Multi-context selection (contexts view)

The context list shows a **SELECTED** column with `+` for each selected context.

| Shortcut | Action |
|----------|--------|
| **Space** | Toggle current context selection |
| **Ctrl-A** | Select all contexts |
| **Ctrl-Space** | Clear selection |
| **Shift-M** | Show nodes from all selected contexts |

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

### Harvester / KubeVirt (VMs)

Alias: `:vm` for VirtualMachine, `:vmi` for VirtualMachineInstance.

| Shortcut | Action | CLI |
|----------|--------|-----|
| **Shift-H** | Open Harvester UI | browser |
| **Shift-S** | Serial console | virtctl console |
| **Shift-V** | VNC console | virtctl vnc |
| **Shift-W** | Start VM | virtctl start / kubectl patch |
| **Shift-X** | Stop VM (confirm) | virtctl stop / kubectl patch |
| **Shift-Z** | Restart VM (confirm) | virtctl restart |
| **Shift-P** | Pause VM | virtctl pause |
| **Shift-Q** | Unpause VM | virtctl unpause |
| **m** | Live-migrate VM (confirm) | virtctl migrate |
| **Shift-Y** | SSH into VM | virtctl ssh |
| **Shift-I** | Guest agent info (OS, FS, users) | virtctl guestosinfo/fslist/userlist |

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
2. The **SELECTED** column shows `+` for selected contexts.
3. **Space** to select contexts, **Ctrl-A** to select all.
4. **Shift-M** runs `kubectl get nodes` across all selected contexts and shows combined output.
5. Type `:rk9s` for a full status page with CLI versions, selected contexts, and nodes.

Plugins on other views also receive `$CONTEXTS` (comma-separated) for multi-cluster operations.

### How to: Trim a Longhorn volume

1. Go to Longhorn volumes (`:vol` or find volumes.longhorn.io).
2. Select a volume.
3. **Shift-T** → confirm → runs `longhornctl trim volume $NAME`.

### How to: Open a VM console in Harvester

1. Go to VirtualMachines (`:vm`) or VirtualMachineInstances (`:vmi`).
2. Select a VM.
3. **Shift-S** for serial console or **Shift-V** for VNC.

### How to: Manage VM lifecycle (start/stop/restart)

1. Go to VirtualMachines (`:vm`).
2. Select a VM.
3. **Shift-W** to start, **Shift-X** to stop, **Shift-Z** to restart (all with confirmation).
4. Requires `virtctl`; falls back to `kubectl patch` for start/stop.

### How to: Pause and unpause a VM

1. Go to VirtualMachines (`:vm`) or VMIs (`:vmi`).
2. Select a VM.
3. **Shift-P** to pause (freeze), **Shift-Q** to unpause.

### How to: SSH into a VM

1. Go to VirtualMachines (`:vm`) or VMIs (`:vmi`).
2. Select a VM.
3. **Shift-Y** connects via `virtctl ssh` as `root`.
4. Set `VIRTCTL_SSH_USER=fedora` (or your user) to change the login user.

### How to: Live-migrate a VM

1. Go to VirtualMachines (`:vm`) or VMIs (`:vmi`).
2. Select the VM.
3. **m** triggers `virtctl migrate` (with confirmation).

### How to: Query VM guest agent info

1. Go to VirtualMachines (`:vm`) or VMIs (`:vmi`).
2. Select a VM with qemu-guest-agent installed.
3. **Shift-I** shows OS info, filesystems, and logged-in users.

### How to: Inspect a Kubewarden policy

1. Go to ClusterAdmissionPolicy or AdmissionPolicy.
2. Select a policy.
3. **Shift-I** runs `kwctl inspect` on the policy module.

### How to: Add custom plugins

Create `~/.config/rk9s/plugins.yaml` or add YAML in `~/.local/share/rk9s/plugins/`:

```yaml
plugins:
  my-plugin:
    shortCut: Shift-2
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
