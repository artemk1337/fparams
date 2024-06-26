package testdata

func invalidArgsAndResultsFuncA(a int, // want "the parameters and returns of the function \"invalidArgsAndResultsFuncA\" should start on a new line"
	b string) (c bool,
	d error) {
	return false, nil
}

func invalidArgsAndResultsFuncB(a int, b int, // want "the parameters and returns of the function \"invalidArgsAndResultsFuncB\" should start on a new line"
	c string) (
	d bool,
	e error) {
	return false, nil
}
