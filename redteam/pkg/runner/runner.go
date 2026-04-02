package runner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Ali627miya/llm-redteam/redteam/pkg/config"
	"github.com/Ali627miya/llm-redteam/redteam/pkg/detectors"
	"github.com/Ali627miya/llm-redteam/redteam/pkg/models"
)

// Invoker calls the application under test. HTTPInvoke implements this for real targets.
type Invoker func(ctx context.Context, cfg *config.Config, prompt string) (text string, status int, err error)

// Run executes attacks with bounded concurrency and returns per-attack results.
func Run(ctx context.Context, cfg *config.Config, attacks []models.AttackCase, invoke Invoker) []models.RunResult {
	if invoke == nil {
		invoke = HTTPInvoke
	}
	max := cfg.Run.MaxAttacks
	if max > 0 && max < len(attacks) {
		attacks = attacks[:max]
	}
	conc := cfg.Run.Concurrency
	if conc < 1 {
		conc = 1
	}

	type job struct {
		idx int
		a   models.AttackCase
	}
	jobs := make(chan job)
	results := make([]models.RunResult, len(attacks))

	var wg sync.WaitGroup
	worker := func() {
		defer wg.Done()
		for j := range jobs {
			if err := ctx.Err(); err != nil {
				results[j.idx] = models.RunResult{
					AttackID: j.a.ID, Category: j.a.Category, Name: j.a.Name,
					Prompt: j.a.Prompt, Error: err.Error(),
					StartedAt: time.Now(), CompletedAt: time.Now(),
				}
				continue
			}
			start := time.Now()
			text, status, err := invoke(ctx, cfg, j.a.Prompt)
			latency := time.Since(start).Milliseconds()
			rr := models.RunResult{
				AttackID:    j.a.ID,
				Category:    j.a.Category,
				Name:        j.a.Name,
				Prompt:      j.a.Prompt,
				Response:    text,
				LatencyMS:   latency,
				StartedAt:   start,
				CompletedAt: time.Now(),
				Metadata:    map[string]string{"http_status": fmt.Sprintf("%d", status)},
			}
			if err != nil {
				rr.Error = err.Error()
			} else {
				rr.Findings = detectors.Analyze(cfg, j.a, text)
			}
			results[j.idx] = rr
		}
	}

	for w := 0; w < conc; w++ {
		wg.Add(1)
		go worker()
	}
	for i, a := range attacks {
		jobs <- job{idx: i, a: a}
	}
	close(jobs)
	wg.Wait()
	return results
}
