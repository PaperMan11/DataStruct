package sort

func QuickSort(arr []int64) {
	quickSort(arr, 0, len(arr)-1)
}

func quickSort(arr []int64, start, end int) {
	if start < end {
		i, j := start, end
		key := arr[(start+end)/2]
		for i <= j {
			for arr[i] < key { // 找到左边第一个比key大的
				i++
			}
			for arr[j] > key { // 找到右边第一个比key小的
				j--
			}
			if i <= j {
				arr[i], arr[j] = arr[j], arr[i]
				i++
				j--
			}
		}
		if start < j {
			quickSort(arr, start, j)
		}
		if end > i {
			quickSort(arr, i, end)
		}
	}
}
