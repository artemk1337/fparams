package testdata

func invalidResultsFuncA() (a bool, // want "the arguments of the function \"invalidResultsFuncA\" should start on a new line"
	b error) {
	return false, nil
}

func invalidResultsFuncB() ( // want "the arguments of the function \"invalidResultsFuncB\" should start on a new line"
	a bool,
	b error) {
	return false, nil
}

func invalidResultsFuncC() ( // want "the arguments of the function \"invalidResultsFuncC\" should start on a new line"
	a bool, b error) {
	return false, nil
}

func invalidResultsFuncD() ( // want "the arguments of the function \"invalidResultsFuncD\" should start on a new line"
	a, b bool,
	c error) {
	return false, false, nil
}
