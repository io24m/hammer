package crypto

import (
	"math/rand"
	"time"
)

const (
	r62 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

//RKey
func RKey(len int) (res []byte) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	res = make([]byte, 0)
	for i := 0; i < len; i++ {
		res = append(res, r62[r.Intn(61)])
	}
	return
}

//ReverseKey
func ReverseKey(key []byte) []byte {
	a := make([]byte, len(key), len(key))
	copy(a, key)
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return a
}
