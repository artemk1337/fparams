package testdata

func invalidArgsFuncA(a int, // want "the parameters and returns of the function \"invalidArgsFuncA\" should start on a new line"
	b string) {
	return
}

func invalidArgsFuncB(a, b int, // want "the parameters and returns of the function \"invalidArgsFuncB\" should start on a new line"
	c string) {
	return
}

func invalidArgsFuncC(a, // want "the parameters and returns of the function \"invalidArgsFuncC\" should start on a new line"
	b int,
	c string,
) {
	return
}

func invalidArgsFuncD( // want "the parameters and returns of the function \"invalidArgsFuncD\" should start on a new line"
	a, b int,
	c string,
) {
	return
}