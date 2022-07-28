//go:build unit
// +build unit

package pathgroup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathGroup(t *testing.T) {
	pathGroups := []string{
		"/foo/bar/:p1",
		"/foo/:p1",
		"/foo",
		"/foo/:p1/bar",
		"/baz/:p1/foo/:p2/bar/:p3",
		"/baz/:p1/qux/:p2/bar/:p3",
	}

	tests := []struct {
		test    string
		pattern string
		ok      bool
	}{
		{test: "/foo/U989nnA-IAAJAJH", pattern: "/foo/:param#1", ok: true},
		{test: "/foo/U989nnA-IAAJA?`[]=+JH/bar", pattern: "/foo/:param#1/bar", ok: true},
		{test: "/foo/", pattern: "/foo", ok: true},
		{test: "/foo", pattern: "/foo", ok: true},
		{test: "/foo/U989nnA/4554541aa", pattern: "/foo/U989nnA/4554541aa", ok: false},
		{test: "/foo/U989nnA-IAAJA?`[]=+JH", pattern: "/foo/:param#1", ok: true},
		{test: "/foo/bar/09090990909sasa", pattern: "/foo/bar/:param#1", ok: true},
		{test: "/foo/bar/09090990909/sasa", pattern: "/foo/bar/09090990909/sasa", ok: false},
		{test: "/foo/bar/", pattern: "/foo/bar", ok: true},
		{test: "/baz/1/foo/2/bar/3", pattern: "/baz/:param#1/foo/:param#2/bar/:param#3", ok: true},
		{test: "/baz/1/qux/2/bar/3", pattern: "/baz/:param#1/qux/:param#2/bar/:param#3", ok: true},
		{test: "/baz/1/qux/2/bar/", pattern: "/baz/:param#1/qux/:param#2/bar", ok: true},
	}

	groups := New(pathGroups)
	for _, tc := range tests {
		p, ok := groups.Find(tc.test)
		assert.Equal(t, tc.pattern, p)
		assert.Equal(t, tc.ok, ok)
	}
}
