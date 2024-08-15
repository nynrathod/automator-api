package users

import (
	"context"
	"fmt"
	"github.com/nynrathod/automator-api/pkg/entities"
	UTL "github.com/nynrathod/automator-api/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Repository interface {
	VerifyEmail(email string) (*entities.User, error)
	Register(user *entities.User) (*entities.User, error)
	GetUser(email string) (*entities.User, error)
}

type repository struct {
	Collection *mongo.Collection
}

// NewRepo is the single instance repo that is being created.
func NewRepo(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}

func (r *repository) VerifyEmail(email string) (*entities.User, error) {
	// nonce := make([]byte, 12) // Generate a unique nonce for each encryption.
	// if _, err := rand.Read(nonce); err != nil {
	// 	log.Fatal("Failed to generate nonce:", err)
	// }
	// hashedEmailToCheck, _ := utilities.EncryptData([]byte(email), nonce)
	filter := bson.M{"email": email}
	var result entities.User
	err := r.Collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	fmt.Println("myres", result)
	return &result, nil
}

func (r *repository) Register(user *entities.User) (*entities.User, error) {

	user.ID = primitive.NewObjectID()
	//user.CreatedAt = time.Now()
	// user.UpdatedAt = time.Now()
	user.UserId = UTL.GenerateRandomID(14, 1)
	_, err := r.Collection.InsertOne(context.Background(), user)

	if err != nil {
		return nil, err
	}
	return user, nil

}

func (r *repository) GetUser(email string) (*entities.User, error) {
	// Create a MongoDB client and establish a connection

	// Define a filter for the query
	filter := bson.M{"email": email}

	// Create a context
	ctx := context.TODO()

	// Perform the query
	var user entities.User
	err := r.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Handle case when no document is found
			return nil, fmt.Errorf("user not found")
		}
		log.Printf("Error finding user: %v\n", err)
		return nil, err
	}

	return &user, nil
}
