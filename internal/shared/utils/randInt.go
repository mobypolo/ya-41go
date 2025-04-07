package utils

import "math/rand"

func RandInt(min, max int) int {
	return min + (max-min+1)*rand.Intn(max-min+1)
}
