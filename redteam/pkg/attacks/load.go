package attacks

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/llm-redteam/redteam/pkg/models"
)

// ParseAttackYAML parses a single attack pack file and assigns IDs using idPrefix (usually file stem).
func ParseAttackYAML(raw []byte, idPrefix string) ([]models.AttackCase, error) {
	var file struct {
		Category string              `yaml:"category"`
		Attacks  []models.AttackCase `yaml:"attacks"`
	}
	if err := yaml.Unmarshal(raw, &file); err != nil {
		return nil, err
	}
	stem := strings.ToLower(strings.TrimSpace(idPrefix))
	for i := range file.Attacks {
		if file.Attacks[i].Category == "" {
			file.Attacks[i].Category = file.Category
		}
		if file.Attacks[i].ID == "" {
			file.Attacks[i].ID = fmt.Sprintf("%s-%d", stem, i+1)
		}
	}
	return file.Attacks, nil
}

// LoadDirectory walks dir recursively for *.yaml / *.yml. If categories is non-empty, only files whose
// stem (basename without extension) is in the set are loaded (same rule as built-in library).
func LoadDirectory(dir string, categories []string) ([]models.AttackCase, error) {
	filter := map[string]struct{}{}
	for _, c := range categories {
		filter[strings.ToLower(strings.TrimSpace(c))] = struct{}{}
	}
	var all []models.AttackCase
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
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
		raw, err := os.ReadFile(path)
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

// LoadBuiltIn is implemented in embed.go and calls the shared parser.
