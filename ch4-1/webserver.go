package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"html/template"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)
type UserModel struct {
	ID primitive.ObjectID	`bson:"_id,omitempty`
	Username string			`bson:"username,omitempty`
	Password string			`bson:"password,omitempty`
}
type User struct {
	Username string
	Authenticated bool
}

type ViewModel struct {
	View string
	User User
}

var store sessions.CookieStore

func getUser(s *sessions.Session) User {
	val := s.Values["user"]
	var user = User{}
	user, ok := val.(User)
	if !ok {
		return User{ Authenticated: false}
	}
	return user
}

func ToPng(b []byte) ([]byte, error){
	switch http.DetectContentType(b){
	case "image/png":
		return b, nil
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	return nil, fmt.Errorf("png conversion error")
}

func UploadFile(w http.ResponseWriter, r *http.Request){
	file, handler, err := r.FormFile("file")

	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Println(filepath.Ext(handler.Filename))

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		fmt.Println(err)
	}

	p, err := ToPng(buf.Bytes())

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w,err.Error(), http.StatusInternalServerError)
		return
	}
	user := getUser(session)

	_ = os.Mkdir("./uploads/" + user.Username, 0777)
	ioutil.WriteFile("./uploads/"+user.Username+"/profile"+"png",p, 0666)

	if err != nil {
		panic(err)
	}
	http.Redirect(w,r,"/", http.StatusFound)
}

func UploadForm(w http.ResponseWriter, r *http.Request){

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w,err.Error(), http.StatusInternalServerError)
		return
	}
	user := getUser(session)

	view := ViewModel{"Upload", user}

	t, _ := template.ParseFiles("./static/home.html","./static/upload.html", "./static/nav.html")

	err = t.Execute(w, view)

	if err != nil{
		fmt.Println(err)
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password),14)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	bytes, err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return string(bytes), err
}

func NewUser(w http.ResponseWriter, r *http.Request){
	apiKey, _ := ioutil.ReadFile("mongo.txt")
	client, err := mongo.NewCient(options.Client().ApplyURI(string(apiKey)))

	if err != nil{
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Backgraoun(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("go-app-db")
	collection := db.Collection("go-app-collection")

	result, err := collection.CountDocuments(ctx, bson.M{"username": bson.D{Key:"$sq", Value: r.FormValue("username")}})

	if err != nil {
		log.Fatal(err)
	}

	if result != 0 {
		fmt.Println("User already exists in db")
	} else {
		hashedPassword, err := HashPassword(r.FormValue("password"))

		if err != nil {
			panic(err)
		}

		u := UserModel {
			Username: r.FormValue("username"),
			Password: hashedPassword,
		}

		insertResult, err := collection.InsertOne(ctx, u)
		if err != nil {
			panic(err)
		}
		fmt.Println(insertResult.InsertedID)
	}

	http.Redirect(w,r,"/login", http.StatusFound)

}

func NewUserForm(w http.ResponseWriter, r *http.Request){

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w,err.Error(), http.StatusInternalServerError)
		return
	}
	user := getUser(session)

	view := ViewModel{"NewUser", user}

	t, _ := template.ParseFiles("./static/home.html","./static/newuser.html", "./static/nav.html")

	err = t.Execute(w, view)

	if err != nil{
		fmt.Println(err)
	}
}

func Login(w http.ResponseWriter, r *http.Request){
	apiKey, _ := ioutil.ReadFile("mongo.txt")
	client, err := mongo.NewCient(options.Client().ApplyURI(string(apiKey)))

	if err != nil{
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Backgraoun(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("go-app-db")
	collection := db.Collection("go-app-collection")

	cursor := collection.FindOne(ctx, bson.M{"username":bson.D{key:"$eq", Value: r.FormValue("username")}})

	if err != nil {
		log.Fatal(err)
	}

	if cursor != nil {
		fmt.Println("user in db.")
		fmt.Println(cursor)

		var result UserModel
		if err = cursor.Decode(&result); err != nil {
			panic(err)
		}
		if CheckPasswordHash(r.FormValue("password"), result.Password){
			fmt.Println("Password matches")
			session, _ := store.Get(r, "session-name")
			session.Options.Path="/"
			user := &User {
				Username: r.FormValue("username"),
				Authenticated: true,
			}
			session.Values["user"] = user

			err = session.Save(r,w)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			fmt.Println("Wrong password.")
			http.Redirect(w,r, "/newuser", http.StatusFound)
			return
		}
	} else {
		fmt.Println("User not in db.")
	}
	http.Redirect(w,r,"/newuser", http.StatusFound)
}

func LoginForm(w http.ResponseWriter, r *http.Request){

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w,err.Error(), http.StatusInternalServerError)
		return
	}
	user := getUser(session)


	view := ViewModel{"Login", user}

	t, _ := template.ParseFiles("./static/home.html","./static/login.html", "./static/nav.html")

	err = t.Execute(w, view)

	if err != nil{
		fmt.Println(err)
	}
}

func Logout(w http.ResponseWriter, r *http.Request){

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w,err.Error(), http.StatusInternalServerError)
		return
	}
	//user := getUser(session)
	
	session.Values["user"] = User{}
	session.Options.MaxAge = -1

	err = session.Save(r,w)
	if err != nil {
		http.Error(w,err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)

}

func Home(w http.ResponseWriter, r *http.Request){

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w,err.Error(), http.StatusInternalServerError)
		return
	}
	user := getUser(session)


	view := ViewModel{"Main", user}

	t, _ := template.ParseFiles("./static/home.html","./static/main.html", "./static/nav.html")

	err = t.Execute(w, view)

	if err != nil{
		fmt.Println(err)
	}
}

func main() {

	// initalize the cookie store
	key, _ := ioutil.ReadFile("key.txt")
	store = sessions.NewCookieStore([]byte(key))

	//register the User struct to the store interface
	gob.Register(User{})

	r:= mux.NewRouter()

	dir := "/uploads"
	flag.StringVar(&dir, "dir", ".", "the directory to serve files from.")
	flag.Parse()
	r.PathPrefix("/uploads").Handler(http.FileServer(http.Dir(dir)))

	r.HandleFunc("/upload", UploadFile).Methods("POST")
	r.HandleFunc("/upload", UploadForm)
	r.HandleFunc("/newuser", NewUser).Methods("POST")
	r.HandleFunc("/newuser", NewUserForm)
	r.HandleFunc("/login", Login).Methods("POST")
	r.HandleFunc("/login", LoginForm)
	r.HandleFunc("/logout", Logout).Methods("POST")

	r.HandleFunc("/", Home)
	log.Fatal(http.ListenAndServe(":80",r))
	
}