package main

import (
	"fmt"
	"io"
	"os"
        "time"
        "context"
        "strconv"
	"math/rand"
        "mime/multipart"
	"net/http"
	"html/template"
	"path/filepath"
        "encoding/json"

        "github.com/gorilla/mux"
        "go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/mongo/options"
        "go.mongodb.org/mongo-driver/bson"
)

type Post struct {
    Name string `json:"name"`
    Expiration int64 `json:"expiration"`
    ExpirationType int `json:"expirationtype"` // -1 none 0 seconds 1 views
}

type ExpirationPayload struct {
    Type int32 `json:"type"`
    Value int64 `json:"value"`
}

type CustomFSHandler = func(w http.ResponseWriter, r *http.Request) (doDefaultFileServe bool)

// this should be variable... once all 5 length fill up move to 6
var HASH_LENGTH = 5
// is this good?
var tpl = template.Must(template.ParseFiles("index.html"))



func main() {
        // thread for cleaning database
        go cleanup()

        // setup and serve
        mux := buildMux()
	http.ListenAndServe(":8000", mux)
}

func cleanup() {
    for range time.Tick(time.Second * 3) {
        filter := bson.M{"expiration": bson.M{"$lt": time.Now().Unix()}, "expirationtype": bson.M{"$eq": 0}}
        col := GetClient().Collection("posts")
        cur, err := col.Find(context.TODO(), filter)
        if err != nil {
            fmt.Println(err)
            continue
        }
        for cur.Next(context.TODO()) {
            post := Post{}
            err := cur.Decode(&post)
            if err != nil {
                fmt.Println(err)
                continue
            }
            fmt.Println("deleting: " + post.Name)
            col.DeleteOne(context.TODO(), bson.M{"name": post.Name})
            // now the file itself
            os.Remove("f/" + post.Name);
        }
    }
}

func buildMux() *mux.Router{
	mux := mux.NewRouter()

        fs := http.FileServer(http.Dir("static"))
        fs2 := FileServerWithMongo(http.Dir("f"))
        mux.Handle("/static/{rest}", http.StripPrefix("/static/", fs))
        mux.Handle("/f/{rest}", http.StripPrefix("/f/", fs2))
        mux.HandleFunc("/", indexHandler)
        mux.HandleFunc("/post", createPostHandler)
        mux.HandleFunc("/expir/{name}", expirGetHandler)

	return mux
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func GetClient() *mongo.Database {
	client, err := mongo.Connect(
        context.Background(),
        options.Client().ApplyURI("mongodb://127.0.0.1/"),
    )

    if err != nil {
        fmt.Println(fmt.Errorf("Error: %v", err))
    }

    return client.Database("chadpaste")
}

func FileServerWithMongo(root http.FileSystem) http.Handler {
    fs := http.FileServer(root)

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // dec expiration by 1 if it is of tpye views
        col := GetClient().Collection("posts")
        filter := bson.M{"name": r.URL, "expirationtype": 1}
        update := bson.D{{"$inc",bson.D{{"expiration", -1}}}}
        _, err := col.UpdateOne(context.TODO(), filter, update)
        if err != nil {
            fmt.Println(err)
        }
        filter = bson.M{"name": r.URL}
        var post bson.M
        if err := col.FindOne(context.TODO(), filter).Decode(&post); err != nil {
            fmt.Println(fmt.Errorf("Error1: %v", err))
        }
        if post != nil {
            if (post["expirationtype"].(int32) == 1 && post["expiration"].(int64) < 0) {
                col.DeleteOne(context.TODO(), bson.M{"name": r.URL})
                os.Remove("f/" + r.URL.Path);
            }
        }

        fs.ServeHTTP(w, r)
    })
}


func expirGetHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {

        params := mux.Vars(r)
        name := params["name"]

        col := GetClient().Collection("posts")

        filter := bson.M{"name": name}
        fmt.Println(filter)
        var post bson.M
        if err := col.FindOne(context.TODO(), filter).Decode(&post); err != nil {
            fmt.Println(fmt.Errorf("Error1: %v", err))
        }
        if post == nil {
            response := ExpirationPayload{}
            response.Type = -2
            json.NewEncoder(w).Encode(response)
            return
        }
        response := ExpirationPayload{}
        response.Type = post["expirationtype"].(int32)
        response.Value = post["expiration"].(int64)
        json.NewEncoder(w).Encode(response)
    }
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20) // 32MB upload limit
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err)
			return
		}
		// TODO: maybe propogate error up to get a sense if it actually saved
                i, err := strconv.ParseInt(r.FormValue("expiration")[0:], 10, 64);
                var exp = -1
                expir_val := r.FormValue("expiration_type")
                fmt.Println(expir_val)
                if expir_val == "seconds" {
                    exp = 0
                } else if expir_val == "views" {
                    exp = 1
                }
		name := saveFile(file, header, i, exp)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(201)
		w.Write([]byte(name))
	} else if r.Method == "GET" {
		w.Write([]byte(fmt.Sprintf("get")))
	}
}

func saveFile(file multipart.File, header *multipart.FileHeader, exp int64, e_type int) string {
	defer file.Close()
	var name = genAndCheckNewURL() + filepath.Ext(header.Filename)
	// check collision
	if _, err := os.Stat("./f/" + name); err == nil {
		// path/to/whatever exists
	}
	f, err := os.OpenFile("./f/" + name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return "/error"
	}
	defer f.Close()
	io.Copy(f, file)
        post := Post{}
        post.Name = name
        if e_type == 0 {
            post.Expiration = expirationEpoch(exp)
        } else {
            post.Expiration = exp
        }
        post.ExpirationType = e_type
        col := GetClient().Collection("posts")
        col.InsertOne(context.TODO(), post)
	return name
}


func expirationEpoch(addedTime int64) int64 {
    now := time.Now().Unix()
    return now + addedTime
}

// calhoun.io inspo 
const urlAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234566789"

var seededRandomization *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomString() string {
    b := make([]byte, HASH_LENGTH)
    for i := range b {
        b[i] = urlAlphabet[seededRandomization.Intn(len(urlAlphabet))]
    }
    return string(b)
}

func genAndCheckNewURL() string {
        potString := randomString()

        // TODO check against database and try a few times....
        // potentially add logic to increase size ofstring if too many

        return potString
}
