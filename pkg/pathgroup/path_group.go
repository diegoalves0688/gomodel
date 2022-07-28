package pathgroup

import (
	"fmt"
	"regexp"
	"strings"
)

type group map[string]node

type node struct {
	nodes group
}

type PathGroup struct {
	nodes group
}

// New assemble and return a PathGroup (tree representation of a group of paths),
// which can be used to filter and group similar requests. Paths are expected to
// follow the pattern "/my/path/:param1/to/:param2".
func New(paths []string) PathGroup {
	return PathGroup{
		nodes: newGroup(paths),
	}
}

func newGroup(paths []string) group {
	// match resource params (e.g., /my/path/:id1/to/:id2)
	r := regexp.MustCompile(`^\:.+$`)

	nodes := group{}

	for _, p := range paths {
		tokens := strings.Split(p, "/")
		s := nodes
		nParams := 0
		for _, t := range tokens {
			// skip all empty tokens
			if t == "" {
				continue
			}
			// check for a parameter and change the current token
			// to a generic placeholder
			if r.MatchString(t) {
				nParams++
				t = fmt.Sprintf(":param#%d", nParams)
			}
			// check for an existing branch, otherwise create a new one
			if n, ok := s[t]; !ok {
				s[t] = node{
					nodes: group{},
				}
				s = s[t].nodes
			} else {
				s = n.nodes
			}
		}
	}
	return nodes
}

// Find returns the group pattern that was matched entirely against the given path.
// For example, the following pattern would be matched successfully with the
// following path:
//
//   PATTERN: /shelves/:id1/books/:id2
//   PATH: /shelves/2/books/1287/pages/42
//
// If there's no match, then the original path is returned instead of the pattern.
// An auxiliary boolean is returned alongside indicating wether there was a match
// or not.
func (g PathGroup) Find(path string) (string, bool) {
	tokens := strings.Split(path, "/")
	s := g.nodes
	nParams := 0
	walked := ""
	for _, t := range tokens {
		// skip all empty tokens
		if t == "" {
			continue
		}
		// check for an existing path
		n, ok := s[t]
		if !ok {
			// if didn't find, check for a generic param
			nParams++
			t = fmt.Sprintf(":param#%d", nParams)
			n, ok := s[t]
			if !ok {
				// if didn't find, then it's not a match
				return path, false
			}
			s = n.nodes
		} else {
			s = n.nodes
		}
		// keep track of the processed path
		walked = fmt.Sprintf("%s/%s", walked, t)
	}

	return walked, true
}
