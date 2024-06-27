package testdata

func invalidResultsFuncA() (a bool, // want "the parameters and returns of the function \"invalidResultsFuncA\" should start on a new line"
	b error) {
	return false, nil
}

func invalidResultsFuncB() ( // want "the parameters and returns of the function \"invalidResultsFuncB\" should start on a new line"
	a bool,
	b error) {
	return false, nil
}

func invalidResultsFuncC() ( // want "the parameters and returns of the function \"invalidResultsFuncC\" should start on a new line"
	a bool, b error) {
	return false, nil
}

func invalidResultsFuncD() ( // want "the parameters and returns of the function \"invalidResultsFuncD\" should start on a new line"
	a, b bool,
	c error) {
	return false, false, nil
}

func invalidResultsFuncE() (bool, bool, // want "the parameters and returns of the function \"invalidResultsFuncE\" should start on a new line"
	error) {
	return false, false, nil
}

func invalidResultsFuncF() ( // want "the parameters and returns of the function \"invalidResultsFuncF\" should start on a new line"
	bool,
	bool,
	error) {
	return false, false, nil
}
