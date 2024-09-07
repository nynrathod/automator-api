package sms

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nynrathod/automator-api/pkg/entities"
	ser "github.com/nynrathod/automator-api/pkg/sms"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

func StoreSms(smsData *entities.Sms, db *mongo.Database, uuid string) (*entities.Sms, error) {

	repository := ser.NewMongoRepository(db, "sms_list")
	service := ser.NewService(repository)
	_, storeErr := service.StoreSms(smsData, uuid)
	if storeErr != nil {
		//fmt.Println("errstore", storeErr)
		return nil, fiber.NewError(fiber.StatusBadRequest)
	}
	return nil, nil
}

func ProcessOtpRequest(requestData *entities.OtpRequest, db *mongo.Database, uuid context.Context) (*entities.Sms,
	error) {
	repository := ser.NewMongoRepository(db, "sms_list")
	service := ser.NewService(repository)
	_, _ = service.ProcessOtpRequest(requestData, uuid)
	return nil, nil
}
