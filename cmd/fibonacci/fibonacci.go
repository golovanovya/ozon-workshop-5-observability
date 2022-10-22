package main

func getNumber(n int) int {
	switch n {
	case 0, 1:
		return n
	default:
		return getNumber(n-1) + getNumber(n-2)
	}
}
