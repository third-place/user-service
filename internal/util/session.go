package util

import (
	"github.com/alexedwards/scs/v2"
	"time"
)

var SessionManager *scs.SessionManager

func init() {
	SessionManager = scs.New()
	SessionManager.Lifetime = 24 * time.Hour
}
