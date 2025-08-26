package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func RandomNumber(input int) string {
	rand.Seed(time.Now().UnixNano())
	randomNumber := ""
	for i := 0; i < input; i++ {
		randomDigit := rand.Intn(10)
		randomNumber += strconv.Itoa(randomDigit)
	}
	return randomNumber
}
