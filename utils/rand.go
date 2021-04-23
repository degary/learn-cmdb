package utils

import (
	"math/rand"
	"time"
)

//随机生成字符串
func RandString(length int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	count := len(letters)
	chars := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		chars[i] = letters[rand.Intn(count)]
	}

	return string(chars)
}

func init() {
	rand.Seed(time.Now().UnixNano())

}
