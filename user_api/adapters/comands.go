package adapters

import (
	"github.com/obarbier/awesome-crypto/user_api/domain"
	"go.mongodb.org/mongo-driver/bson"
)

type createUserCommand struct {
	Collection string        `bson:"insert"`
	User       []domain.User `bson:"documents"`
}

type findUserCommand struct {
	Collection string      `bson:"find"`
	Filter     interface{} `bson:"filter"`
}
type updateStatement struct {
	Query   bson.D `bson:"q"`
	Updates bson.D `bson:"u"`
}

type updateUserCommand struct {
	Collection string            `bson:"update"`
	Updates    []updateStatement `bson:"updates"`
}
