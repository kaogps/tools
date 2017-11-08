package utils

import (
	"net/http"
	"testing"
)

type userModel struct {
	Account  string
	Password string
}

func TestResolveHTTPRespToInterface(t *testing.T) {
	var testURL = "https://smartinterface.narro.me/async/user/login"
	var data userModel
	data.Account = "testAccount"
	data.Password = "testPassword"
	res, err := http.PostForm(testURL, *ResolveStructToValues(data))
	if err != nil {
		t.Error(err)
	}
	var m FailReturn
	err = ResolveHTTPRespToInterface(res, &m)
	if err != nil {
		t.Error(err)
	}
}
