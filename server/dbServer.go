package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongodb struct {
	client *mongo.Client
	dbName string
}

func NewDB() *Mongodb {
	return &Mongodb{}
}

func (m *Mongodb) Connect(ip string, port string, name string) error {
	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprint("mongodb://", ip, ":", port)))
	if err != nil {
		return err
	} else {
		m.client = client
		m.dbName = name
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = m.client.Connect(ctx)

	if err != nil {
		return err
	}

	// Check the connection
	err = m.client.Ping(context.TODO(), nil)

	if err != nil {
		return err
	}

	return nil

}

func (m *Mongodb) Verify(ctx context.Context, col string, username string, password string) (bool, interface{}) {
	collection := m.client.Database(m.dbName).Collection(col)

	query := collection.FindOne(ctx, bson.M{"username": username})
	var dec bson.M
	err := query.Decode(&dec)
	if err != nil {
		log.Println("err decoding query")
		return false, nil
	}

	i := dec["identity"].(string)
	p := dec["password"].(string)
	if p != password {
		return false, nil
	} else {
		return true, i
	}
}

func (m *Mongodb) Add(ctx context.Context, col string, data interface{}) (interface{}, error) {
	collection := m.client.Database(m.dbName).Collection(col)

	result, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (m *Mongodb) DeleteOne(ctx context.Context, col string, key string, val interface{}) (*mongo.DeleteResult, error) {
	collection := m.client.Database(m.dbName).Collection(col)

	filter := bson.M{key: val}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Println("delete err")
		return nil, err
	}
	return result, nil
}

func (m *Mongodb) Update(ctx context.Context, col string, key string, val interface{}, data interface{}) *mongo.SingleResult {
	collection := m.client.Database(m.dbName).Collection(col)

	q := bson.M{key: val}
	set := bson.D{{"$set", data}}
	ops := options.FindOneAndUpdate().SetUpsert(true)
	cur := collection.FindOneAndUpdate(ctx, q, set, ops)

	return cur
}

type Post struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Tag       string             `json:"tag" bson:"tag"`
	Title     string             `json:"title" bson:"title"`
	User      string             `json:"user" bson:"user"`
	Factory   string             `json:"factory" bson:"factory"`
	Market    string             `json:"market" bson:"market"`
	Date      string             `json:"date" bson:"date"`
	Amount    string             `json:"amount" bson:"amount"`
	Progress  string             `json:"progress" bson:"progress"`
	Paperwork string             `json:"paperwork" bson:"paperwork"`
	BCData    string             `json:"bcdata" bson:"bcdata"`
}

type BCdata struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Factory  string             `json:"factory" bson:"factory"`
	Date     string             `json:"date" bson:"date"`
	Chain    string             `json:"chain" bson:"chain"`
	BlockNum int                `json:"blocknum" bson:"blocknum"`
	Hash     string             `json:"hash" bson:"hash"`
	Image    string             `json:"image" bson:"image"`
}

func (m *Mongodb) QueryOne(ctx context.Context, col string, key string, val interface{}) *mongo.SingleResult {
	collection := m.client.Database(m.dbName).Collection(col)

	cur := collection.FindOne(ctx, bson.D{{key, val}})

	return cur
}

func (m *Mongodb) Query(ctx context.Context, col string, key string, val interface{}) (*mongo.Cursor, error) {
	collection := m.client.Database(m.dbName).Collection(col)
	cur, err := collection.Find(ctx, bson.D{{key, val}})

	if err != nil {
		fmt.Println("find err")
		return nil, err
	}

	return cur, nil

}

func (m *Mongodb) QueryAll(ctx context.Context, col string) (*mongo.Cursor, error) {
	collection := m.client.Database(m.dbName).Collection(col)
	cur, err := collection.Find(ctx, bson.D{})

	if err != nil {
		fmt.Println("find err")
		return nil, err
	}

	return cur, nil

}

func (m *Mongodb) hey() {
	collection := m.client.Database("testing").Collection("members")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// res, err := collection.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159})
	// id := res.InsertedID
	// fmt.Println(id)

	fmt.Println("Connected to MongoDB!")
	// ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
	cur := collection.FindOne(ctx, bson.M{"name": "brandon"})
	var p interface{}
	err := cur.Decode(&p)
	if err != nil {
		log.Fatal(err)
	}
	j, _ := json.Marshal(p)
	fmt.Println(string(j))
	_ = json.Unmarshal(j, &p)

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer cur.Close(ctx)
	/// for cur.Next(ctx) {
	/// 	var result bson.M
	/// 	err := cur.Decode(&result)
	/// 	if err != nil {
	/// 		log.Fatal(err)
	/// 	}
	/// 	// do something with result....
	/// 	fmt.Println(result)
	/// }
	// if err := cur.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	// Rest of the code will go here
	fmt.Println(p)
}

// func main() {
// 	a := &Mongodb{}
// 	a.Connect("localhost", "27017", "testing")
// 	// a.hey()
// 	// g := map[string]string{"hello": "hi", "asdf": "2345", "1234": "5678"}
// 	// a.Update("test", "hello", "hi", g)
// 	c, d := a.QueryAll("test")
// 	if d != nil {
// 		panic(d)
// 	}
// 	fmt.Println(c)
//
// }
