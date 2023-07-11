package utils

import "math/rand"

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int) string {
	res := make([]byte, n)
	for i := range res {
		res[i] = letters[rand.Intn(len(letters))]
	}
	return string(res)
}

func RandomNumber(signed bool) int64 {
	// if signed true then return integer
	// else return only positive number
	return int64(rand.Intn(99999))
}
