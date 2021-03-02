package utils

import (
	"crypto/md5"
	"fmt"
)

//生成一个加盐的字符串
func Md5Salt(text, salt string) string {
	if salt == "" {
		salt = RandString(8)
	}

	bytes := []byte(fmt.Sprintf("%s:%s", salt, text))
	sum := md5.Sum(bytes)
	return fmt.Sprintf("%s:%x", salt, sum)
}
