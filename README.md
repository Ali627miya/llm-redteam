<p align="center">
  <img src="redteam/docs/logo.svg" alt="LLM Red Team" width="88" height="88" />
</p>

<h1 align="center">LLM Red Team</h1>

<p align="center"><strong>Automatically test your LLM apps for prompt injection, jailbreaks, and data leaks.</strong></p>

<p align="center">
  <a href="https://github.com/Ali627miya/llm-redteam/actions/workflows/ci.yml"><img src="https://github.com/Ali627miya/llm-redteam/actions/workflows/ci.yml/badge.svg" alt="CI" /></a>
  <a href="https://github.com/Ali627miya/llm-redteam/actions/workflows/redteam.yml"><img src="https://github.com/Ali627miya/llm-redteam/actions/workflows/redteam.yml/badge.svg" alt="Scan workflow" /></a>
  <a href="https://github.com/Ali627miya/llm-redteam/blob/main/redteam/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License" /></a>
  <a href="https://github.com/Ali627miya/llm-redteam/blob/main/redteam/go.mod"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white" alt="Go version" /></a>
  <a href="https://github.com/Ali627miya/llm-redteam"><img src="https://img.shields.io/github/stars/Ali627miya/llm-redteam?style=social" alt="GitHub stars" /></a>
</p>

---

## Why this exists

Shipping an LLM feature without probing **prompt injection**, **jailbreaks**, and **context leakage** is risky. **LLM Red Team** is a small **Go CLI** that runs a library of attacks against **your HTTP API** (OpenAI-shaped, Anthropic, Ollama, or your gateway), then scores responses with **PII / toxicity / topic / leakage heuristics** and writes **JSON + HTML** reports—ideal for **CI** and **security reviews**.

**No proxy.** **No vendor lock-in.** The engine is open source; a future **hosted** product can add history, schedules, SSO, and compliance exports.

## Quick start

### Option A — Install with Go

```bash
go install github.com/Ali627miya/llm-redteam/redteam/cmd/redteam@latest
redteam version
```

### Option B — Prebuilt binary (no Go)

Download the archive for your OS from [**Releases**](https://github.com/Ali627miya/llm-redteam/releases), extract `redteam` (or `redteam.exe`), and put it on your `PATH`.

> **First release:** after you publish, create a tag: `git tag v0.1.0 && git push origin v0.1.0` — [GoReleaser](https://goreleaser.com) builds **Linux, macOS, Windows** (amd64 + arm64) via [`.github/workflows/release.yml`](.github/workflows/release.yml).

### Option C — Clone and build

```bash
git clone https://github.com/Ali627miya/llm-redteam.git
cd llm-redteam/redteam
go build -o bin/redteam ./cmd/redteam
./bin/redteam list
```

### Run a CI-safe scan (no API keys, no network to your model)

```bash
cd redteam   # config paths are relative to this folder
redteam run --config examples/redteam-mock.yaml --output report.html --format html --mock --fail-on-findings
```

### One-shot local demo (mock HTTP target + scan)

```bash
cd redteam && ./scripts/run-local.sh
```

## Minimal configs (copy & adapt)

<details>
<summary><strong>OpenAI</strong> (Chat Completions)</summary>

```yaml
target:
  url: https://api.openai.com/v1/chat/completions
  method: POST
  headers:
    Authorization: Bearer ${OPENAI_API_KEY}
    Content-Type: application/json
  body_template: |
    {"model":"gpt-4o-mini","messages":[{"role":"user","content":{{toJSON .Prompt}}}],"temperature":0.2}
  response_path: choices.0.message.content
```

Full example: [`redteam/examples/redteam.yaml`](redteam/examples/redteam.yaml).

</details>

<details>
<summary><strong>Anthropic</strong> (Messages API)</summary>

See [`redteam/examples/redteam.anthropic.yaml`](redteam/examples/redteam.anthropic.yaml) — set `ANTHROPIC_API_KEY`.

</details>

<details>
<summary><strong>Local Ollama</strong></summary>

See [`redteam/examples/redteam.ollama.yaml`](redteam/examples/redteam.ollama.yaml) — default `http://127.0.0.1:11434/v1/chat/completions`.

</details>

## Live demo & video

- **Static marketing + sample report:** [`redteam/website/`](redteam/website/) — deploy `redteam/website` to Vercel/Netlify or GitHub Pages ([`redteam/DEMO.md`](redteam/DEMO.md)).
- **Terminal recording:** add your [asciinema](https://asciinema.org/) or YouTube link here after you record one (2–3 minutes: `run-local.sh` + HTML report).

> The CLI does **not** phone home. Usage analytics would only be added with an **explicit opt-in** flag and clear documentation.

## GitHub Actions

- **CI:** [`.github/workflows/ci.yml`](.github/workflows/ci.yml) — `go test` + build.
- **Example scan job:** [`.github/workflows/redteam.yml`](.github/workflows/redteam.yml) — mock scan on every push.
- **Reusable workflow (install + scan):** [`.github/actions/redteam/action.yml`](.github/actions/redteam/action.yml)

Example using the composite action **from this repo** (after `actions/checkout`):

```yaml
- uses: actions/checkout@v4
- uses: ./.github/actions/redteam
  with:
    ref: main
    working-directory: redteam
    config: examples/redteam-mock.yaml
    extra-args: "--mock --fail-on-findings"
```

From **another** repository, put `redteam.yaml` in your repo and pin a tag:

```yaml
- uses: actions/checkout@v4
- uses: Ali627miya/llm-redteam/.github/actions/redteam@v0.1.0
  with:
    ref: v0.1.0
    working-directory: .
    config: redteam.yaml
```

(Replace `v0.1.0` after your first [Release](https://github.com/Ali627miya/llm-redteam/releases). Until then, `uses: ...@main` and `ref: main` work for `go install`.)

## “Powered by LLM Red Team”

Add to your README:

```markdown
[![LLM Red Team](https://img.shields.io/badge/security-LLM%20Red%20Team-3ee8b5?logo=shield)](https://github.com/Ali627miya/llm-redteam)
```

## Roadmap (feedback-driven)

- More providers & **attack packs** (YAML — contributions welcome).
- **Custom detection rules** in YAML.
- **OWASP LLM Top 10** mapping in reports.
- **Cloud** waitlist: open an issue or star the repo — we’ll link a proper waitlist when ready.

## Docs in this repo

| Doc | Purpose |
|-----|---------|
| [`redteam/README.md`](redteam/README.md) | Technical details, detectors, project layout |
| [`LAUNCH.md`](LAUNCH.md) | HN, Reddit, outreach checklist |
| [`CONTRIBUTING.md`](CONTRIBUTING.md) | PRs, issues, attack pack format |
| [`SECURITY.md`](SECURITY.md) | Reporting vulnerabilities |
| [`redteam/DEMO.md`](redteam/DEMO.md) | Buyer demo & hosting the static site |

## License

MIT — see [`redteam/LICENSE`](redteam/LICENSE).

## Star history

If the project helps you, a star on [GitHub](https://github.com/Ali627miya/llm-redteam) helps others discover it.
