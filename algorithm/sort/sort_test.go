package sort

import (
	"fmt"
	"testing"
)

func TestQSort(t *testing.T) {
	arr := []int64{2, 3, 4, 5, 6, 7, 6}
	QuickSort(arr)
	fmt.Println(arr)
}

func TestMSort(t *testing.T) {
	arr := []int64{2, 3, 4, 5, 6, 7, 6}
	arr = MergeSort(arr)
	fmt.Println(arr)
}
