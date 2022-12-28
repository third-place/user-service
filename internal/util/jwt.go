package util

import "os"

var JwtKey = []byte(os.Getenv("JWT_KEY"))
