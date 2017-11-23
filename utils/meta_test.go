package utils

import "testing"
import "regexp"
import "strconv"
import "encoding/json"

func TestTrimPathName(t *testing.T) {
	var testURL = "//seg1////seg2"
	var res = TrimPathName(testURL)
	if res != "/seg1/seg2" {
		t.Errorf("trim url %v", testURL)
		t.Errorf("expected %v", "/seg1/seg2")
		t.Errorf("     got %v", res)
	}

	res = TrimPathName("")
	if res != "" {
		t.Errorf("trim url %v", "")
		t.Errorf("expected %v", "")
		t.Errorf("     got %v", res)
	}
}

func TestTrimSpace(t *testing.T) {
	var testStr = "  	data  		  "
	var res = TrimSpace(testStr)
	if res != "data" {
		t.Errorf("trim space %v", testStr)
		t.Errorf("expected %v", "data")
		t.Errorf("     got %v", res)
	}

	res = TrimSpace("")
	if res != "" {
		t.Errorf("trim space %v", "")
		t.Errorf("  expected %v", "")
		t.Errorf("       got %v", res)
	}
}

func TestRandStringBytesMaskImprSrc(t *testing.T) {
	var length = 32
	var token = RandStringBytesMaskImprSrc(length)
	if len(token) != length {
		t.Errorf("expected string length %v", length)
		t.Errorf("                   got %v", len(token))
	}
	var reg = regexp.MustCompile(`^[\da-zA-Z]{` + strconv.Itoa(length) + `}$`)
	var succ = reg.Match([]byte(token))
	if !succ {
		t.Errorf("unexpected token %v", token)
	}
}

func TestResolveStructToValues(t *testing.T) {
	var m FailReturn
	m.Code = 20000
	m.Message = "test ResolveStructToValues"
	var values = ResolveStructToValues(m)
	var code = values.Get("Code")
	if code != "20000" {
		t.Errorf("expected Code value %v", 20000)
		t.Errorf("                got %v", code)
	}

	var msg = values.Get("Message")
	if msg != "test ResolveStructToValues" {
		t.Errorf("expected Code value %v", "test ResolveStructToValues")
		t.Errorf("                got %v", msg)
	}
}

func TestInterfaceToStruct(t *testing.T) {
	var testJson = `{"Code":0,"Message":"succ"}`
	var m FailReturn
	var iData interface{}
	var err = json.Unmarshal([]byte(testJson), &iData)
	if err != nil {
		t.Error(err)
	}
	err = InterfaceToStruct(&iData, &m)
	if err != nil {
		t.Error(err)
	}
	if m.Code != 0 {
		t.Errorf("expected code %v", 0)
		t.Errorf("		    got %v", m.Code)
	}
	if m.Message != "succ" {
		t.Errorf("expected message %v", "succ")
		t.Errorf("		       got %v", m.Message)
	}
}

func TestMapToStruct(t *testing.T) {
	var mapData = make(map[string]interface{}, 0)
	var m FailReturn
	mapData["Code"] = 0
	mapData["Message"] = "succ"
	var err = MapToStruct(mapData, &m)
	if err != nil {
		t.Error(err)
	}
	if m.Code != 0 {
		t.Errorf("expected code %v", 0)
		t.Errorf("		    got %v", m.Code)
	}
	if m.Message != "succ" {
		t.Errorf("expected message %v", "succ")
		t.Errorf("		       got %v", m.Message)
	}
}
