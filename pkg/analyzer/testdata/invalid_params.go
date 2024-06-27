package testdata

func invalidArgsFuncA(a int, // want "the parameters of the function \"invalidArgsFuncA\" should be on separate lines"
	b string) {
	return
}

func invalidArgsFuncB(a, b int, // want "the parameters of the function \"invalidArgsFuncB\" should be on separate lines"
	c string) {
	return
}

func invalidArgsFuncC(a, // want "the parameters of the function \"invalidArgsFuncC\" should be on separate lines"
	b int,
	c string,
) {
	return
}

func invalidArgsFuncD( // want "the parameters of the function \"invalidArgsFuncD\" should be on separate lines"
	a, b int,
	c string,
) {
	return
}
