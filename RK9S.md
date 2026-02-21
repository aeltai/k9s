# rk9s – SUSE/Rancher Kubernetes TUI

> NOTE: This guide is now also merged into the root project README under **"rk9s Fork Guide (Merged in project root)"**.

rk9s is a **SUSE/Rancher-focused** fork of upstream k9s. It adds plugins and aliases for Rancher, Fleet, RKE2, K3s, Longhorn, Harvester, Kubewarden, and Traefik/NGINX Ingress. Multi-context operations run in **parallel** (inspired by [kubectl-mc](https://github.com/jonnylangefeld/kubectl-mc)).

- Upstream repository: [derailed/k9s](https://github.com/derailed/k9s)
- Upstream docs (core behavior): [https://k9scli.io](https://k9scli.io)
- This document covers rk9s fork-specific behavior and implementation details.

**Press `?` in rk9s to open the Help view.** The **RK9S** section in the legend lists all our shortcuts. They also appear in **RESOURCE** when viewing the matching resource (e.g. nodes, volumes).

---

## Implementation details: upstream baseline + rk9s logic

rk9s keeps upstream k9s core behavior (resource discovery, watchers, rendering, command mode) and layers Rancher/SUSE-focused defaults on top.

### Startup sync and config layout

At startup, rk9s initializes XDG paths and syncs fork defaults:

1. Embedded plugins in `internal/config/default_plugins/*.yaml` are synced to `~/.local/share/rk9s/plugins/`.
2. Dashboard/navigation hotkeys are synced into `~/.config/rk9s/hotkeys.yaml`.
   - `rk9s-*` hotkey entries are force-updated each startup.
   - Non-rk9s user hotkeys are preserved.

### Action resolution order (key handling)

In a resource view, rk9s resolves actions in this order:

1. Base table/view actions.
2. Plugin shortcuts loaded from plugin files.
3. Hotkeys loaded from hotkeys config.
4. View-specific bindings are re-applied last so core per-view actions still win when needed.

Collision behavior:

- A plugin or hotkey must set `override: true` to replace an existing shortcut.
- Without `override: true`, duplicate bindings are ignored and logged.

### Plugin loading precedence

Plugins are merged from multiple locations. If two plugins have the same key/name, the one loaded later wins.

Load order:

1. `~/.config/rk9s/plugins.yaml`
2. Active cluster/context plugin file (`.../clusters/<cluster>/<context>/plugins.yaml`)
3. XDG plugin directories:
   - `$XDG_DATA_DIRS/rk9s/plugins`
   - `$XDG_DATA_HOME/rk9s/plugins`
   - `$XDG_CONFIG_HOME/rk9s/plugins`

In practice, files in `~/.config/rk9s/plugins/` can override bundled defaults that are synced to `~/.local/share/rk9s/plugins/`.

### Plugin execution behavior

- If a plugin has **no** `args` and **no** `pipes`, rk9s treats `command` as an in-TUI command and navigates directly.
- If `inView: true`, command output is rendered inside a Details panel.
- Otherwise, rk9s executes the command as a shell process (foreground/background based on plugin settings).

### Multi-context execution model

Selected contexts are persisted in `~/.config/rk9s/selected_contexts` (or equivalent XDG path), one context per line.

When 2+ contexts are selected and you open a Kubernetes resource view:

- rk9s enables multi-context table mode automatically.
- Data is listed in parallel (up to 10 contexts concurrently) with per-context dynamic clients.
- Unreachable contexts are skipped (warning logged), not fatal.
- A `CLUSTER` column is injected into the table.
- Internal row IDs are encoded as `<context>@@<resource-path>`.

For row actions (`describe`, `yaml`, `edit`, delete, plugin commands), rk9s detects row context and executes in that context.

Plugin env context values:

- `$CONTEXT` = context of the selected row
- `$CONTEXTS` = comma-separated list of all selected contexts

**Example:** On a node, press **Shift-C** to run `kubectl debug node/$NAME ...` in the selected context and show RKE2/K3s config.

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

### Global navigation hotkeys (synced on startup)

| Key | Command | Opens |
|-----|---------|-------|
| **F1** | `home` | Home dashboard |
| **F2** | `clusters.management.cattle.io` | Rancher clusters |
| **F3** | `helmcharts.helm.cattle.io` | RKE2/K3s distro view |
| **F4** | `nodes node-role.kubernetes.io/control-plane=true` | etcd/control-plane nodes |
| **F5** | `nodes` | Nodes |
| **F6** | `gitrepos.fleet.cattle.io` | Fleet GitRepos |
| **F7** | `volumes.longhorn.io` | Longhorn volumes |
| **F8** | `virtualmachines.kubevirt.io` | KubeVirt/Harvester VMs |
| **F9** | `rk9s` | rk9s status dashboard |
| **F10** | `context` | Contexts view |

### Multi-context selection (contexts view)

The context list shows a **SELECTED** column with `+` for each selected context.

| Shortcut | Action |
|----------|--------|
| **Space** | Toggle current context selection |
| **Ctrl-A** | Select all contexts |
| **Ctrl-Space** | Clear selection |
| **Enter** | Switch active context (normal single-context switch) |

After selecting at least 2 contexts, open any Kubernetes resource view (`:nodes`, `:pods`, `:volumes.longhorn.io`, etc.). rk9s will automatically switch that view into multi-context mode, inject a `CLUSTER` column, and fetch data in parallel. `$CONTEXTS` is then available to plugins for multi-cluster actions.

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

1. Open contexts (`:contexts` or **F10**).
2. The **SELECTED** column shows `+` for selected contexts.
3. **Space** to select contexts, **Ctrl-A** to select all.
4. Open any resource view (for example **F5** / `:nodes` or `:pods`).
5. With 2+ contexts selected, the table auto-switches to multi-context mode and shows the `CLUSTER` column.
6. Use normal row actions (`y`, `d`, `e`, plugins). rk9s executes each row action in that row's context.
7. Type `:rk9s` (or `:status`) for a status dashboard with selected contexts and CLI/tooling checks.

Plugins on other views also receive `$CONTEXTS` (comma-separated) for multi-cluster operations.

> **Implementation note:** Multi-context listing is parallel with a max concurrency of 10 contexts and skips unreachable contexts instead of failing the full view.

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
    description: Show selected resource identity
    scopes: [pods]
    command: echo
    args: ["$CONTEXT/$NAMESPACE/$NAME"]

  my-multi-context-plugin:
    shortCut: Shift-3
    description: Check selected pod across all selected contexts
    scopes: [pods]
    command: bash
    inView: true
    args:
      - -c
      - |
        IFS=',' read -r -a CTX_ARR <<< "$CONTEXTS"
        for C in "${CTX_ARR[@]}"; do
          echo "=== [$C] ==="
          kubectl --context "$C" -n "$NAMESPACE" get pod "$NAME" -o wide 2>/dev/null || echo "not found"
        done
```

Use `$CONTEXT` for current row context and `$CONTEXTS` for the selected context set.

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

See [README.md](README.md) for project-level notes and [k9scli.io](https://k9scli.io) for upstream k9s behavior.
