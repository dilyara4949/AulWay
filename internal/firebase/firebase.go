package firebase

import (
	"context"
	"firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
	"log"
)

func InitializeFirebase() (*auth.Client, error) {
	opt := option.WithCredentialsFile("firebase/aulway-9670a-firebase-adminsdk-fbsvc-f483a0bc6a.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v", err)
		return nil, err
	}

	// Get Auth client
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Firebase Auth client: %v", err)
		return nil, err
	}

	return authClient, nil
}
