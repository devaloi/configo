package configo

import (
	"reflect"
	"testing"
)

func TestFlatten(t *testing.T) {
	tests := []struct {
		name string
		in   map[string]any
		want map[string]any
	}{
		{
			name: "flat map unchanged",
			in:   map[string]any{"key": "value"},
			want: map[string]any{"key": "value"},
		},
		{
			name: "nested map",
			in: map[string]any{
				"database": map[string]any{
					"host": "localhost",
					"port": 5432,
				},
			},
			want: map[string]any{
				"database.host": "localhost",
				"database.port": 5432,
			},
		},
		{
			name: "deep nesting",
			in: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": "deep",
					},
				},
			},
			want: map[string]any{"a.b.c": "deep"},
		},
		{
			name: "empty nested map",
			in: map[string]any{
				"empty": map[string]any{},
			},
			want: map[string]any{"empty": map[string]any{}},
		},
		{
			name: "mixed types",
			in: map[string]any{
				"server": map[string]any{
					"port": 8080,
					"ssl":  true,
				},
				"name": "app",
			},
			want: map[string]any{
				"server.port": 8080,
				"server.ssl":  true,
				"name":        "app",
			},
		},
		{
			name: "array values preserved",
			in: map[string]any{
				"tags": []string{"a", "b"},
			},
			want: map[string]any{
				"tags": []string{"a", "b"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Flatten(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Flatten() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnflatten(t *testing.T) {
	tests := []struct {
		name string
		in   map[string]any
		want map[string]any
	}{
		{
			name: "flat key unchanged",
			in:   map[string]any{"key": "value"},
			want: map[string]any{"key": "value"},
		},
		{
			name: "dot-notation to nested",
			in: map[string]any{
				"database.host": "localhost",
				"database.port": 5432,
			},
			want: map[string]any{
				"database": map[string]any{
					"host": "localhost",
					"port": 5432,
				},
			},
		},
		{
			name: "deep nesting",
			in:   map[string]any{"a.b.c": "deep"},
			want: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": "deep",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Unflatten(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unflatten() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlattenUnflattenRoundTrip(t *testing.T) {
	original := map[string]any{
		"database": map[string]any{
			"host": "localhost",
			"port": 5432,
		},
		"server": map[string]any{
			"name": "app",
		},
	}

	flat := Flatten(original)
	restored := Unflatten(flat)

	if !reflect.DeepEqual(original, restored) {
		t.Errorf("round trip failed: got %v, want %v", restored, original)
	}
}
