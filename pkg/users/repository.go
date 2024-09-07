package users

import (
	"context"
	"fmt"
	"github.com/nynrathod/automator-api/pkg/entities"
	UTL "github.com/nynrathod/automator-api/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	VerifyEmail(email string, isLogin bool) (*entities.User, error)
	Register(user *entities.User) (*entities.User, error)
	GetUser(email string) (*entities.User, error)
	AddUser(data *entities.User) (*entities.User, error)
	ListUser(data string) (*entities.User, error)
}

type repository struct {
	Collection *mongo.Collection
}

func NewRepo(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}

func (r *repository) VerifyEmail(email string, isLogin bool) (*entities.User, error) {

	filter := bson.M{"email": email}
	var result entities.User
	err := r.Collection.FindOne(context.Background(), filter).Decode(&result)
	if isLogin && err != nil {
		//fmt.Println("llerrr", err)
		return nil, mongo.ErrNoDocuments
	}
	if isLogin {
		//fmt.Println("yeslog", &result)
		return &result, nil
	}

	if err == nil {
		//fmt.Println("exitserr", err)
		return nil, err
	}

	return nil, mongo.ErrNoDocuments
}

func (r *repository) Register(user *entities.User) (*entities.User, error) {

	user.ID = primitive.NewObjectID()
	user.UserId = UTL.GenerateRandomID(14, 1)
	_, err := r.Collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (r *repository) GetUser(email string) (*entities.User, error) {

	filter := bson.M{"email": email}

	ctx := context.TODO()

	var user entities.User
	err := r.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		//log.Printf("Error finding user: %v\n", err)
		return nil, err
	}

	return &user, nil
}

func (r *repository) AddUser(data *entities.User) (*entities.User, error) {
	data.ID = primitive.NewObjectID()
	_, err := r.Collection.InsertOne(context.Background(), data)
	if err != nil {
		//fmt.Println("errr repo", err)
		//return nil, err
	}

	return data, nil
}

func (r *repository) ListUser(data string) (*entities.User, error) {
	var result entities.User
	filter := bson.D{{"email", data}}
	err := r.Collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		//fmt.Println("errr repo", err)
		//return nil, err
	}

	return &result, nil
}
