package utils

import (
	"strconv"
)

func AtoI(text string) int {
	intItem, err := strconv.Atoi(text)
	if err != nil {
		return 0
	}
	return intItem
}

func PointerBoolToBool(input *bool) bool {
	if input == nil {
		return false
	}
	return *input
}
