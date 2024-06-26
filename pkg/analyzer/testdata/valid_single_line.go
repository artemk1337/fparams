package testdata

func singleLineFuncA(a int, b string) (c bool, d error) {
	return false, nil
}

func singleLineFuncB(a int) (b bool, c error) {
	return false, nil
}

func singleLineFuncC(a int, b string) (c error) {
	return nil
}

func singleLineFuncD(int, string) error {
	return nil
}

func singleLineFuncE(_ int, _ string) error {
	return nil
}
