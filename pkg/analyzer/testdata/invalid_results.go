package testdata

func invalidResultsFuncA() (a bool, // want "the return values of the function \"invalidResultsFuncA\" should be on separate lines"
	b error) {
	return false, nil
}

func invalidResultsFuncB() ( // want "the return values of the function \"invalidResultsFuncB\" should be on separate lines"
	a bool,
	b error) {
	return false, nil
}

func invalidResultsFuncC() ( // want "the return values of the function \"invalidResultsFuncC\" should be on separate lines"
	a bool, b error) {
	return false, nil
}

func invalidResultsFuncD() ( // want "the return values of the function \"invalidResultsFuncD\" should be on separate lines"
	a, b bool,
	c error) {
	return false, false, nil
}

func invalidResultsFuncE() (bool, bool, // want "the return values of the function \"invalidResultsFuncE\" should be on separate lines"
	error) {
	return false, false, nil
}

func invalidResultsFuncF() ( // want "the return values of the function \"invalidResultsFuncF\" should be on separate lines"
	bool,
	bool,
	error) {
	return false, false, nil
}
