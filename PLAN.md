# Plan: Extract jack msg into mtx

Standalone Matrix messaging CLI for agentic coordination. Extracted from `github.com/zoobzio/jack/msg`.

## Source

All source files live at `../jack/msg/`. Copy them — don't move — so jack stays buildable until it cuts over.

## Step 1: Scaffold the module

The repo is currently the samoa template. Transform it:

1. **go.mod**: rename module from `github.com/zoobzio/samoa` to `github.com/zoobzio/mtx`
2. **api.go**: delete it (samoa placeholder)
3. **testing/**: keep the infrastructure, clear the placeholder content
4. **.goreleaser.yml**: update `project_name` to `mtx`, change `builds` from `skip: true` to build a binary from `cmd/mtx/main.go` targeting linux/darwin/amd64/arm64
5. **.claude/MISSION.md**: rewrite for mtx (see below)
6. **go.sum**: will be regenerated

## Step 2: Create the package structure

```
mtx/
├── cmd/mtx/main.go     # CLI entrypoint
├── mtx.go              # Client, types, helpers (from msg.go)
├── board.go            # Board commands + ProvisionGlobalBoard/AnnounceOnBoard
├── register.go         # register subcommand
├── login.go            # login subcommand
├── send.go             # send subcommand
├── read.go             # read subcommand
├── watch.go            # watch subcommand
├── create.go           # create subcommand
├── dm.go               # dm send/read subcommands
├── invite.go           # invite subcommand
├── invites.go          # invites subcommand
├── join.go             # join subcommand
├── leave.go            # leave subcommand
├── members.go          # members subcommand
├── rooms.go            # rooms subcommand
├── whoami.go           # whoami subcommand
├── who.go              # who subcommand (see step 4)
├── *_test.go           # all corresponding test files
├── go.mod
└── cmd/
    └── mtx/
        └── main.go
```

## Step 3: Rename and decouple

These are the mechanical changes needed across every copied file:

### 3a. Package declaration

Every file: `package msg` -> `package mtx`

### 3b. Environment variables

In `mtx.go` (was `msg.go`), rename the env var keys:

| Old | New | Fallback |
|-----|-----|----------|
| `JACK_MSG_TOKEN` | `MTX_TOKEN` | Also check `JACK_MSG_TOKEN` for transition period |
| `JACK_TEAM` | `MTX_TEAM` | Also check `JACK_TEAM` for transition period |

Update `TokenFromEnv()` and `TeamFromEnv()` accordingly. Also update `envFromFile()` to walk up looking for `.mtx/env` in addition to `.jack/env`.

### 3c. Package-level globals -> Config struct

Replace the three package-level globals:

```go
// Before (msg.go)
var (
    Homeserver        string
    RegistrationToken string
    DataDir           string
)
```

With a config struct and loader:

```go
// Config holds mtx configuration.
type Config struct {
    Homeserver        string `yaml:"homeserver"`
    RegistrationToken string `yaml:"registration_token"`
    DataDir           string `yaml:"data_dir"`
}
```

Load from `~/.config/mtx/config.yaml` (respect `MTX_CONFIG_DIR` env var). Every subcommand that currently reads package globals should get config through cobra's PersistentPreRunE on the root command.

### 3d. Root command

Rename the root cobra command:

```go
// Before
var Cmd = &cobra.Command{
    Use:   "msg",
    Short: "Matrix messaging commands",
}

// After
var rootCmd = &cobra.Command{
    Use:   "mtx",
    Short: "Matrix messaging CLI for agentic coordination",
}
```

### 3e. cmd/mtx/main.go

Minimal entrypoint:

```go
package main

import (
    "os"
    "github.com/zoobzio/mtx"
)

func main() {
    if err := mtx.Execute(); err != nil {
        os.Exit(1)
    }
}
```

Add an `Execute()` function in the root package that calls `rootCmd.Execute()`.

## Step 4: Handle the `who` command

`who` is the one command with a jack-specific dependency — it reads jack's `registry.yaml` to list session users. Two options:

**Option A (recommended)**: Keep `who` but make the data dir configurable. If `MTX_DATA_DIR` or config `data_dir` points at jack's data dir, it works. If not configured, `who` prints an error saying no registry is available. This keeps mtx useful as a jack companion without hard-coupling it.

**Option B**: Drop `who` entirely from mtx. It's a jack concern — jack can reimplement it locally.

Go with Option A for now. The `registryData` struct and `loadRegistry()` function stay in `who.go` — they have no jack imports, just YAML parsing.

## Step 5: Add dependencies

```sh
go get github.com/spf13/cobra@latest
go get gopkg.in/yaml.v3@latest
```

These are the only two external deps. The Matrix client is hand-rolled HTTP.

## Step 6: Verify

```sh
make ci          # vet + lint + test
go build ./cmd/mtx/
./mtx --help
```

All existing tests from `jack/msg/*_test.go` should pass with only the package rename.

## Step 7: Config file format

`~/.config/mtx/config.yaml`:

```yaml
homeserver: https://matrix.example.com
registration_token: secret
data_dir: ~/.jack  # optional, for who command compatibility
```

## Step 8: Update .claude/MISSION.md

```markdown
# Mission: mtx

Standalone Matrix messaging CLI for agentic coordination.

## Purpose

Provide a lightweight, framework-agnostic CLI for agent-to-agent messaging
over the Matrix protocol. Extracted from jack's msg subsystem to serve as
an independent tool that any agent orchestrator can use.

## What This Package Contains

- Matrix client library (v3 client API, no external Matrix SDK)
- CLI for rooms, messaging, DMs, presence, and board coordination
- Board system for team and global announcements

## What This Package Does NOT Contain

- Agent orchestration or session management (that's jack)
- Sandboxing or process isolation
- Git/GitHub integration
```

## Files to NOT copy

- Nothing from jack outside of `msg/` is needed
- No jack config types, no age encryption, no tmux code

## Ordering

1. Scaffold (steps 1-2)
2. Copy + rename (step 3)
3. Handle who (step 4)
4. Wire deps (step 5)
5. Verify (step 6)
6. Config + docs (steps 7-8)

All of this is a single PR on the mtx repo.
