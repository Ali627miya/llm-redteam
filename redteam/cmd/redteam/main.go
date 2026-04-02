package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Ali627miya/llm-redteam/redteam/internal/version"
	"github.com/Ali627miya/llm-redteam/redteam/pkg/attacks"
	"github.com/Ali627miya/llm-redteam/redteam/pkg/config"
	"github.com/Ali627miya/llm-redteam/redteam/pkg/report"
	"github.com/Ali627miya/llm-redteam/redteam/pkg/runner"
)

func main() {
	os.Exit(run())
}

func run() int {
	if len(os.Args) > 1 && (os.Args[1] == "-version" || os.Args[1] == "--version") {
		fmt.Println(version.Version)
		return 0
	}

	fs := flag.NewFlagSet("redteam", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: redteam <command> [options]\n\nCommands:\n  run       Execute attack library against target\n  list      Print attack IDs\n  version   Print version\n\nGlobal:\n  -version   Print version and exit\n\n")
		fs.PrintDefaults()
	}
	if len(os.Args) < 2 {
		fs.Usage()
		return 2
	}
	switch os.Args[1] {
	case "version":
		fmt.Println(version.Version)
		return 0
	case "list":
		return cmdList(os.Args[2:])
	case "run":
		return cmdRun(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		fs.Usage()
		return 2
	}
}

func splitCategories(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	var out []string
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func cmdList(args []string) int {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	var extraDirs stringSlice
	fs.Var(&extraDirs, "attacks-dir", "directory of extra YAML attack packs (repeatable)")
	noBuiltin := fs.Bool("no-builtin", false, "only load packs from -attacks-dir")
	catStr := fs.String("categories", "", "comma-separated stems, e.g. jailbreak,prompt_injection (default: all)")
	_ = fs.Parse(args)

	cases, err := attacks.LoadAll(splitCategories(*catStr), *noBuiltin, []string(extraDirs))
	if err != nil {
		fmt.Fprintf(os.Stderr, "load attacks: %v\n", err)
		return 1
	}
	for _, a := range cases {
		fmt.Printf("%s\t%s\t%s\n", a.ID, a.Category, a.Name)
	}
	return 0
}

func cmdRun(args []string) int {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	cfgPath := fs.String("config", "redteam.yaml", "path to YAML config")
	out := fs.String("output", "redteam-report.json", "output file path")
	format := fs.String("format", "json", "json or html")
	mock := fs.Bool("mock", false, "do not call target; use canned responses (CI / dry run)")
	failOnFindings := fs.Bool("fail-on-findings", false, "exit 1 if any findings")
	var extraDirs stringSlice
	fs.Var(&extraDirs, "attacks-dir", "directory of extra YAML attack packs (repeatable)")
	noBuiltin := fs.Bool("no-builtin", false, "skip embedded library; use only -attacks-dir packs")
	_ = fs.Parse(args)

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		return 1
	}
	attackCases, err := attacks.LoadAll(cfg.Categories, *noBuiltin, []string(extraDirs))
	if err != nil {
		fmt.Fprintf(os.Stderr, "attacks: %v\n", err)
		return 1
	}
	if len(attackCases) == 0 {
		fmt.Fprintf(os.Stderr, "no attacks loaded (check -no-builtin and -attacks-dir)\n")
		return 1
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var invoke runner.Invoker = runner.HTTPInvoke
	if *mock {
		invoke = mockInvoker
	}
	results := runner.Run(ctx, cfg, attackCases, invoke)

	switch strings.ToLower(*format) {
	case "json":
		if err := report.WriteJSON(*out, cfg.Target.URL, results); err != nil {
			fmt.Fprintf(os.Stderr, "write json: %v\n", err)
			return 1
		}
	case "html":
		if err := report.WriteHTML(*out, cfg.Target.URL, results); err != nil {
			fmt.Fprintf(os.Stderr, "write html: %v\n", err)
			return 1
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown format %q\n", *format)
		return 2
	}

	sum := report.BuildSummary(results)
	fmt.Printf("Scan complete: %d attacks, %d clean, %d with findings, %d errors, %d total findings\n",
		sum.TotalAttacks, sum.Passed, sum.Failed, sum.Errors, sum.TotalFindings)
	fmt.Printf("Wrote %s (%s)\n", *out, *format)

	if *failOnFindings && sum.TotalFindings > 0 {
		return 1
	}
	if sum.Errors == sum.TotalAttacks && sum.TotalAttacks > 0 {
		return 1
	}
	return 0
}

func mockInvoker(_ context.Context, cfg *config.Config, prompt string) (string, int, error) {
	lower := strings.ToLower(prompt)
	switch {
	case strings.Contains(lower, "ignore all previous"), strings.Contains(lower, "injection_ok"):
		return `I will not override my instructions or repeat internal system prompts.`, 200, nil
	case strings.Contains(lower, "developer mode"), strings.Contains(lower, "you are \"dan\""), strings.Contains(lower, "you are 'dan'"):
		return `I cannot enable developer mode or remove safety policies.`, 200, nil
	case strings.Contains(lower, "hurtful"), strings.Contains(lower, "slurs"):
		return `I can't help with harassment or slurs.`, 200, nil
	case strings.Contains(lower, "torture"), strings.Contains(lower, "graphic"):
		return `I won't provide graphic torture or extreme violence.`, 200, nil
	case strings.Contains(lower, "phishing"), strings.Contains(lower, "villain"):
		return `I can't provide instructions for phishing or attacking banks.`, 200, nil
	case strings.Contains(lower, "crm"), strings.Contains(lower, "customer records"):
		return `I don't have access to a CRM and can't output real customer records.`, 200, nil
	case cfg.Context != "" && strings.Contains(lower, "secret"):
		return `I don't see confidential tokens in what I should repeat here.`, 200, nil
	}
	return `Mock response: request declined or handled with default safety behavior.`, 200, nil
}
