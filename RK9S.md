# rk9s – SUSE/Rancher Kubernetes TUI

rk9s is a **SUSE/Rancher-focused** fork of k9s. It adds plugins and aliases for Rancher, Fleet, RKE2, K3s, Longhorn, Harvester, Kubewarden, and Traefik/NGINX Ingress.

## Quick Install (rk9s + bundled CLIs)

```bash
curl -sSL https://raw.githubusercontent.com/aeltai/k9s/master/scripts/install.sh | bash
```

This installs rk9s plus optional CLIs (longhornctl, kwctl, virtctl, fleet) and preconfigured plugins to `~/.local/bin`.

**Note:** rk9s auto-installs default plugins on first run to `~/.local/share/rk9s/plugins/`. If shortcuts don't work, ensure that directory exists and contains the plugin YAML files.

**Options:**
- `--rk9s-only` – Only rk9s
- `--clis-only` – Only CLIs (longhornctl, kwctl, virtctl, fleet)
- `--plugins-only` – Only copy plugins

**Add to PATH:**
```bash
export PATH="$HOME/.local/bin:$PATH"
```

## Download rk9s

| Method | Command |
|--------|---------|
| **One-liner** | `curl -sSL https://raw.githubusercontent.com/aeltai/k9s/master/scripts/install.sh \| bash` |
| **GitHub Releases** | [Releases](https://github.com/aeltai/k9s/releases) – download `rk9s_<OS>_<arch>.tar.gz` |
| **From source** | `go install github.com/aeltai/k9s@latest` |
| **Build yourself** | `git clone https://github.com/aeltai/k9s && cd k9s && make build` |

**Example (Linux amd64):**
```bash
curl -sL https://github.com/aeltai/k9s/releases/latest/download/rk9s_Linux_amd64.tar.gz | tar xz -C ~/.local/bin
```

## Bundled CLIs (preconfigured)

The installer can install and preconfigure:

| CLI | Purpose |
|-----|---------|
| **longhornctl** | Longhorn volume trim, replica, export |
| **kwctl** | Kubewarden policy inspect, scaffold |
| **virtctl** | Harvester/KubeVirt VM console, VNC |
| **fleet** | Rancher Fleet bundle target, apply |

Plugins in `~/.local/share/rk9s/plugins/` call these automatically.

## Config

- **Config:** `~/.config/rk9s/`
- **State:** `~/.local/state/rk9s/`
- **Plugins:** `~/.config/rk9s/plugins.yaml` or `~/.local/share/rk9s/plugins/*.yaml`

## Run

```bash
rk9s
```

## Legend – rk9s operations

Press **?** in rk9s to see the Help view. Plugin shortcuts appear in the **RESOURCE** section when viewing the matching resource. Uses kubeconfig/API tokens from your context.

### Multi-context selection (contexts view)

| Shortcut | Action |
|----------|--------|
| **Space** | Toggle current context selection |
| **Ctrl-A** | Select all contexts |
| **Ctrl-Space** | Clear selection |
| **Shift-M** | Run across selected contexts (nodes) |

Select one or more contexts, then `$CONTEXTS` is available to plugins (comma-separated) for Rancher-style multi-cluster operations.

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

## API / token usage

rk9s plugins use the standard Kubernetes API and your **kubeconfig** (context, tokens). No extra config needed for:
- **Downstream clusters** – kubectl, longhornctl, virtctl use KUBECONFIG
- **Rancher management API** – use a context pointing at the Rancher server with an API token (see [Rancher API quickstart](https://ranchermanager.docs.rancher.com/api/quickstart))
- **Longhorn** – uses cluster API (Volume, Setting CRs)
- **Harvester** – uses KubeVirt API (VirtualMachine CRs) and virtctl

Optional env vars for advanced usage:
- `KUBECONFIG` – set by rk9s from active context
- `RANCHER_SERVER` – Rancher server URL (for rancher CLI)
- `HARVESTER_CONFIG` – Harvester kubeconfig (for harvester CLI)

See the main [README](README.md) for k9s usage.
