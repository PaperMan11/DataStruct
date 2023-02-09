package unionfind

import (
	"testing"
)

func TestUnionFind(t *testing.T) {
	uf := NewUnionFind(10)
	uf.UnionElements(0, 1)
	uf.UnionElements(1, 2)
	uf.UnionElements(2, 3)
	uf.UnionElements(3, 4)
	uf.UnionElements(4, 5)
	uf.UnionElements(5, 6)
	uf.UnionElements(6, 7)
	uf.UnionElements(7, 8)

	if uf.Find(4) != 1 {
		t.Fatalf("root error")
	}

	if uf.Find(1) != uf.Find(7) {
		t.Fatalf("should equal")
	}

	if uf.IsConnected(1, 9) {
		t.Fatalf("should not be equal")
	}
}
