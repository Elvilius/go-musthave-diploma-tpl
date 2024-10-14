package utils

import (
	"strconv"
	"strings"
)

func CheckLuhn(cnn string) bool {
	sum := 0
	arrayStr := strings.Split(cnn, "")
	if len(arrayStr) == 0 {
		return false
	}
	parity := len(arrayStr) % 2
	for i := 0; i < len(arrayStr); i++ {
		digital, err := strconv.Atoi(arrayStr[i])
		if err != nil {
			return false
		}
		if i%2 == parity {
			digital = digital * 2
			if digital > 9 {
				digital = digital - 9
			}
		}
		sum = digital + sum

	}
	return sum%10 == 0
}
