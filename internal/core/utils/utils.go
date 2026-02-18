package utils

func Clamp(a, bottom, top int) int {

	if a < bottom {
		a = bottom
	}

	if a > top {
		a = top
	}

	return a
}

func ClampBottom(a, bottom int) int {

	if a < bottom {
		a = bottom
	}

	return a
}

func ClampTop(a, top int) int {

	if a > top {
		a = top
	}

	return a
}
