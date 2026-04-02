package runner

import "testing"

func TestExtractResponsePath(t *testing.T) {
	body := []byte(`{"choices":[{"message":{"content":"hello"}}]}`)
	got, err := extractResponsePath(body, "choices.0.message.content")
	if err != nil {
		t.Fatal(err)
	}
	if got != "hello" {
		t.Fatalf("got %q", got)
	}
}

func TestRenderBody(t *testing.T) {
	b, err := renderBody(`{"x":{{toJSON .Prompt}}}`, `say "hi"`, "")
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"x":"say \"hi\""}` {
		t.Fatalf("got %s", b)
	}
}
