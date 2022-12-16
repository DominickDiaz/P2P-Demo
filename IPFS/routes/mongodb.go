package routes

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
	"zendx.io/P2P-Drive/models"
	"zendx.io/P2P-Drive/sec"
)

// Model for  storting database client and is used for closing connection.
type MongoDb struct {
	Client *mongo.Client
}

// -------------------------- Establish DB Connection/Client --------------------------\\

func Connection() *MongoDb {
	url := "mongodb+srv://admin:DOM123@domsdb.agpuaxn.mongodb.net/?retryWrites=true&w=majority"
	//url := "mongodb+srv://" + os.Getenv("USER") + os.Getenv("PASS") + "@domsdb.agpuaxn.mongodb.net/?retryWrites=true&w=majority"
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(url).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return &MongoDb{Client: client}
}

//-------------------------- Register User into DB with client --------------------------\\

func (connection *MongoDb) DBregister(userInfo *models.RegisterRequest) {
	db := connection.Client.Database("P2P")
	coll := db.Collection("Users")
	userInfo.Token = uuid.New().String()
	MASTER := "e8f93939-7dd3-488e-9af3-c7093b3fdeaa"

	pass := sec.EncryptPass(&userInfo.Email, &MASTER)
	userInfo.LastLogin = time.Now().Format("09-07-2017")
	docs := bson.M{"_id": userInfo.Email, "Username": userInfo.Username, "UserPassword": pass, "Number": userInfo.Number, "Email": userInfo.Email,
		"Fname": userInfo.FirstName, "Lname": userInfo.LastName, "Token": userInfo.Token, "LastLogin": userInfo.LastLogin}
	result, err := coll.InsertOne(context.TODO(), docs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Inserted document with ID %v\n", result.InsertedID)
}

// -------------------------- Check If Email in DB --------------------------\\

func (connection *MongoDb) DBemailCheck(email string) string {
	var info models.RegisterRequest

	db := connection.Client.Database("P2P") //Set Database
	coll := db.Collection("Users")          //Set Collection
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	fmt.Println("Retreiving information...")
	filter := bson.M{"_id": email} //Set Filter
	//filter := bson.M{"Email": email}

	i, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return "Found"
	}
	for i.Next(context.TODO()) {
		var result models.RegisterRequest
		if err := i.Decode(&result); err != nil {
			panic(err)
		}
		info.Email = result.Email
	}
	fmt.Println("Successfully Retrieved")
	print(info.Email)
	if info.Email == "" {
		return "Not Found"
	} else {
		return "Found"
	}
}

// -------------------------- Get User Token from DB with client --------------------------\\

func (connection *MongoDb) Login(user *models.LoginRequest) string {
	db := connection.Client.Database("P2P")
	coll := db.Collection("Users")
	var result models.RegisterRequest
	MASTER := "e8f93939-7dd3-488e-9af3-c7093b3fdeaa"

	fmt.Println("Retreiving information...")
	fmt.Println(userLogin.Username)
	fmt.Println(userLogin.UserPassword)

	pass := sec.EncryptPass(&userLogin.Username, &MASTER)
	filter := bson.M{"Email": user.Username}

	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully Retrieved")
	if result.UserPassword == pass {

		filter := bson.M{"_id": result.Email}
		update := bson.M{"$set": bson.M{"LastLogin": time.Now().Format("09-07-2017")}}
		if _, err := coll.UpdateOne(context.Background(), filter, update); err != nil {
			return "er"
		}
		return result.Token
	} else {
		return "Incorrect Password"
	}
}

//-------------------------- Get User file info from DB with client --------------------------\\

func (connection *MongoDb) GetUserFiles(owner string) []models.AddResponse {
	var files []models.AddResponse
	db := connection.Client.Database("P2P")
	coll := db.Collection("User_Files")

	fmt.Println("Retreiving information...")

	filter := bson.M{"Owner": owner}
	fmt.Println("Retreiving information...")
	i, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return []models.AddResponse{}
	}

	fmt.Println("Retreiving information...")

	for i.Next(context.TODO()) {
		var result models.AddResponse

		if err := i.Decode(&result); err != nil {
			panic(err)
		}
		files = append(files, result)

	}

	fmt.Println("Successfully Retrieved")
	return files
}

//-------------------------- Upload File Data to DB with client --------------------------\\

func (connection *MongoDb) DBupload(file models.AddResponse) {

	db := connection.Client.Database("P2P")
	coll := db.Collection("User_Files")
	docs := bson.M{"Hash": file.Hash, "Name": file.Name,
		"Size": file.Size, "Link": file.Link, "Owner": file.Owner}
	result, err := coll.InsertOne(context.TODO(), docs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("inserted document with ID %v\n", result.InsertedID)
}

//-------------------------- Close Client --------------------------\\

func (connection *MongoDb) CloseClientDB() {
	if connection == nil {
		return
	}
	err := connection.Client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	// TODO optional you can log your closed MongoDB client
	fmt.Println("Connection to MongoDB closed.")
}

//-------------------------- Update User Tokens --------------------------\\

// CheckLastLogin checks the "LastLogin" field of all documents in a MongoDB collection and
// updates the "Token" field of any documents where the "LastLogin" date is more than 4 days old.
func (connection *MongoDb) CheckLastLogin() error {
	// Get a handle to the collection we want to query.
	db := connection.Client.Database("P2P")
	coll := db.Collection("Users")
	// Set up a pipeline to find documents where the "LastLogin" field is more than 4 days old.
	// We'll use the current time minus 4 days as a reference point for comparison.
	fourDaysAgo := time.Now().Format("09-07-2017") //.Add(-4 * 24 * time.Hour)
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"LastLogin": bson.M{
					"$lt": fourDaysAgo,
				},
			},
		},
	}

	// Execute the pipeline and get a cursor to the resulting documents.
	cursor, err := coll.Aggregate(context.Background(), pipeline)
	if err != nil {
		return fmt.Errorf("error finding documents with old LastLogin: %w", err)
	}

	// Iterate through the cursor and update the "Token" field of each document with a new UUID.
	for cursor.Next(context.Background()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return fmt.Errorf("error decoding document: %w", err)
		}

		// Generate a new UUID for the "Token" field.
		newToken := uuid.New()

		// Update the document with the new "Token" value.
		filter := bson.M{"_id": result["_id"]}
		update := bson.M{"$set": bson.M{"Token": newToken}}
		if _, err := coll.UpdateOne(context.Background(), filter, update); err != nil {
			return fmt.Errorf("error updating Token: %w", err)
		}
	}

	return nil
}
