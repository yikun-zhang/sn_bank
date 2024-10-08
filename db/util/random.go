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

func RandomInt(min, max int) int64 {
    return int64(min + rand.Intn(max-min+1))
}


func RandomString(n int) string {
	var sb strings.Builder

	for i:=0 ; i<n ;i++ {
		c := alphabet[rand.Intn(n)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0,1000)
}

func RandomCurrency() string {
	currencies := []string{"EUR","USD","CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}