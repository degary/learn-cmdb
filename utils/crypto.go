package utils

import (
	"crypto/md5"
	"fmt"
	"strings"
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

//分割字符串,返回 盐和string
func SplitMd5Salt(text string) (string, string) {
	nodes := strings.Split(text, ":")
	if len(nodes) >= 2 {
		//返回盐和加密的密码
		return nodes[0], nodes[1]
	} else {
		return "", nodes[0]
	}
}
