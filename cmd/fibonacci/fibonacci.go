package main

import "context"

func getNumber(ctx context.Context, n int) int {
	switch n {
	case 0, 1:
		return n
	default:
		return getNumber(ctx, n-1) + getNumber(ctx, n-2)
	}
}
