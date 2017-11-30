package utils

// CheckError 检验error是否为nil，若不为nil，则panic
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
