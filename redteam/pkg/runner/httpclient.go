package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/Ali627miya/llm-redteam/redteam/pkg/config"
)

type bodyVars struct {
	Prompt  string
	Context string
}

func renderBody(tpl string, prompt, context string) ([]byte, error) {
	funcs := template.FuncMap{
		"toJSON": func(s string) (string, error) {
			b, err := json.Marshal(s)
			if err != nil {
				return "", err
			}
			return string(b), nil
		},
	}
	t, err := template.New("body").Funcs(funcs).Parse(tpl)
	if err != nil {
		return nil, fmt.Errorf("body_template parse: %w", err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, bodyVars{Prompt: prompt, Context: context}); err != nil {
		return nil, fmt.Errorf("body_template execute: %w", err)
	}
	return buf.Bytes(), nil
}

func extractResponsePath(body []byte, path string) (string, error) {
	if path == "" || path == "." {
		return string(body), nil
	}
	var v any
	if err := json.Unmarshal(body, &v); err != nil {
		return "", fmt.Errorf("response json: %w", err)
	}
	parts := strings.Split(path, ".")
	cur := v
	for _, p := range parts {
		if p == "" {
			continue
		}
		switch x := cur.(type) {
		case map[string]any:
			cur, _ = x[p]
		case []any:
			var idx int
			if _, err := fmt.Sscanf(p, "%d", &idx); err != nil || idx < 0 || idx >= len(x) {
				return "", fmt.Errorf("invalid array index %q in path %q", p, path)
			}
			cur = x[idx]
		default:
			return "", fmt.Errorf("cannot traverse %T at %q", cur, p)
		}
	}
	switch s := cur.(type) {
	case string:
		return s, nil
	case nil:
		return "", fmt.Errorf("empty value at path %q", path)
	default:
		b, err := json.Marshal(s)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
}

// HTTPInvoke posts to the configured target and returns model text and status.
func HTTPInvoke(ctx context.Context, cfg *config.Config, prompt string) (text string, status int, err error) {
	body, err := renderBody(cfg.Target.BodyTemplate, prompt, cfg.Context)
	if err != nil {
		return "", 0, err
	}
	ctxTimeout := time.Duration(cfg.Run.TimeoutSeconds) * time.Second
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, cfg.Target.Method, cfg.Target.URL, bytes.NewReader(body))
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range cfg.Target.Headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return "", resp.StatusCode, err
	}
	text, err = extractResponsePath(respBody, cfg.Target.ResponsePath)
	if err != nil {
		return string(respBody), resp.StatusCode, fmt.Errorf("extract response: %w", err)
	}
	return text, resp.StatusCode, nil
}
