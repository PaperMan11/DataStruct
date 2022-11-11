package algorithm

func Kmp(s1, s2 string) int {
	var i, j = 0, 0
	nextArr := make([]int, len(s2))
	getNextArr(s2, nextArr)
	// fmt.Println(nextArr)
	for i < len(s1) && j < len(s2) {
		if j == -1 || s1[i] == s2[j] {
			i++
			j++
		} else {
			j = nextArr[j] //i不变,j后退，现在知道为什么这样让子串回退了吧
		}
	}
	if j >= len(s2) {
		return i - len(s2) //返回第一个匹配的下标
	}
	return -1
}

func getNextArr(str string, nextArr []int) {
	var (
		i int = 0
		k int = -1 // 数组值(即i坐标前str的最长公共前后缀)
	)
	nextArr[0] = -1
	for i < len(str)-1 {
		if k == -1 || str[i] == str[k] {
			i++
			k++
			if nextArr[i] != nextArr[k] {
				nextArr[i] = k
			} else {
				nextArr[i] = nextArr[k]
			}
		} else {
			k = nextArr[k]
		}
	}
}
