package sms

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/nynrathod/automator-api/pkg/entities"
	ser "github.com/nynrathod/automator-api/pkg/sms"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

func StoreSms(smsData *entities.Sms, db *mongo.Database, uuid string) (*entities.Sms, error) {

	//fmt.Println("myEmail", smsData.Email)
	//fmt.Println("myUserId", smsData.UserId)
	//fmt.Println("myApp", smsData.App)
	//fmt.Println("myOtp", smsData.Otp)
	//fmt.Println("myTimestamp", smsData.TimeStamp)

	//type InsertRequest struct {
	//	Name   string    `json:"name"`
	//	Expiry time.Time `json:"expiry"`
	//}
	//

	fmt.Println("\ncalling handler", smsData)

	fmt.Println("validall")

	// Instantiate the service
	repository := ser.NewMongoRepository(db, "sms_list")
	// Instantiate the service with the repository
	service := ser.NewService(repository)
	// Store the SMS using the serviceZ
	_, storeErr := service.StoreSms(smsData, uuid)
	if storeErr != nil {
		fmt.Println("errstore", storeErr)
		return nil, fiber.NewError(fiber.StatusBadRequest)
	}
	return nil, nil
}

func ProcessOtpRequest(requestData *entities.OtpRequest, db *mongo.Database, uuid context.Context) (*entities.Sms,
	error) {
	fmt.Println("requestData", requestData)

	repository := ser.NewMongoRepository(db, "sms_list")
	// Instantiate the service with the repository
	service := ser.NewService(repository)

	_, _ = service.ProcessOtpRequest(requestData, uuid)
	//filter := bson.D{{"mobile_number", requestData.Recipient}}
	//err := collection.FindOne(context.Background(), filter).Decode(&result)

	return nil, nil
}
