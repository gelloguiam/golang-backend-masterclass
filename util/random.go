package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// generate random int
func RandomInt(min int64, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// generate random string
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// generate random money
func RandomMoney() int64 {
	return RandomInt(0, 100)
}

// generate random owner
func RandomOwner() string {
	return RandomString(15)
}

// generate random currency
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "PHP"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
