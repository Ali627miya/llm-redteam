# Buyer / customer demo guide

Use this when presenting **llm-redteam** to prospects, investors, or security buyers.

## What to show (5 minutes)

1. **Landing page** (`website/index.html`) — problem, open-core model, path to SaaS.
2. **Sample report** (`website/sample-report.html`) — tangible output: findings, severities, prompts/responses (synthetic).
3. **Live terminal** (optional) — run `./scripts/run-local.sh` and open the generated `report.html` to prove the CLI is real.

**One-liner positioning:** “We automate LLM red-teaming in CI and produce audit-friendly reports; the engine is open source, the moat is hosted history, collaboration, and enterprise controls.”

## Host the marketing site

### Option A — Vercel

1. Push the repo to GitHub (or GitLab).
2. [Vercel](https://vercel.com) → New Project → import the repo.
3. Set **Root Directory** to `redteam/website`.
4. Framework preset: **Other** (static). Deploy.

Custom domain: Project → Settings → Domains.

### Option B — Netlify

1. Netlify → Add new site → Import from Git.
2. Base directory: `redteam/website`, publish directory: `.` (same folder).
3. Deploy.

### Option C — GitHub Pages

1. Repository **Settings → Pages**.
2. Source: **GitHub Actions** or deploy the `website/` folder as static assets from a workflow that uploads `redteam/website` contents to `gh-pages`.

### Option D — Docker (local screen-share)

From the repository root:

```bash
docker compose -f redteam/docker-compose.website.yml up
```

Open **http://localhost:8080** during a Zoom call.

## Before the call — checklist

- [ ] Replace placeholder GitHub URL in `website/index.html` with your real repo.
- [ ] Add your email or Calendly link in the CTAs if you want inbound.
- [ ] Run `./scripts/run-local.sh` once to ensure Go and ports work on the demo machine.
- [ ] Optional: print or PDF the sample report page for offline review.

## Docker image (CLI)

Build the scanner image for technical buyers:

```bash
docker build -t llm-redteam ./redteam
docker run --rm llm-redteam version
docker run --rm -v "$(pwd)/redteam/examples:/app/examples:ro" llm-redteam run \
  --config /app/examples/redteam-mock.yaml --output /tmp/r.json --format json --mock
```

(Adjust volume mounts for your paths; write output to a mounted volume if you need the file on the host.)

## Legal tone

The footer states this is **not** a full pentest or certification. Keep that for enterprise conversations; upsell professional services or deeper tooling separately if needed.
