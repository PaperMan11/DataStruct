package algorithm

import (
	"fmt"
	"testing"
)

func TestKmp(t *testing.T) {
	str1 := "ABABDABABAE"
	str2 := "BDA" // nextArr [-1 0 0]

	fmt.Println(Kmp(str1, str2))
}
