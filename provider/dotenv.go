package provider

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// DotEnv loads configuration from a .env file.
type DotEnv struct {
	Path string
}

func NewDotEnv(path string) *DotEnv {
	return &DotEnv{Path: path}
}

func (p *DotEnv) Load() (map[string]any, error) {
	f, err := os.Open(p.Path)
	if err != nil {
		return nil, fmt.Errorf("dotenv provider: %w", err)
	}
	defer func() { _ = f.Close() }()

	out := make(map[string]any)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)
		out[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("dotenv provider: %w", err)
	}
	return out, nil
}
