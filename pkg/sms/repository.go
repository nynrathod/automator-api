package sms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/contrib/socketio"
	"github.com/nynrathod/automator-api/pkg/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var ws *socketio.Websocket

type Repository interface {
	StoreSms(smsData *entities.Sms, uuid string) (*entities.Sms, error)
	ProcessOtpRequest(requestData *entities.OtpRequest, uuid context.Context) (*entities.OtpRequest, error)
}

type mongoRepository struct {
	Collection      *mongo.Collection
	db              *mongo.Database
	UsersCollection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database, collectionName string) Repository {
	collection := db.Collection(collectionName)
	return &mongoRepository{
		Collection:      collection,
		UsersCollection: db.Collection("users"),
		db:              db,
	}
}

func (r *mongoRepository) StoreSms(smsData *entities.Sms, uuid string) (*entities.Sms, error) {
	//fmt.Println("myEmail", smsData.Email)
	//fmt.Println("myUserId", smsData.UserId)
	//fmt.Println("myApp", smsData.App)
	//fmt.Println("myOtp", smsData.Otp)
	//fmt.Println("myTimestamp", smsData.TimeStamp)
	fmt.Println("smsdata", smsData)
	//document := bson.D{
	//	{Key: "email", Value: smsData.Email},
	//	{Key: "user_id", Value: smsData.UserId},
	//	{Key: "app", Value: smsData.App},
	//	{Key: "otp", Value: smsData.Otp},
	//	{Key: "timestamp", Value: smsData.TimeStamp},
	//	{Key: "expiry", Value: smsData.Expiry},
	//}

	storageData := &entities.SmsForStorage{
		UserId:    smsData.UserId,
		App:       smsData.App,
		Otp:       smsData.Otp,
		Expiry:    smsData.Expiry,
		TimeStamp: smsData.TimeStamp,
	}

	fmt.Println("hello", smsData)
	res, err := r.Collection.InsertOne(context.Background(), storageData)
	if err != nil {
		fmt.Println("ASdasdvfsdfds", err)
		//return c.Status(fiber.StatusInternalServerError).SendString("Failed to insert document")
	}
	fmt.Printf("Insert result: %+v\n", res)
	aa := map[string]string{
		"type": "ack",
	}

	// Marshal the map to JSON
	jsonData, err := json.Marshal(aa)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		//return
	}
	ws.EmitTo(uuid, jsonData, socketio.TextMessage)
	return smsData, nil

	//return nil, nil
}

func (r *mongoRepository) DeleteDocument(filter bson.D) (*mongo.DeleteResult, error) {
	var deleteResult *mongo.DeleteResult
	var err error

	fmt.Println("deletingdocs")
	// Retry mechanism
	const maxRetries = 3
	const retryDelay = 3 * time.Second

	for i := 0; i < maxRetries; i++ {
		fmt.Println("insideretry")
		deleteResult, err = r.Collection.DeleteOne(context.Background(), filter)
		if err == nil {
			fmt.Println("errerdetelete", err)
			// If the operation succeeds, return the result
			return deleteResult, nil
		}

		// If there is an error, log it and retry after a delay
		fmt.Printf("Error deleting document, retrying in %v: %v\n", retryDelay, err)
		time.Sleep(retryDelay)
	}

	// If all retries fail, return the last error
	return nil, fmt.Errorf("failed to delete document after %d retries: %w", maxRetries, err)
}

func (r *mongoRepository) ProcessOtpRequest(requestData *entities.OtpRequest, ctx context.Context) (*entities.OtpRequest, error) {

	// ensure UUID found
	uuid, ok := ctx.Value("UUID").(string)
	if !ok {
		return nil, fmt.Errorf("UUID not found in context")
	}

	fmt.Println("requestData.AppName", requestData.AppName)

	// Check if user has access of app
	var userHasAccess entities.Sms // Replace with appropriate type
	userHasAccessFilter := bson.D{
		{"sharedWith", requestData.Requester},
		{"app", requestData.AppName},
	}
	accErr := r.Collection.Database().Collection("shared_access").FindOne(
		context.Background(),
		userHasAccessFilter,
		options.FindOne(),
	).Decode(&userHasAccess)
	if accErr != nil {
		fmt.Println("accErr", accErr)
		return nil, accErr
	}

	// Check if that user exits whose mobile number added
	var result entities.User
	filter := bson.D{{"mobileNumber", requestData.Recipient}}
	errUser := r.UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	if errUser != nil {
		fmt.Println("other err", errUser)
		if errors.Is(errUser, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("no user found with mobile number: %s", requestData.Recipient)
		}
		return nil, errUser
	}

	// filter to find sms in table
	userIDStr := result.ID.Hex()
	var resultFromOtherCollection entities.OtpResponse
	otherCollectionFilter := bson.D{
		{"userId", userIDStr},
		{"app", requestData.AppName},
	}

	//time.Sleep(5 * time.Second)

	// Ping user if still connected else delete first sms matching with this request
	// Retry because sms may need some time to receive and store in database
	const maxRetries = 3
	const retryDelay = 3 * time.Second
	var errSms error

	errWs := ws.EmitTo(uuid, []byte("your message"), socketio.TextMessage)
	if errWs != nil {

		for i := 0; i < maxRetries; i++ {
			res, errDelete := r.Collection.DeleteOne(context.Background(), otherCollectionFilter)
			if errDelete == nil {
				if res.DeletedCount > 0 {
					fmt.Println("Document deleted successfully")
					break
				} else {
					fmt.Println("No document found to delete, retrying...")
				}
			} else {
				fmt.Printf("Error deleting document, retrying in %v: %v\n", retryDelay, errDelete)
			}

			time.Sleep(retryDelay)
		}

		fmt.Printf("Error on ping %v\n", errWs)
		return nil, nil
	}

	// find  matching sms with retry
	for i := 0; i < maxRetries; i++ {
		errSms = r.Collection.FindOne(context.Background(), otherCollectionFilter).Decode(&resultFromOtherCollection)
		if errSms == nil {
			break
		}

		fmt.Printf("Error finding document, retrying in %v: %v\n", retryDelay, errSms)
		time.Sleep(retryDelay)
	}

	// No return on sms not found
	if errSms != nil {
		fmt.Println("Error finding document in other collection:", errSms)
		return nil, errSms
	}

	resultFromOtherCollection.Event = "OTP_RESPONSE"

	jsonData, errMarshal := json.Marshal(resultFromOtherCollection)
	if errMarshal != nil {
		fmt.Printf("Error marshalling resultFromOtherCollection: %v\n", errMarshal)
		return nil, errMarshal
	}

	fmt.Printf("Document found in other collection: %+v\n", string(jsonData))

	// Once sms found emit to user and delete this document
	errEmit := ws.EmitTo(uuid, jsonData, socketio.TextMessage)
	if errEmit != nil {
		fmt.Printf("Error emitting to WebSocket: %v\n", errEmit)
		return nil, errEmit
	}

	_, _ = r.Collection.DeleteOne(context.Background(), otherCollectionFilter)

	return nil, nil
}
