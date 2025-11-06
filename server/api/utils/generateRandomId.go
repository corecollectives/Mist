package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomId() int64 {
	rand.Seed(time.Now().UnixNano())
	min := 100000                           
	max := 999999                           
	randomNum := rand.Intn(max-min+1) + min 
	return int64(randomNum)

}
