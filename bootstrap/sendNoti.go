package bootstrap

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

func SendNoti() {
	// Path to your service account key file
	opt := option.WithCredentialsFile("./serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	ctx := context.Background()
	_, err = app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	// This registration token comes from the client FCM SDKs.
	registrationToken := "fN93S8s5SkypwBtv3iW7Ys:APA91bEQWakxnJYuouZg0tadYYV7Eb8xY95ZQ9IZeEsrZ6o-BQ2Zu9mCDGNbL65pW4CielJkbOWZYS-3aPhx5JRzBT6hKJPd_d3nUP6DP7KNo_QC7We__d24qcvzxDI7nvto4ahYQBXp"

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Notification: &messaging.Notification{
			Title: "Your Notification Title",
			Body:  "Your notification body text.",
		},
		Token: registrationToken,
	}

	// Send a message to the device corresponding to the provided registration token.
	//response, err := client.Send(ctx, message)
	//if err != nil {
	//	log.Fatalln("myerrror", err)
	//}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", message)
}
