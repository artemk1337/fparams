package testdata

func invalidArgsAndResultsFuncA(a int, // want "the parameters and return values of the function \"invalidArgsAndResultsFuncA\" should be on separate lines"
	b string) (c bool,
	d error) {
	return false, nil
}

func invalidArgsAndResultsFuncB(a int, b int, // want "the parameters and return values of the function \"invalidArgsAndResultsFuncB\" should be on separate lines"
	c string) (
	d bool,
	e error) {
	return false, nil
}
