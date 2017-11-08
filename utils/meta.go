package utils

import (
	"math/rand"
	"regexp"
	"time"
)

// TrimPathName //seg1////seg2 --> /seg1/seg2
func TrimPathName(s string) string {
	if len(s) == 0 {
		return ""
	}
	var reg = regexp.MustCompile(`/+`)
	return reg.ReplaceAllString(s, "/")
}

// TrimSpace 去除字符串首尾的空格符 \t \n 等
func TrimSpace(s string) string {
	if len(s) == 0 {
		return ""
	}
	var reg = regexp.MustCompile(`(^\s+)|(\s+$)`)
	return reg.ReplaceAllString(s, "")
}

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// RandStringBytesMaskImprSrc 生成随机n位长度的token
func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}
