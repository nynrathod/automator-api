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
	"time"
)

var ws *socketio.Websocket

type Repository interface {
	StoreSms(smsData *entities.Sms) (*entities.Sms, error)
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

func (r *mongoRepository) StoreSms(smsData *entities.Sms) (*entities.Sms, error) {
	//fmt.Println("myEmail", smsData.Email)
	//fmt.Println("myUserId", smsData.UserId)
	//fmt.Println("myApp", smsData.App)
	//fmt.Println("myOtp", smsData.Otp)
	//fmt.Println("myTimestamp", smsData.TimeStamp)

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

	uuid, ok := ctx.Value("UUID").(string)
	if !ok {
		return nil, fmt.Errorf("UUID not found in context")
	}

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

	userIDStr := result.ID.Hex()

	fmt.Println("founduser ", requestData.AppName)

	var resultFromOtherCollection entities.Sms // Replace with appropriate type
	otherCollectionFilter := bson.D{
		{"userId", userIDStr},
		{"app", requestData.AppName},
	}

	time.Sleep(5 * time.Second)

	const maxRetries = 3
	const retryDelay = 3 * time.Second
	var errSms error
	//var errDelete error

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

	for i := 0; i < maxRetries; i++ {
		errSms = r.Collection.FindOne(context.Background(), otherCollectionFilter).Decode(&resultFromOtherCollection)
		if errSms == nil {
			break
		}

		fmt.Printf("Error finding document, retrying in %v: %v\n", retryDelay, errSms)
		time.Sleep(retryDelay)
	}

	if errSms != nil {
		fmt.Println("Error finding document in other collection:", errSms)
		return nil, errSms
	}

	jsonData, errMarshal := json.Marshal(resultFromOtherCollection)
	if errMarshal != nil {
		fmt.Printf("Error marshalling resultFromOtherCollection: %v\n", errMarshal)
		return nil, errMarshal
	}

	fmt.Printf("Document found in other collection: %+v\n", string(jsonData))

	errEmit := ws.EmitTo(uuid, jsonData, socketio.TextMessage)
	if errEmit != nil {
		fmt.Printf("Error emitting to WebSocket: %v\n", errEmit)
		return nil, errEmit
	}

	_, _ = r.Collection.DeleteOne(context.Background(), otherCollectionFilter)

	return nil, nil
}
