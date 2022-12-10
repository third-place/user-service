package util

import "math/rand"

var numbers = []rune("0123456789")
var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateCode() string {
	l := make([]rune, 3)
	n := make([]rune, 3)
	for i := range l {
		l[i] = letters[rand.Intn(len(letters))]
	}
	for i := range n {
		n[i] = numbers[rand.Intn(len(numbers))]
	}
	return string(l) + "-" + string(n)
}
