package attacks

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/Ali627miya/llm-redteam/redteam/pkg/models"
)

//go:embed library/*.yaml
var libraryFS embed.FS

// LoadBuiltIn returns all attack cases from the embedded library, optionally filtered by category name (file stem).
func LoadBuiltIn(categories []string) ([]models.AttackCase, error) {
	filter := map[string]struct{}{}
	for _, c := range categories {
		filter[strings.ToLower(strings.TrimSpace(c))] = struct{}{}
	}
	var all []models.AttackCase
	err := fs.WalkDir(libraryFS, "library", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml") {
			return nil
		}
		stem := strings.ToLower(strings.TrimSuffix(strings.TrimSuffix(filepath.Base(path), ".yaml"), ".yml"))
		if len(filter) > 0 {
			if _, ok := filter[stem]; !ok {
				return nil
			}
		}
		raw, err := libraryFS.ReadFile(path)
		if err != nil {
			return err
		}
		cases, err := ParseAttackYAML(raw, stem)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		all = append(all, cases...)
		return nil
	})
	return all, err
}
