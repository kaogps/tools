package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
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

// ResolveHTTPRespToInterface 解析http response到结构体
func ResolveHTTPRespToInterface(input *http.Response, output interface{}) error {
	if input == nil {
		return errors.New("http response not found")
	}
	var body, err = ioutil.ReadAll(input.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &output)
	return err
}

// ResolveStructToValues 解析结构体数据到url.Values中
// 如果data不是结构体，则返回nil
func ResolveStructToValues(data interface{}) *url.Values {
	var values url.Values = make(map[string][]string, 0)
	var d = newStructData(data)
	d.innerData()
	if d.v.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < d.c; i++ {
		values.Add(d.t.Field(i).Name, fmt.Sprintf("%v", d.v.Field(i).Interface()))
	}
	return &values
}

type structData struct {
	v reflect.Value
	t reflect.Type
	c int
}

func newStructData(data interface{}) structData {
	var model structData
	model.v = reflect.ValueOf(data)
	model.t = reflect.TypeOf(data)
	model.c = model.v.NumField()

	return model
}

func (s structData) innerData() {
	if s.v.Kind() == reflect.Ptr {
		s.v = s.v.Elem()
	}
	if s.t.Kind() == reflect.Ptr {
		s.t = s.t.Elem()
	}
}
