# LLM Red Team — technical reference

> **Start here:** the main project README is at the [repository root](../README.md) (install, badges, quick start, configs).

## Layout

| Path | Role |
|------|------|
| `cmd/redteam` | CLI entry (`run`, `list`, `version`) |
| `cmd/mocktarget` | Local OpenAI-compatible mock server for demos |
| `pkg/attacks` | Embedded YAML library + disk loader |
| `pkg/detectors` | PII, toxicity blocklist, topics, leakage, jailbreak heuristics |
| `pkg/runner` | HTTP invoker + concurrency |
| `pkg/report` | JSON + HTML |
| `examples/` | `redteam.yaml`, Anthropic, Ollama, CI mock, custom packs |
| `website/` | Static landing + sample report for buyers |
| `scripts/run-local.sh` | Build mock server + scan in one go |

## Detectors (v1)

Regex / heuristic only — tune for your product; expect false positives. Replace blocklists with classifiers (e.g. Perspective API) for production-grade signal.

## Module path

```text
github.com/Ali627miya/llm-redteam/redteam
```

## Customer / buyer demo

See [DEMO.md](DEMO.md) and [website/](website/).

## Docker

```bash
docker build -t llm-redteam ./redteam   # from repo root
docker run --rm llm-redteam version
```

## TypeScript

[`integrations/typescript`](integrations/typescript) — spawn the CLI from Node (`REDTEAM_BIN` or `PATH`).
