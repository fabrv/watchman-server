package database

import (
	"sync"

	"github.com/gofiber/fiber/v2/middleware/session"
)

var sessionLock = &sync.Mutex{}
var sessionInstance *session.Store

func SessionInstance() *session.Store {
	if sessionInstance == nil {
		sessionLock.Lock()
		defer sessionLock.Unlock()

		if sessionInstance == nil {
			sessionInstance = session.New(session.Config{
				Storage: RedisStorageInstance(),
			})
		}
	}
	return sessionInstance
}
