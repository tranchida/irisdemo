package main

import (
	"context"
	"log"
	"time"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// const uri = "mongodb://localhost:27017"
const uri = "mongodb+srv://giampa0l0:PkZ66tL$X9@cluster0.7xoftco.mongodb.net/?retryWrites=true&w=majority"

type Car struct {
	ID                   string `bson:"_id"`
	ListingURL           string `bson:"listing_url"`
	Name                 string `bson:"name"`
	Summary              string `bson:"summary"`
	Space                string `bson:"space"`
	Description          string `bson:"description"`
	NeighborhoodOverview string `bson:"neighborhood_overview"`
	Notes                string `bson:"notes"`
	Transit              string `bson:"transit"`
	Access               string `bson:"access"`
	Interaction          string `bson:"interaction"`
	HouseRules           string `bson:"house_rules"`
	PropertyType         string `bson:"property_type"`
	RoomType             string `bson:"room_type"`
	BedType              string `bson:"bed_type"`
	MinimumNights        string `bson:"minimum_nights"`
	MaximumNights        string `bson:"maximum_nights"`
	CancellationPolicy   string `bson:"cancellation_policy"`
	LastScraped          struct {
		Date time.Time `bson:"$date"`
	} `bson:"last_scraped"`
	CalendarLastScraped struct {
		Date time.Time `bson:"$date"`
	} `bson:"calendar_last_scraped"`
	FirstReview struct {
		Date time.Time `bson:"$date"`
	} `bson:"first_review"`
	LastReview struct {
		Date time.Time `bson:"$date"`
	} `bson:"last_review"`
	Accommodates    int      `bson:"accommodates"`
	Bedrooms        int      `bson:"bedrooms"`
	Beds            int      `bson:"beds"`
	NumberOfReviews int      `bson:"number_of_reviews"`
	Bathrooms       float64  `bson:"bathrooms"`
	Amenities       []string `bson:"amenities"`
	Price           struct {
		NumberDecimal string `bson:"$numberDecimal"`
	} `bson:"price"`
	WeeklyPrice struct {
		NumberDecimal string `bson:"$numberDecimal"`
	} `bson:"weekly_price"`
	MonthlyPrice struct {
		NumberDecimal string `bson:"$numberDecimal"`
	} `bson:"monthly_price"`
	CleaningFee struct {
		NumberDecimal string `bson:"$numberDecimal"`
	} `bson:"cleaning_fee"`
	ExtraPeople struct {
		NumberDecimal string `bson:"$numberDecimal"`
	} `bson:"extra_people"`
	GuestsIncluded struct {
		NumberDecimal string `bson:"$numberDecimal"`
	} `bson:"guests_included"`
	Images struct {
		ThumbnailURL string `bson:"thumbnail_url"`
		MediumURL    string `bson:"medium_url"`
		PictureURL   string `bson:"picture_url"`
		XlPictureURL string `bson:"xl_picture_url"`
	} `bson:"images"`
	Host struct {
		HostID                 string   `bson:"host_id"`
		HostURL                string   `bson:"host_url"`
		HostName               string   `bson:"host_name"`
		HostLocation           string   `bson:"host_location"`
		HostAbout              string   `bson:"host_about"`
		HostResponseTime       string   `bson:"host_response_time"`
		HostThumbnailURL       string   `bson:"host_thumbnail_url"`
		HostPictureURL         string   `bson:"host_picture_url"`
		HostNeighbourhood      string   `bson:"host_neighbourhood"`
		HostResponseRate       int      `bson:"host_response_rate"`
		HostIsSuperhost        bool     `bson:"host_is_superhost"`
		HostHasProfilePic      bool     `bson:"host_has_profile_pic"`
		HostIdentityVerified   bool     `bson:"host_identity_verified"`
		HostListingsCount      int      `bson:"host_listings_count"`
		HostTotalListingsCount int      `bson:"host_total_listings_count"`
		HostVerifications      []string `bson:"host_verifications"`
	} `bson:"host"`
	Address struct {
		Street         string `bson:"street"`
		Suburb         string `bson:"suburb"`
		GovernmentArea string `bson:"government_area"`
		Market         string `bson:"market"`
		Country        string `bson:"country"`
		CountryCode    string `bson:"country_code"`
		Location       struct {
			Type            string    `bson:"type"`
			Coordinates     []float64 `bson:"coordinates"`
			IsLocationExact bool      `bson:"is_location_exact"`
		} `bson:"location"`
	} `bson:"address"`
	Availability struct {
		Availability30  int `bson:"availability_30"`
		Availability60  int `bson:"availability_60"`
		Availability90  int `bson:"availability_90"`
		Availability365 int `bson:"availability_365"`
	} `bson:"availability"`
	ReviewScores struct {
		ReviewScoresAccuracy      int `bson:"review_scores_accuracy"`
		ReviewScoresCleanliness   int `bson:"review_scores_cleanliness"`
		ReviewScoresCheckin       int `bson:"review_scores_checkin"`
		ReviewScoresCommunication int `bson:"review_scores_communication"`
		ReviewScoresLocation      int `bson:"review_scores_location"`
		ReviewScoresValue         int `bson:"review_scores_value"`
		ReviewScoresRating        int `bson:"review_scores_rating"`
	} `bson:"review_scores"`
	Reviews []struct {
		ID   string `bson:"_id"`
		Date struct {
			Date time.Time `bson:"$date"`
		} `bson:"date"`
		ListingID    string `bson:"listing_id"`
		ReviewerID   string `bson:"reviewer_id"`
		ReviewerName string `bson:"reviewer_name"`
		Comments     string `bson:"comments"`
	} `bson:"reviews"`
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
