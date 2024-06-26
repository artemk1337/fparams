package testdata

func multiLineFuncA(
	a int,
	b string,
) (
	c bool,
	d error,
) {
	return false, nil
}

func multiLineFuncB(
	a int,
	b string,
) {
	return
}

func multiLineFuncC() (
	a bool,
	b error,
) {
	return false, nil
}
