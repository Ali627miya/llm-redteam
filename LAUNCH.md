# Launch checklist — LLM Red Team

Use this after your first **GitHub Release** (`v0.1.0`) is live.

## Before you post

- [ ] Tag release: `git tag v0.1.0 && git push origin v0.1.0` (triggers [GoReleaser](.github/workflows/release.yml)).
- [ ] Confirm **Releases** page has binaries for Linux / macOS / Windows.
- [ ] Record **2–3 min demo** (asciinema or YouTube): clone → `run-local.sh` or `--mock` → open HTML report.
- [ ] Add the demo link to the [root README](README.md) “Live demo & video” section.
- [ ] Deploy `redteam/website` (Vercel / Netlify / Pages) and add URL to README if you want a public landing page.
- [ ] Optional: simple **cloud waitlist** (Google Form, Buttondown, or Mailchimp) and link from README + website.

## Communities (respect each subreddit’s self-promo rules)

| Channel | Suggested angle |
|---------|-----------------|
| **Hacker News** | “Show HN: LLM Red Team — open-source security testing for LLM apps (Go CLI, CI-friendly)” |
| **r/golang** | Technical: architecture, `go install`, no proxy |
| **r/LLMDevs** / **r/LocalLLaMA** | Ollama example config, local testing |
| **r/cybersecurity** / **r/netsec** | Position as *shift-left* probe tool, not a full pentest |
| **r/devops** | GitHub Action + `fail-on-findings` gate |
| **Discord** | LangChain, Latent Space, OWASP — one short message + link, offer to help wire HTTP target |

## Social

- Short clip + link; tag a few AI/security folks **sparingly** (quality over spam).

## Blog post outline

1. Problem: LLM features ship fast; abuse tests are ad hoc.  
2. Solution: HTTP-pluggable scanner + attack library + HTML/JSON for audits.  
3. Architecture: Go CLI, YAML packs, detectors, no proxy.  
4. How to run: `go install`, Releases binary, GitHub Action.  
5. Roadmap: SaaS, OWASP mapping, more providers.  
Publish on **Dev.to**, **Medium**, or your own site; cross-post summary to HN.

## Early adopters

- DM 10 AI startups: offer **free feedback session** using their staging API + your tool; ask for a **testimonial** or GitHub star if it helps.
- “Hall of Fame” in README: list companies **with permission**.

## Feedback loop

- **GitHub Issues** for bugs and features.  
- Ask: *What’s missing? Biggest pain in LLM security testing?*  
- Prioritize: new providers, attack templates, reusable Action improvements, rule packs.

## Long term

- Detection rules as YAML; OWASP LLM Top 10 mapping; SOC2-oriented export narratives (with legal review).  
- Monetization when traction is real: hosted dashboard, SSO, org policies — **not** before you have users.
