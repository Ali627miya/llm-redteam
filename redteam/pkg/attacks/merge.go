package attacks

import (
	"strings"

	"github.com/Ali627miya/llm-redteam/redteam/pkg/models"
)

// LoadAll returns built-in attacks (unless skipBuiltin) plus every YAML file under extraDirs (recursive).
// categories filters by YAML file stem for both built-in and on-disk packs.
func LoadAll(categories []string, skipBuiltin bool, extraDirs []string) ([]models.AttackCase, error) {
	var out []models.AttackCase
	if !skipBuiltin {
		b, err := LoadBuiltIn(categories)
		if err != nil {
			return nil, err
		}
		out = append(out, b...)
	}
	for _, dir := range extraDirs {
		dir = strings.TrimSpace(dir)
		if dir == "" {
			continue
		}
		c, err := LoadDirectory(dir, categories)
		if err != nil {
			return nil, err
		}
		out = append(out, c...)
	}
	return out, nil
}
