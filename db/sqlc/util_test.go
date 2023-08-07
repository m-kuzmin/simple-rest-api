package sqlc_test

import (
	"fmt"
	"math/rand"
)

func RandomString(length uint) string {
	runes := make([]rune, length)

	for i := uint(0); i < length; i++ {
		runes[i] = rune(RandomNumber(int('a'), int('z')))
	}

	return string(runes)
}

func RandomNumber(min, max int) int {
	if min >= max {
		panic(fmt.Sprintf("min (%d) should be less than max (%d)", min, max))
	}

	return min + rand.Intn(max-min)
}
