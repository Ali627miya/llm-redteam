package attacks

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDirectory(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "custom.yaml")
	content := `category: x
attacks:
  - name: T1
    prompt: hello
`
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cases, err := LoadDirectory(dir, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) != 1 || cases[0].Name != "T1" || cases[0].ID != "custom-1" {
		t.Fatalf("%+v", cases)
	}
}

func TestLoadAllSkipBuiltin(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "solo.yaml")
	if err := os.WriteFile(p, []byte(`category: z
attacks:
  - name: Only
    prompt: x
`), 0o644); err != nil {
		t.Fatal(err)
	}
	cases, err := LoadAll(nil, true, []string{dir})
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) != 1 || cases[0].Name != "Only" {
		t.Fatalf("%+v", cases)
	}
}
