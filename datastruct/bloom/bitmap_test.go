package bloom

import (
	"fmt"
	"log"
	"math"
	"testing"
)

func TestBitNum(t *testing.T) {
	log.Println(bitNum)
}

func TestInt(t *testing.T) {
	var num uint = math.MaxUint - 8
	fmt.Println(num)
	fmt.Println(int(num + 7))
}

func TestBitMap(t *testing.T) {
	bm := NewBitMap(24)
	fmt.Printf("%08b\n", bm.bits)
	bm.Set(24)
	bm.Set(1026)
	fmt.Printf("%08b\n", bm.bits)
	has := bm.Check(24)
	fmt.Println(has)
	bm.ReSet(24)
	fmt.Printf("%08b\n", bm.bits)
}
