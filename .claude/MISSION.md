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
