package entities

type Access struct {
	AppOwner   string `json:"appOwner" bson:"appOwner,omitempty" validate:"required"`
	SharedWith string `json:"sharedWith" bson:"sharedWith,omitempty" validate:"required"`
	App        string `json:"app" bson:"app,omitempty" validate:"required"`
	Access     string `json:"access" bson:"access,omitempty"`
	TimeStamp  string `json:"timeStamp" bson:"timeStamp,omitempty" validate:"required"`
}
