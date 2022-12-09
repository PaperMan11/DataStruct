package bloom

const bitNum = (32 << (^uint(0) >> 63)) / 8 // (自动判断32位或64位) / 8

type BitMap struct {
	bits []byte
	vmax uint // 当前已存储的最大边界
}

// NewBitMap 创建 bitmap
func NewBitMap(max_val ...uint) *BitMap {
	var (
		bitMap *BitMap = &BitMap{}
		max    uint    = 8192
	)
	if len(max_val) > 0 && max_val[0] > 0 {
		max = max_val[0]
	}

	sz := (max + bitNum - 1) / bitNum
	bitMap.bits = make([]byte, sz)
	bitMap.vmax = max
	return bitMap
}

func (bm *BitMap) Set(num uint) {
	if num >= bm.vmax { // 扩容
		bm.vmax += 1024
		if bm.vmax < num {
			bm.vmax = num
		}

		dd := int((bm.vmax+bitNum-1)/bitNum) - len(bm.bits) // bitmap 需要增加的长度
		if dd > 0 {
			tmp_arr := make([]byte, dd)
			bm.bits = append(bm.bits, tmp_arr...)
		}
	}

	bm.bits[num/bitNum] |= 1 << (num % bitNum)
}

func (bm *BitMap) ReSet(num uint) {
	if num >= bm.vmax {
		return
	}
	bm.bits[num/bitNum] &^= 1 << (num % bitNum)
}

func (bm *BitMap) Check(num uint) bool {
	if num >= bm.vmax {
		return false
	}
	return (bm.bits[num/bitNum] & (1 << (num % bitNum))) != 0
}
