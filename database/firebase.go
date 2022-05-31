package database

import (
	"context"
	"sync"

	firebase "firebase.google.com/go"
	"github.com/fabrv/watchman-server/utils"
	"google.golang.org/api/option"
)

var firebaseLock = &sync.Mutex{}
var firebaseInstance *firebase.App

func FirebaseInstance() *firebase.App {
	if firebaseInstance == nil {
		firebaseLock.Lock()
		defer firebaseLock.Unlock()

		opt := option.WithCredentialsFile(utils.GetEnv("FIREBASE_CREDENTIALS_FILE", "./firebase-credentials.json"))

		if firebaseInstance == nil {
			fi, err := firebase.NewApp(context.Background(), nil, opt)
			if err != nil {
				panic(err)
			}

			firebaseInstance = fi
		}
	}
	return firebaseInstance
}
