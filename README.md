# dot-agent

Sync AI coding assistant configuration across machines using git.

Supports **Claude Code** (`~/.claude`), **Cursor** (`~/.cursor`), and **[Hermes Agent](https://github.com/NousResearch/hermes-agent)** (`~/.hermes`) on macOS, Linux, and Windows.

## Install

```bash
go install github.com/cthulhu/dot-agent/cmd/dot-agent@latest
```

Or build from source:

```bash
git clone https://github.com/cthulhu/dot-agent.git
cd dot-agent
go build -o dot-agent ./cmd/dot-agent
```

Requires **git** in your `PATH`.

## Quick start

Create a private git repo (GitHub, GitLab, etc.), then on your primary machine:

```bash
dot-agent init --repo git@github.com:you/ai-config.git
dot-agent add claude
dot-agent add cursor
dot-agent add hermes
dot-agent push
```

On a new machine:

```bash
dot-agent init --repo git@github.com:you/ai-config.git
dot-agent pull --apply
```

## Commands

| Command | Description |
|---------|-------------|
| `init [--repo URL] [--path DIR]` | Create or clone the source git repo |
| `add [claude\|cursor\|hermes]` | Capture local config into the repo |
| `apply [claude\|cursor\|hermes]` | Write repo config to local directories |
| `diff [claude\|cursor\|hermes]` | Show differences (source vs local) |
| `status` | Git status + config drift |
| `push` | Commit (if needed) and push |
| `pull [--apply]` | Pull remote; optionally apply locally |
| `doctor` | Validate git, paths, and manifest |

### Flags

- `add --dry-run` / `apply --dry-run` — preview changes
- `apply --backup` — backup overwritten local files
- `apply --force` — apply even when local files differ
- `--source DIR` — override source repo path (default stored in user config)

## Source repo layout

```
dot-agent.yaml
assistants/
  claude/
  cursor/
  hermes/
```

Hermes syncs portable config (`config.yaml`, `SOUL.md`, `memories/`, `skills/`, `cron/`) and skips secrets (`.env`, `auth.json`), sessions, logs, and the installed source tree (`hermes-agent/`). On native Windows, Hermes may use `%LOCALAPPDATA%\hermes` instead of `~/.hermes` — override `target` in `dot-agent.yaml` if needed.

Default paths:

| OS | Source repo | User config |
|----|-------------|-------------|
| macOS / Linux | `~/.local/share/dot-agent/source` | `~/.config/dot-agent/config.yaml` |
| Windows | `%LOCALAPPDATA%\dot-agent\source` | `%APPDATA%\dot-agent\config.yaml` |

## What gets synced

Sensitive and machine-local paths are ignored by default (caches, projects, extensions, credentials). Edit `dot-agent.yaml` in your source repo to customize ignore patterns.

## License

MIT
