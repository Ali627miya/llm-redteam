# LLM Red Team

## Customer / buyer demo

- **Marketing site (static):** `website/` — deploy `website` to [Vercel](https://vercel.com), [Netlify](https://netlify.com), or GitHub Pages. See **`DEMO.md`** for step-by-step hosting and a 5-minute pitch script.
- **Sample report:** open `website/sample-report.html` in a browser (no install).
- **Local preview:** `docker compose -f docker-compose.website.yml up` → http://localhost:8080

---

Developer-first CLI for **automated LLM security probes**: jailbreaks, prompt injection, toxic output probes, coarse PII / secret leakage heuristics, and disallowed-topic patterns. Runs **locally or in CI** against any HTTP endpoint you configure (OpenAI-compatible APIs, your own gateway, or a staging app).

This repository is the **open-source core**. A natural commercial extension is a hosted control plane (scheduled runs, history, dashboards, SSO, org-wide policies) while keeping the engine free.

## Quick start

Requirements: **Go 1.22+**

```bash
cd redteam
make build
./bin/redteam version
./bin/redteam list
./bin/redteam run --config examples/redteam-mock.yaml --output report.html --format html --mock --fail-on-findings
```

- `list` / `run` — print built-in attack IDs or execute the library; `--mock` uses a deterministic safe stub (no network) for CI.
- `--fail-on-findings` — exit code `1` if any detector flags output (use in CI when you expect zero issues).
- `-attacks-dir` (repeatable) — merge extra YAML packs from disk; `-no-builtin` uses only those packs.
- `version` or `redteam -version` — print release string (overridable with `-ldflags` at build time).

`make run-mock` runs the same with the Makefile.

### Local HTTP mock target

For realistic JSON request/response testing without OpenAI:

```bash
# terminal A
make run-mockserver-vulnerable

# terminal B — expect findings (vulnerable persona triggers detectors)
make run-against-mockserver
```

`cmd/mocktarget` implements `POST /v1/chat/completions` with `-persona safe` or `vulnerable`.

## Configuration

See `examples/redteam.yaml` for a template aimed at OpenAI’s chat completions API. Important fields:

| Field | Purpose |
|--------|---------|
| `target.url` | POST endpoint |
| `target.body_template` | Go template; use `{{toJSON .Prompt}}` and optional `{{toJSON .Context}}` |
| `target.response_path` | Dot path into JSON for model text (e.g. `choices.0.message.content`) |
| `context` | Optional synthetic context for **leakage** detectors (use fake PII in CI, never production secrets) |
| `categories` | Optional list of suite names matching files in `pkg/attacks/library/` (stems): `jailbreak`, `prompt_injection`, `toxicity`, `data_exfil` |
| `detectors` | Per-detector toggles (`pii`, `toxicity`, …); omitted keys default to **on** |

Headers and URLs expand **`${VAR}`** via the environment (e.g. `Authorization: Bearer ${OPENAI_API_KEY}`).

## Attack library

YAML lives under `pkg/attacks/library/`. Add files with:

```yaml
category: my_suite
attacks:
  - name: Short title
    description: Why this matters
    prompt: |
      Multi-line attack prompt
    tags: [injection]
```

Example extra pack: `examples/custom_pack/compliance.yaml`. Run with:

`./bin/redteam run --config examples/redteam-mock.yaml --attacks-dir examples/custom_pack --mock`

## Reports

- **JSON** — stable `schema_version`, full prompts/responses, findings (for SaaS ingestion later).
- **HTML** — single-file dark-theme summary for humans.

## CI

- **GitHub Actions**: `.github/workflows/redteam.yml` at the **repository root** (expects a `redteam/` subdirectory). If this project is the repo root, change `working-directory` and artifact `path` to `.`.
- **GitLab**: see `examples/gitlab-ci.yml`.

## TypeScript wrapper

`integrations/typescript` provides `runRedteamSync()` which spawns the binary (`redteam` on `PATH` or `REDTEAM_BIN`).

```bash
cd integrations/typescript && npm install && npm run build
```

## Framework integrations

The scanner does not import LangChain/LlamaIndex/Vercel directly; it exercises **your deployed HTTP contract**. Examples:

- `examples/langchain_openai_http.py` — minimal chain shape; wrap with an HTTP API for scanning.
- `examples/vercel_ai_sdk_snippet.ts` — notes for a non-streaming test route.

LlamaIndex: same pattern — expose the engine over HTTP, align `body_template` / `response_path` with your JSON.

## Detectors (Phase 1)

| Detector | Notes |
|-----------|--------|
| `pii` | Regex-style email, phone, SSN-shaped sequences, long digit runs, Stripe-like keys, Bearer tokens |
| `toxicity` | Small blocklist (replace with Perspective or a classifier in product hardening) |
| `disallowed_topics` | Substring probes for a few high-risk instruction classes |
| `context_leakage` | Verbatim-ish echo of distinctive context tokens |
| `jailbreak_success_heuristic` | Phrases that often indicate compliance with jailbreak prompts |

Tune rules for your product; false positives are expected with regex-only checks.

## Roadmap → paid SaaS

1. **OSS engine** (this repo) — CI adoption, attack packs, JSON schema stability.
2. **Cloud** — store runs, diff over time, schedules, team alerts, managed secret injection.
3. **Enterprise** — custom attack packs, compliance exports, SSO, VPC agents.

## License

MIT — see `LICENSE`.
