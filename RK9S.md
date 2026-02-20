# rk9s – SUSE/Rancher Kubernetes TUI

rk9s is a **SUSE/Rancher-focused** fork of k9s. It adds plugins and aliases for Rancher, Fleet, RKE2, K3s, Longhorn, Harvester, Kubewarden, and Traefik/NGINX Ingress.

## Quick Install (rk9s + bundled CLIs)

```bash
curl -sSL https://raw.githubusercontent.com/aeltai/k9s/master/scripts/install.sh | bash
```

This installs rk9s plus optional CLIs (longhornctl, kwctl, virtctl, fleet) and preconfigured plugins to `~/.local/bin`.

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

See the main [README](README.md) for k9s usage.
