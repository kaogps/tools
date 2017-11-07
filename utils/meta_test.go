package utils

import "testing"

func TestTrimPathName(t *testing.T) {
	var testURL = "//seg1////seg2"
	var res = TrimPathName(testURL)
	if res != "/seg1/seg2" {
		t.Errorf("trim url %v", testURL)
		t.Errorf("expected %v", "/seg1/seg2")
		t.Errorf("     got %v", res)
	}
}
