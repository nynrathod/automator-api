package services

import (
	"github.com/nynrathod/automator-api/api/router"
	"github.com/nynrathod/automator-api/pkg/access"
	"github.com/nynrathod/automator-api/pkg/users"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppServiceInitializer struct {
	db *mongo.Database
}

func NewAppServiceInitializer(db *mongo.Database) *AppServiceInitializer {
	return &AppServiceInitializer{db: db}
}

func (initializer *AppServiceInitializer) InitializeAppServices(db *mongo.Database) router.AppServices {
	userCollection := db.Collection("users")
	smsCollection := db.Collection("sms_list")

	accessCollection := db.Collection("shared_access")
	//appsCollection := db.Collection("apps")

	userRepo := users.NewRepo(userCollection)
	accessRepo := access.NewRepo(accessCollection)
	//smsRepo := sms.NewMongoRepository(smsCollection)
	//appsRepo := apps.NewRepo(appsCollection)

	userService := users.NewService(userRepo)
	accessService := access.NewService(accessRepo)
	//smsService := sms.NewService(smsRepo)
	//appsService := apps.NewService(appsRepo)

	return router.AppServices{
		UserService:   userService,
		AccessService: accessService,
		//SmsService:  smsService,
		//AppsService:    appsService,
		UserCollection: userCollection,
		SmsCollection:  smsCollection,

		AccessCollection: accessCollection,
		//AppsCollection: appsCollection,
	}
}
