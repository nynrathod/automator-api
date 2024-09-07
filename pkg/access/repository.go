package access

import (
	"context"
	"fmt"
	"github.com/nynrathod/automator-api/pkg/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Repository interface {
	ToggleAccess(data *entities.Access) (any, error)

	SharedUser(data *entities.Access) (any, error)
	SharedAccess(data *entities.Access) (any, error)
}

type repository struct {
	Collection *mongo.Collection
}

func NewRepo(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}

func (r *repository) ToggleAccess(data *entities.Access) (any, error) {
	collection := r.Collection

	filter := bson.M{
		"appOwner":   data.AppOwner,
		"sharedWith": data.SharedWith,
		"app":        data.App,
	}

	// Find the document
	var foundDocument bson.M
	err := collection.FindOne(context.Background(), filter).Decode(&foundDocument)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			storageData := &entities.Access{
				AppOwner:   data.AppOwner,
				SharedWith: data.SharedWith,
				App:        data.App,
				Access:     data.Access,
				TimeStamp:  time.Now().Format(time.RFC3339Nano),
			}

			_, err := collection.InsertOne(context.Background(), storageData)
			if err != nil {
				//fmt.Println("Error inserting document:", err)
				return nil, err
			}

			return storageData, nil
		} else {
			//fmt.Println("Error finding document:", err)
			return nil, err
		}
	}

	updatedAccess := data.Access
	update := bson.M{
		"$set": bson.M{
			"access": updatedAccess,
		},
	}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		//fmt.Println("Error updating document:", err)
		return nil, err
	}

	if result.ModifiedCount == 0 {
		//fmt.Println("No document updated, possibly no change needed.")
	}

	// Optionally return the updated document or some confirmation
	return bson.M{
		"appOwner":   data.AppOwner,
		"sharedWith": data.SharedWith,
		"app":        data.App,
		"access":     updatedAccess,
	}, nil
}

func (r *repository) SharedUser(data *entities.Access) (any, error) {
	var accessResults []entities.Access
	filter := bson.D{
		{"appOwner", data.AppOwner},
		{"app", data.App},
	}

	cursor, errIfn := r.Collection.Find(context.Background(), filter)
	if errIfn != nil {
		//fmt.Println("errIfn", errIfn)
		return nil, errIfn
	}

	if errIfn = cursor.All(context.TODO(), &accessResults); errIfn != nil {
		//fmt.Println("errCur", errIfn)
		return nil, errIfn
	}

	//fmt.Println("findall", accessResults)

	var userResults []entities.User
	for _, access := range accessResults {
		var user entities.User

		objectID, err := primitive.ObjectIDFromHex(access.SharedWith)
		if err != nil {
			fmt.Println("Invalid ObjectID format:", access.SharedWith)
			continue
		}
		userFilter := bson.D{{"_id", objectID}}

		projection := bson.D{
			{"_id", 1},
			{"email", 1},
			{"firstName", 1},
			{"lastName", 1},
		}

		uErr := r.Collection.Database().Collection("users").FindOne(
			context.Background(),
			userFilter,
			options.FindOne().SetProjection(projection),
		).Decode(&user)

		if uErr != nil {
			//fmt.Println("usererrr", uErr)
			if uErr == mongo.ErrNoDocuments {
				//fmt.Println("No document found for filter:", userFilter)
				continue
			}
			return nil, uErr
		}

		user.Access = access.Access

		userResults = append(userResults, user)
	}

	//fmt.Println("userDetails", userResults)

	return userResults, nil
}

func (r *repository) SharedAccess(data *entities.Access) (any, error) {
	var accessResults []entities.Access
	filter := bson.D{
		{"sharedWith", data.SharedWith},
	}

	cursor, errIfn := r.Collection.Find(context.Background(), filter)
	if errIfn != nil {
		//fmt.Println("errIfn", errIfn)
		return nil, errIfn
	}

	if errIfn = cursor.All(context.TODO(), &accessResults); errIfn != nil {
		//fmt.Println("errCur", errIfn)
		return nil, errIfn
	}

	//fmt.Println("findall", accessResults)

	return accessResults, nil
}
