package types

import (
	"math/rand"
	"time"
)

func Random(min, max int) int {
	if max <= min {
		return 0
	}

	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Min32(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

func Min64(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}
