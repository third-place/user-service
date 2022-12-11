package util

import (
	"github.com/goombaio/namegenerator"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var NameGenerator namegenerator.Generator

func init() {
	seed := time.Now().UTC().UnixNano()
	NameGenerator = namegenerator.NewNameGenerator(seed)
}

func RandomUsername() string {
	return NameGenerator.Generate()
}

func RandomEmailAddress() string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	num := r.Intn(100000)
	return os.Getenv("EMAIL_PREFIX") + "+" + strconv.Itoa(num) + "@" + os.Getenv("EMAIL_DOMAIN")
}
