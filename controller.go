package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection = db().Database("mydb").Collection("users")

type user struct {
	Id   string `json:"id" bson:"_id"`
	Name string `json:"name"`
	City string `json:"city"`
	Age  int    `json:"age"`
}

func createProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	var person user

	err := json.NewDecoder(r.Body).Decode(&person)

	if err != nil {
		fmt.Print(err)
	}

	person.Id = fmt.Sprintf(`%d`, time.Now().UnixMilli())

	insertResult, err := userCollection.InsertOne(context.TODO(), person)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(insertResult.InsertedID)
}

func getUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	var body user
	e := json.NewDecoder(r.Body).Decode(&body)

	vars := mux.Vars(r)

	if e != nil {
		fmt.Print(e)
	}
	var result primitive.M

	err := userCollection.FindOne(context.TODO(), bson.D{{"_id", vars[`id`]}}).Decode(&result)

	if err != nil {
		fmt.Print(err)

	}
	json.NewEncoder(w).Encode(result)

}

func updateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	type updateBody struct {
		Name string `json:"name"`
		City string `json:"city"`
	}

	var body updateBody
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {
		fmt.Print(e)
	}
	filter := bson.D{{"name", body.Name}}
	after := options.After
	returnOpt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	update := bson.D{{"$set", bson.D{{"city", body.City}}}}

	updateResult := userCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result primitive.M
	_ = updateResult.Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func deleteProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	params := mux.Vars(r)["id"]
	_id, err := primitive.ObjectIDFromHex(params)

	if err != nil {
		fmt.Print(err.Error)
	}
	opts := options.Delete().SetCollation(&options.Collation{})

	res, err := userCollection.DeleteOne(context.TODO(), bson.D{{"_id", _id}}, opts)

	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(res.DeletedCount)
}
