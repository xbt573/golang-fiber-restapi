package database

import (
	"context"
	"fmt"
	"restapi/types"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	db     types.Database
	client *mongo.Client
)

func Connect(connectionString string, databaseName string) error {
	c, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(connectionString),
	)
	if err != nil {
		return err
	}

	err = c.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return err
	}

	client = c
	db = types.Database{
		Tasks: client.Database(databaseName).Collection("tasks"),
	}
	return nil
}

func Disconnect() error {
	if err := client.Disconnect(context.Background()); err != nil {
		return err
	}

	return nil
}

func InsertTask(task types.Task) (types.Task, error) {
	task.Id = uuid.New()

	_, err := db.Tasks.InsertOne(context.Background(), bson.D{
		{Key: "Id", Value: task.Id},
		{Key: "Name", Value: task.Name},
		{Key: "Description", Value: task.Description},
	})
	if err != nil {
		return task, err
	}

	return task, nil
}

func DeleteTask(id uuid.UUID) (types.Task, error) {
	res := db.Tasks.FindOneAndDelete(context.Background(), bson.D{
		{Key: "Id", Value: id},
	})

	task := types.Task{}
	err := res.Decode(&task)

	if err != nil {
		return task, err
	}

	return task, nil
}

func FindTask(id uuid.UUID) (types.Task, error) {
	res := db.Tasks.FindOne(context.Background(), bson.D{
		{Key: "Id", Value: id},
	})

	task := types.Task{}
	err := res.Decode(&task)
	if err != nil {
		return task, err
	}

	return task, nil
}

func AllTasks() ([]types.Task, error) {
	cur, err := db.Tasks.Find(context.Background(), bson.D{})
	if err != nil {
		return make([]types.Task, 0), err
	}

	tasks := make([]types.Task, 1)

	for cur.Next(context.Background()) {
		task := types.Task{}

		if err := cur.Decode(&task); err != nil {
			return make([]types.Task, 0), err
		}

		tasks = append(tasks, task)
	}

	if err := cur.Err(); err != nil {
		return make([]types.Task, 0), err
	}

	return tasks, nil
}

func UpdateTask(id uuid.UUID, task types.Task) (types.Task, error) {
	filter := bson.D{
		{Key: "Id", Value: id},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "Name", Value: task.Name},
			{Key: "Description", Value: task.Description},
		}},
	}

	_, err := db.Tasks.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println(err)
		return task, err
	}

	return task, nil
}
