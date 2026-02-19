package configo

import (
	"fmt"
	"strings"
)

// Flatten converts a nested map to a flat map with dot-notation keys.
func Flatten(m map[string]any) map[string]any {
	out := make(map[string]any)
	flatten("", m, out)
	return out
}

func flatten(prefix string, m map[string]any, out map[string]any) {
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch val := v.(type) {
		case map[string]any:
			if len(val) == 0 {
				out[key] = val
			} else {
				flatten(key, val, out)
			}
		case map[any]any:
			converted := convertMap(val)
			if len(converted) == 0 {
				out[key] = converted
			} else {
				flatten(key, converted, out)
			}
		default:
			out[key] = v
		}
	}
}

func convertMap(m map[any]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[fmt.Sprintf("%v", k)] = v
	}
	return out
}

// Unflatten converts a flat dot-notation map to a nested map.
func Unflatten(m map[string]any) map[string]any {
	out := make(map[string]any)
	for k, v := range m {
		parts := strings.Split(k, ".")
		current := out
		for i, part := range parts {
			if i == len(parts)-1 {
				current[part] = v
			} else {
				if next, ok := current[part]; ok {
					if nextMap, ok := next.(map[string]any); ok {
						current = nextMap
					} else {
						newMap := make(map[string]any)
						current[part] = newMap
						current = newMap
					}
				} else {
					newMap := make(map[string]any)
					current[part] = newMap
					current = newMap
				}
			}
		}
	}
	return out
}
