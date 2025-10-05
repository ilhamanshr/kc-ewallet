package strhelper

import (
	"math/rand"
	"strconv"
	"time"
)

// RandomStrNumGenerator random string number generator
// with input length of desired string
func RandomStrNumGenerator(n int) (res string) {
	for len(res) < n {
		seed := rand.NewSource(time.Now().UnixNano())
		gen := rand.New(seed)
		i := gen.Intn(10)
		res += strconv.Itoa(i)
	}
	return
}
