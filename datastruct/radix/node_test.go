package radix

import (
	"fmt"
	"strings"
	"testing"
)

func TestRadixNode(t *testing.T) {
	patterns := []string{
		"api/v1/user/:id/post/:pid/paper",
		"api/v1/user/:id/post",
		"api/v2/user/:id/post/:pid/paper",
		"api/v1/user/post/:pid/paper",
		"api/v1/user/post/paper",
	}
	root := newNode("")
	for _, pattern := range patterns {
		parts := strings.Split(pattern, "/")
		fmt.Println(parts)
		root.insert(pattern, parts, 0)
	}
	i := root.inOrder("")
	t.Log(i...)
	m, b := root.mate("/api/v1/user/1/post")
	t.Log(m, b)

	// root.erase([]string{"api", "v1", "user", ":id", "post"}, 0)
	ok := root.delete([]string{"api", "v1", "user", ":id", "post"}, 0)
	t.Log(ok)
	i = root.inOrder("")
	t.Log(i...)
}
