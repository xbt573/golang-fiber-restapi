package types

import "go.mongodb.org/mongo-driver/mongo"

type Database struct {
	Tasks *mongo.Collection
}
