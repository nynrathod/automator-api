package services

import (
	"github.com/nynrathod/automator-api/pkg/users"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppServices struct {
	UserService users.Service
	//AppsService    apps.Service
	UserCollection *mongo.Collection
	//AppsCollection *mongo.Collection
}
