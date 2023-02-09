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

// 三路快速排序
func QuickSort3Ways(arr []int64) {
	quickSort3Ways(arr, 0, len(arr)-1)
}

func quickSort3Ways(arr []int64, l, r int) {
	if l >= r {
		return
	}

	var (
		lt    = l     // arr[l+1...lt] < pivot
		gt    = r + 1 // arr[gt...r] > pivot
		i     = l     // arr[lt+1...i) == pivot
		pivot = arr[l]
	)
	for i < gt {
		if arr[i] < pivot {
			swap(arr, l, i)
			i++
			lt++
		} else if arr[i] > pivot {
			swap(arr, i, gt-1)
			gt--
		} else {
			i++
		}
	}
	if arr[l] > arr[lt] {
		swap(arr, l, lt)
	}
	quickSort3Ways(arr, l, lt)
	quickSort3Ways(arr, gt, r)
}

func swap(arr []int64, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

