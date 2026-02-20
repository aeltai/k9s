package mc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInjectContext(t *testing.T) {
	tests := map[string]struct {
		args    []string
		ctx     string
		expect  []string
	}{
		"simple": {
			args:   []string{"get", "nodes"},
			ctx:    "my-ctx",
			expect: []string{"get", "nodes", "--context", "my-ctx"},
		},
		"with double dash": {
			args:   []string{"exec", "pod", "--", "ls"},
			ctx:    "my-ctx",
			expect: []string{"exec", "pod", "--context", "my-ctx", "--", "ls"},
		},
		"empty args": {
			args:   []string{},
			ctx:    "c1",
			expect: []string{"--context", "c1"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := injectContext(tc.args, tc.ctx)
			assert.Equal(t, tc.expect, got)
		})
	}
}

func TestFormatResults(t *testing.T) {
	results := []Result{
		{Context: "ctx-1", Output: "node1 Ready\nnode2 Ready\n"},
		{Context: "ctx-2", Err: assert.AnError, Output: "connection refused"},
	}
	out := FormatResults(results)
	assert.Contains(t, out, "ctx-1")
	assert.Contains(t, out, "node1 Ready")
	assert.Contains(t, out, "ctx-2")
	assert.Contains(t, out, "(error) connection refused")
}
