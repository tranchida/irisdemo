package main

import (
	"context"
	"log"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://localhost:27017"

type Car struct {
	Id    primitive.ObjectID `bson:"_id,omitempty"`
	Model string             `bson:"model,omitempty"`
	Color string             `bson:"color,omitempty"`
}

var carCollection *mongo.Collection

func main() {

	var err error
	carCollection, err = connectMongo()
	if err != nil {
		log.Fatal(err)
	}

	app := createApp()
	app.Listen(":8080")

}

func connectMongo() (*mongo.Collection, error) {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	car := client.Database("viacar").Collection("car")

	if err != nil {
		return nil, err
	}

	return car, nil
}

func createApp() *iris.Application {

	app := iris.New()

	//f, _ := os.Create("iris.log")
	//app.Logger().SetOutput(f)

	carParty := app.Party("/car")

	carParty.Get("/", listAllCars)

	carParty.Get("/{id}", findCarById)

	carParty.Post("/", addCar)

	return app
}

func listAllCars(ctx iris.Context) {

	cursor, err := carCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	var cars []Car

	if err = cursor.All(context.TODO(), &cars); err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	ctx.JSON(cars)
}

func findCarById(ctx iris.Context) {

	id := ctx.Params().Get("id")
	docId, err := primitive.ObjectIDFromHex(id)
	println(id)
	if err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	var car Car
	err = carCollection.FindOne(context.TODO(), bson.M{"_id": docId}).Decode(&car)
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}
	ctx.JSON(car)

}

func addCar(ctx iris.Context) {

	var car Car
	err := ctx.ReadBody(&car)
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	result, err := carCollection.InsertOne(context.TODO(), car)
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	ctx.Application().Logger().Log(golog.InfoLevel, result.InsertedID)

}
