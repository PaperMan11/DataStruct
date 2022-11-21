package sort

func MergeSort(arr []int64) []int64 {
	if len(arr) < 2 {
		return arr
	}
	mid := len(arr) / 2
	leftArr := MergeSort(arr[:mid])
	rightArr := MergeSort(arr[mid:])
	return mergeSort(leftArr, rightArr)
}

func mergeSort(leftArr, rightArr []int64) []int64 {
	l_len := len(leftArr)
	r_len := len(rightArr)
	all_len := l_len + r_len
	var (
		l   = 0
		r   = 0
		res = make([]int64, all_len)
	)
	for i := 0; i < all_len; i++ {
		if l >= l_len {
			res[i] = rightArr[r]
			r++
		} else if r >= r_len {
			res[i] = leftArr[l]
			l++
		} else if leftArr[l] <= rightArr[r] {
			res[i] = leftArr[l]
			l++
		} else {
			res[i] = rightArr[r]
			r++
		}
	}
	return res
}

// 第二种写法
// func msort(arr []int, left int, right int) {
// 	if left >= right {
// 		return
// 	}
// 	mid := (left + right) / 2
// 	msort(arr, left, mid)
// 	msort(arr, mid+1, right)
// 	merge(arr, left, mid, right)
// }

// func merge(arr []int, left int, mid int, right int) {
// 	left1 := left
// 	left2 := mid + 1
// 	var arr2 [100]int
// 	k := 0
// 	for left1 <= mid && left2 <= right {
// 		if arr[left1] < arr[left2] {
// 			arr2[k] = arr[left1]
// 			k++
// 			left1++
// 		} else {
// 			arr2[k] = arr[left2]
// 			k++
// 			left2++
// 		}
// 	}
// 	for left1 <= mid {
// 		arr2[k] = arr[left1]
// 		k++
// 		left1++
// 	}
// 	for left2 <= right {
// 		arr2[k] = arr[left2]
// 		k++
// 		left2++
// 	}
// 	for i := 0; i < k; i++ {
// 		arr[i+left] = arr2[i]
// 	}
// }
