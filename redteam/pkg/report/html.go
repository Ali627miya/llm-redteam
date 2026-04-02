package report

import (
	"fmt"
	"html"
	"os"
	"strings"
	"time"

	"github.com/Ali627miya/llm-redteam/redteam/pkg/models"
)

// WriteHTML renders a minimal self-contained report page.
func WriteHTML(path string, target string, results []models.RunResult) error {
	sum := BuildSummary(results)
	var b strings.Builder
	b.WriteString(`<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><title>LLM Red Team Report</title>`)
	b.WriteString(`<style>
body{font-family:system-ui,sans-serif;margin:2rem;max-width:1100px;background:#0f1419;color:#e6edf3}
h1{font-size:1.5rem} h2{font-size:1.1rem;margin-top:1.5rem}
.stats{display:flex;gap:1rem;flex-wrap:wrap;margin:1rem 0}
.stat{background:#161b22;border:1px solid #30363d;border-radius:8px;padding:0.75rem 1rem}
.critical{color:#ff7b72}.high{color:#ffa657}.ok{color:#3fb950}
.card{border:1px solid #30363d;border-radius:8px;padding:1rem;margin:0.75rem 0;background:#161b22}
.meta{color:#8b949e;font-size:0.85rem}
pre{white-space:pre-wrap;word-break:break-word;background:#010409;padding:0.75rem;border-radius:6px;border:1px solid #21262d;font-size:0.8rem}
.tag{display:inline-block;background:#21262d;padding:0.15rem 0.4rem;border-radius:4px;margin-right:0.25rem;font-size:0.75rem}
</style></head><body>`)
	b.WriteString(`<h1>LLM red team scan</h1>`)
	b.WriteString(`<p class="meta">Generated ` + html.EscapeString(time.Now().UTC().Format(time.RFC3339)) + ` · Target: ` + html.EscapeString(target) + `</p>`)
	b.WriteString(`<div class="stats">`)
	b.WriteString(fmt.Sprintf(`<div class="stat"><strong>%d</strong><br>attacks</div>`, sum.TotalAttacks))
	b.WriteString(fmt.Sprintf(`<div class="stat ok"><strong>%d</strong><br>clean</div>`, sum.Passed))
	b.WriteString(fmt.Sprintf(`<div class="stat critical"><strong>%d</strong><br>with findings</div>`, sum.Failed))
	b.WriteString(fmt.Sprintf(`<div class="stat high"><strong>%d</strong><br>errors</div>`, sum.Errors))
	b.WriteString(fmt.Sprintf(`<div class="stat"><strong>%d</strong><br>findings</div>`, sum.TotalFindings))
	b.WriteString(`</div>`)

	for _, r := range results {
		b.WriteString(`<div class="card">`)
		title := html.EscapeString(r.Name)
		if r.Error != "" {
			b.WriteString(fmt.Sprintf(`<h2 class="critical">%s</h2>`, title))
			b.WriteString(`<p class="meta">Error: ` + html.EscapeString(r.Error) + `</p>`)
		} else if len(r.Findings) > 0 {
			b.WriteString(fmt.Sprintf(`<h2 class="high">%s</h2>`, title))
		} else {
			b.WriteString(fmt.Sprintf(`<h2 class="ok">%s</h2>`, title))
		}
		b.WriteString(`<p class="meta">` + html.EscapeString(r.Category) + ` · ` + html.EscapeString(r.AttackID) + ` · ` + fmt.Sprintf("%d ms", r.LatencyMS) + `</p>`)
		b.WriteString(`<h3>Prompt</h3><pre>` + html.EscapeString(r.Prompt) + `</pre>`)
		if r.Response != "" {
			b.WriteString(`<h3>Response</h3><pre>` + html.EscapeString(trunc(r.Response, 4000)) + `</pre>`)
		}
		if len(r.Findings) > 0 {
			b.WriteString(`<h3>Findings</h3><ul>`)
			for _, f := range r.Findings {
				cls := "meta"
				if f.Severity == models.SeverityCritical {
					cls = "critical"
				} else if f.Severity == models.SeverityHigh {
					cls = "high"
				}
				b.WriteString(fmt.Sprintf(`<li class="%s"><span class="tag">%s</span> <strong>%s</strong> — %s`,
					cls, html.EscapeString(string(f.Severity)), html.EscapeString(f.Detector), html.EscapeString(f.Message)))
				if f.Evidence != "" {
					b.WriteString(` <code>` + html.EscapeString(f.Evidence) + `</code>`)
				}
				b.WriteString(`</li>`)
			}
			b.WriteString(`</ul>`)
		}
		b.WriteString(`</div>`)
	}
	b.WriteString(`</body></html>`)
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func trunc(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "\n…"
}
