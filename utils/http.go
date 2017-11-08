package utils

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
)

// PostMultipartForm 使用http post方式提交文件请求
func PostMultipartForm(url string, value map[string]string, file map[string]io.Reader) (*http.Response, error) {
	var buf = new(bytes.Buffer)
	var w = multipart.NewWriter(buf)
	// 写参数
	for k, v := range value {
		w.WriteField(k, v)
	}
	// 写文件
	for k, v := range file {
		var fw, err = w.CreateFormFile(k, k)
		if err != nil {
			Logger.Error(err)
			return nil, err
		}
		_, err = io.Copy(fw, v)
		if err != nil {
			Logger.Error(err)
			return nil, err
		}
	}
	w.Close()
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	return res, nil
}
