package utils

import (
	"testing"
)

func TestFriendlyLogger(t *testing.T) {
	var data struct {
		Name string
	}
	data.Name = "susan"
	Logger.Debug(data)
	Debug(data, data)

}
