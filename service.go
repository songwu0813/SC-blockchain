package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"io/ioutil"
	"log"
	"mongo/server"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// type service struct {
type db interface {
	Add(string, interface{}) (interface{}, error)
}

type service struct {
	db   *server.Mongodb
	ip   string
	port string
}

type User struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Username   string             `json:"username" bson:"username"`
	Name       string             `json:"name,omitempty" bson:"name"`
	Department string             `json:"department,omitempty" bson:"department"`
	Email      string             `json:"email,omitempty" bson:"email"`
	Phone      string             `json:"phone,omitempty" bson:"phone"`
	Password   string             `json:"password" bson:"password"`
	Identity   string             `json:"identity" bson:"identity"`
}

type BCdataa struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Tag     string             `json:"tag" bson:"tag"`
	Name    string             `json:"name,omitempty" bson:"name"`
	Factory string             `json:"factory" bson:"factory"`
	Date    string             `json:"date" bson:"date"`
	Chain   string             `json:"chain" bson:"chain"`
	Hash    []string           `json:"hash" bson:"hash"`
	ImgHash []string           `json:"imghash" bson:"imghash"`
	Image   string             `json:"image" bson:"image"`
	Lat     string             `json:"lat" bson:"lat"`
	Long    string             `json:"long" bson:"long"`
	Dir     string             `json:"dir" bson:"dir"`
	FocLen  string             `json:"foclen" bson:"foclen"`
	DDDH    string             `json:"dddh" bson:"dddh"`
}

type BCdata struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Tag     string             `json:"tag" bson:"tag"`
	Name    string             `json:"name,omitempty" bson:"name"`
	Factory string             `json:"factory" bson:"factory"`
	Date    string             `json:"date" bson:"date"`
	Chain   string             `json:"chain" bson:"chain"`
	Hash    string             `json:"hash" bson:"hash"`
	ImgHash string             `json:"imghash" bson:"imghash"`
	Image   string             `json:"image" bson:"image"`
	Lat     string             `json:"lat" bson:"lat"`
	Long    string             `json:"long" bson:"long"`
	Dir     string             `json:"dir" bson:"dir"`
	FocLen  string             `json:"foclen" bson:"foclen"`
	DDDH    string             `json:"dddh" bson:"dddh"`
}

type Post struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Tag       string             `json:"tag" bson:"tag"`
	Title     string             `json:"title" bson:"title"`
	User      string             `json:"user" bson:"user"`
	Name      string             `json:"name" bson:"name"`
	Factory   string             `json:"factory" bson:"factory"`
	Market    string             `json:"market" bson:"market"`
	Date      string             `json:"date,omitempty" bson:"date"`
	Amount    int                `json:"amount,string" bson:"amount"`
	Progress  string             `json:"progress" bson:"progress"`
	Paperwork string             `json:"paperwork" bson:"paperwork"`
	// BCData    string             `json:"bcdata" bson:"bcdata"`
}

type Image struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Hash    []string           `json:"hash" bson:"hash"`
	ImgHash []string           `json:"imghash" bson:"imghash"`
	Img     string             `json:"img" bson"img"`
}

// type Image struct {
// 	ID      primitive.ObjectID `json:"id" bson:"_id"`
// 	Hash    string             `json:"hash" bson:"hash"`
// 	ImgHash string             `json:"imghash" bson:"imghash"`
// 	Img     string             `json:"img" bson"img"`
// }

func NewService(ip string, port string) *service {
	return &service{db: server.NewDB(), ip: ip, port: port}
}

func (s *service) Start(dbip string, dbport string, dbname string) error {
	err := s.db.Connect(dbip, dbport, dbname)
	if err != nil {
		return err
	}
	log.Println("db connected!")

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/user", s.newUser).Methods("POST")
	r.HandleFunc("/user", s.allUser).Methods("GET")
	r.HandleFunc("/user/{id}", s.user).Methods("GET")
	r.HandleFunc("/user/{id}", s.updateUser).Methods("PUT")
	r.HandleFunc("/user/{id}", s.deleteUser).Methods("DELETE")
	r.HandleFunc("/verifyuser", s.verifyUser).Methods("POST")
	r.HandleFunc("/verify", s.verifyUserAndReturnPost).Methods("POST")
	r.HandleFunc("/post", s.allPost).Methods("GET")
	r.HandleFunc("/post/{id}", s.post).Methods("GET")
	r.HandleFunc("/post/{id}", s.deletePost).Methods("DELETE")
	r.HandleFunc("/post", s.newPost).Methods("POST")
	r.HandleFunc("/post/{id}", s.updatePost).Methods("PUT")
	r.HandleFunc("/uploadfile", s.uploadFile).Methods("POST")
	r.HandleFunc("/uploadfile", s.allBcPost).Methods("GET")
	r.HandleFunc("/uploadfile/{id}", s.bcPost).Methods("GET")
	r.HandleFunc("/bcpost", s.allBcPost).Methods("GET")
	r.HandleFunc("/bcpost/{id}", s.bcPost).Methods("GET")
	r.HandleFunc("/bcpost", s.newBcPost).Methods("POST")

	r.HandleFunc("/verifyhash/{imghash}/{txhash}", s.verifyHash).Methods("GET")

	r.HandleFunc("/test", s.test).Methods("POST")

	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	r.PathPrefix("/file/").Handler(http.StripPrefix("/file/", http.FileServer(http.Dir("./file"))))

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	log.Println("server starting processing...")
	err = http.ListenAndServe(":"+s.port, handlers.CORS(headersOk, originsOk, methodsOk)(r))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *service) test(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	var buf interface{}
	err = json.Unmarshal(body, &buf)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(buf)
}

func (s *service) uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("uploadfile called")
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("image")
	id := r.FormValue("id")
	if err != nil {
		fmt.Println("form file image err")
		fmt.Println(err)
		return
	}
	fmt.Println("id:", id)
	for k, _ := range r.MultipartForm.File {
		fmt.Println(k)
	}
	defer file.Close()

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur := s.db.QueryOne(ctx, "posts", "tag", id)
	if err != nil {
		fmt.Println("err query one")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	k := &Post{}
	err = cur.Decode(k)
	if err != nil {
		fmt.Println("err decode cur")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	idd := k.ID.Hex()
	// _, err = s.db.Add(ctx, "testing", image)
	path := "images/" + id + "/"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}
	date := time.Now().Add(time.Hour * 8).Format("2006-01-02_150405")
	// date := time.Now().Add(time.Hour*8).Format("2006-01-02")
	fmt.Println("filename begin")
	filename := path + date + ".jpeg"
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	io.Copy(f, file)
	defer f.Close()
	fmt.Println("open file begin")
	g, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("python3 command begin")
	cmd := exec.Command("python3", "./agri/tmp.py", "add", ":"+filename)
	out, err := cmd.Output()
	if err != nil {
		log.Println("err output command")
		fmt.Println("bc err")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("python3 command complete")

	// hash := strings.TrimSpace(string(out))
	hash := strings.Split(string(out), "\n")
	_id, err := primitive.ObjectIDFromHex(idd)
	if err != nil {
		log.Println("err objectid from hex")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	img := "" + path
	// image := &Image{ID: _id, Hash: hash, Img: img}

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	// _, err = s.db.Add(ctx, "testing", image)

	bc, err := s.findBcPost(ctx, "bcposts", "_id", _id)
	if err != nil {
		log.Println("err find bc post")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	exif.RegisterParsers(mknote.All...)
	x, err := exif.Decode(g)
	if err != nil {
		fmt.Println("err decodeing exif " + filename)
		log.Fatal(err)
	}
	str := x.String()
	fmt.Println(str)
	lat, long, _ := x.LatLong()
	str_lat := strconv.FormatFloat(lat, 'f', 5, 64)
	str_long := strconv.FormatFloat(long, 'f', 5, 64)
	fmt.Println(str_lat, str_long)

	focal, _ := x.Get(exif.FocalLength)
	numer, denom, _ := focal.Rat2(0) // retrieve first (only) rat. value
	focallen := fmt.Sprintf("%.3f", float64(numer)/float64(denom))
	fmt.Println(focallen)

	imgdir, _ := x.Get(exif.GPSImgDirection)
	a, b, _ := imgdir.Rat2(0) // retrieve first (only) rat. value
	gps_dir := fmt.Sprintf("%.15f", float64(a)/float64(b))
	fmt.Println(gps_dir)
	cmdd := exec.Command("ssh", "-i", "awsEC-ubuntu.pem", "ubuntu@18.219.71.129", "source ~/.bashrc", ";", "python3", "CalCadAddr.py", "-a", str_long, "-b", str_lat, "-c", gps_dir, "-d", focallen)
	outt, err := cmdd.Output()
	if err != nil {
		fmt.Println(err)
	}
	var kk map[string]interface{}
	err = json.Unmarshal(outt, &kk)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(kk)
	// 地段地號
	dddh, ok := (kk["result"].(map[string]interface{})["cadaddr"]).(string)
	if ok {
		bc.DDDH = dddh
		fmt.Println(dddh)
	}

	bc.Lat = str_lat
	bc.Long = str_long
	bc.Hash = append(bc.Hash, hash[1])
	bc.ImgHash = append(bc.ImgHash, hash[0])
	bc.Image = img
	bc.Date = date
	bc.Chain = "Ropsten"
	fmt.Println(bc)
	res := s.db.Update(ctx, "bcposts", "_id", _id, bc)
	bcc := BCdataa{}
	err = res.Decode(&bcc)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(&bc)

	// image := &Image{ID: _id, ImgHash: hash[0], Hash: hash[1], Img: img}
	err = json.NewEncoder(w).Encode(bc)
	if err != nil {
		log.Println("err encoding return")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("Image Hash:" + hash[0])
	fmt.Println("uploaded!")
}

func (s *service) file(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	args := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(args["id"])
	if err != nil {
		log.Println("err object id from hex")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur := s.db.QueryOne(ctx, "bcposts", "_id", id)

	u := &Image{}
	err = cur.Decode(u)
	if err != nil {
		log.Println("err decoding user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		log.Println("err encoding image to json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (s *service) allFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Total-Count", "319")
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")

	var images []*Image
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur, err := s.db.QueryAll(ctx, "testing")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for cur.Next(ctx) {
		i := &Image{}

		err := cur.Decode(i)
		if err != nil {
			log.Println("err decoding cur")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		images = append(images, i)
	}

	err = json.NewEncoder(w).Encode(images)
	if err != nil {
		log.Println("err encoding images to json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *service) newUser(w http.ResponseWriter, r *http.Request) {
	log.Println("newuser called")
	w.Header().Set("Content-Type", "application/json")
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var reqData map[string]interface{}
	err = json.Unmarshal(d, &reqData)
	if err != nil {
		log.Println("err decoding input data")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user := &User{}
	user.ID = primitive.NewObjectID()
	for i, j := range reqData {
		fmt.Println(i, j)
		trimed := strings.TrimSpace(j.(string))
		switch i {
		case "username":
			user.Username = trimed
		case "name":
			user.Name = trimed
		case "department":
			user.Department = trimed
		case "email":
			user.Email = trimed
		case "phone":
			user.Phone = trimed
		case "password":
			user.Password = trimed
		case "identity":
			user.Identity = trimed
		default:
		}
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = s.db.Add(ctx, "testuser", user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if d, err := json.Marshal(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(d)
		}
	}
	log.Println("newuser call succeed!")

}

func (s *service) verifyUser(w http.ResponseWriter, r *http.Request) {
	log.Println("verifyuser called")
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	d := &User{}
	err = json.Unmarshal(body, d)
	if err != nil {
		log.Println("err unmarshaling user to struct")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur := s.db.QueryOne(ctx, "testuser", "username", d.Username)
	if cur.Err() != nil {
		log.Println("err querying user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	dec := &User{}
	err = cur.Decode(dec)
	if err != nil {
		log.Println("err decoding cur")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if dec.Password != d.Password {
		log.Println("wrong password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(dec.Identity)
	if err != nil {
		log.Println("err encoding json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *service) deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteuser called")
	args := mux.Vars(r)

	id, err := primitive.ObjectIDFromHex(args["id"])
	if err != nil {
		log.Println("err object from hex")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := s.db.DeleteOne(ctx, "testuser", "_id", id)
	if err != nil {
		log.Println("err delete one")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ret := map[string]interface{}{"id": result.DeletedCount}
	err = json.NewEncoder(w).Encode(ret)
	if err != nil {
		log.Println("err encoding to json response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("deleteuser succeed")
}

func (s *service) verifyUserAndReturnPost(w http.ResponseWriter, r *http.Request) {
	log.Println("verifyuser called")
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	d := &User{}
	err = json.Unmarshal(body, d)
	if err != nil {
		log.Println("err unmarshaling user to struct")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur := s.db.QueryOne(ctx, "testuser", "username", d.Username)
	if cur.Err() != nil {
		log.Println("err querying user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	dec := &User{}
	err = cur.Decode(dec)
	if err != nil {
		log.Println("err decoding cur")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if dec.Password != d.Password {
		log.Println("wrong password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	curr, err := s.db.Query(ctx, "posts", "user", dec.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer curr.Close(ctx)

	var k []*Post
	for curr.Next(ctx) {
		// elem := &bson.D{}
		// var elem map[string]interface{}
		elem := &Post{}
		if err := curr.Decode(&elem); err != nil {
			log.Println(err)
		}
		fmt.Println(elem)
		k = append(k, elem)
	}

	err = json.NewEncoder(w).Encode(k)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("verifyuser call succeed!")
}

func (s *service) user(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	args := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(args["id"])
	if err != nil {
		log.Println("err object id from hex")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur := s.db.QueryOne(ctx, "testuser", "_id", id)

	u := &User{}
	err = cur.Decode(u)
	if err != nil {
		log.Println("err decoding user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		log.Println("err encoding user to json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (s *service) allUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Total-Count", "319")
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")

	var users []*User
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur, err := s.db.QueryAll(ctx, "testuser")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for cur.Next(ctx) {
		u := &User{}

		err := cur.Decode(u)
		if err != nil {
			log.Println("err decoding cur")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		users = append(users, u)
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Println("err encoding users to json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *service) updateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("updateuser called")
	args := mux.Vars(r)
	val, err := primitive.ObjectIDFromHex(args["id"])
	if err != nil {
		fmt.Println("err converting id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("err read data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var reqData User
	err = json.Unmarshal(data, &reqData)
	if err != nil {
		fmt.Println("err marshal data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur := s.db.Update(ctx, "testuser", "_id", val, reqData)

	ret := &User{}
	err = cur.Decode(ret)
	if err != nil {
		fmt.Println("err cur decode")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ret)
	if err != nil {
		fmt.Println("err marshal json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("updateuser call succeed!")
}

func (s *service) newPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("new post called")
	r.ParseMultipartForm(32 << 20)
	formData := r.MultipartForm
	files := formData.File["file"]

	tag := r.FormValue("tag")
	title := r.FormValue("title")
	user := r.FormValue("user")
	name := r.FormValue("name")
	factory := r.FormValue("factory")
	market := r.FormValue("market")
	amount, err := strconv.Atoi(r.FormValue("amount"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	progress := r.FormValue("progress")

	path := "file/" + tag + "/"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}
	for k, _ := range files {
		file, err := files[k].Open()
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		p := path + files[k].Filename
		fmt.Println(p)
		f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		io.Copy(f, file)
		f.Close()
	}

	filename := "" + path

	_id := primitive.NewObjectID()
	_date := time.Now().Add(time.Hour * 8).Format(time.ANSIC)

	post := &Post{ID: _id, Tag: tag, Title: title, User: user, Name: name, Date: _date, Factory: factory, Market: market, Amount: amount, Progress: progress, Paperwork: filename}
	bc := &BCdataa{ID: _id, Tag: tag, Name: name, Factory: factory, ImgHash: []string{}, Hash: []string{}}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = s.db.Add(ctx, "posts", post)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		_, err = s.db.Add(ctx, "bcposts", bc)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if d, err := json.Marshal(post); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(d)
		}
	}
}

func (s *service) deletePost(w http.ResponseWriter, r *http.Request) {
	log.Println("deletepost called")
	args := mux.Vars(r)

	id, err := primitive.ObjectIDFromHex(args["id"])
	if err != nil {
		log.Println("err object from hex")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := s.db.DeleteOne(ctx, "posts", "_id", id)
	if err != nil {
		log.Println("err delete one")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err = s.db.DeleteOne(ctx, "bcposts", "_id", id)
	if err != nil {
		log.Println("err delete one")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ret := map[string]interface{}{"id": result.DeletedCount}
	err = json.NewEncoder(w).Encode(ret)
	if err != nil {
		log.Println("err encoding to json response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("deletepost succeed")
}

func (s *service) newBcPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// var reqData map[string]string
	var reqData BCdata
	err = json.Unmarshal(d, &reqData)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(reqData)

	id := reqData.ID
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	reqData.ID = id
	// reqData.Date = reqData.ID.Timestamp().Add(8 * time.Hour).String()
	reqData.Date = reqData.ID.Timestamp().Local().String()

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = s.db.Add(ctx, "bcposts", reqData)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if d, err := json.Marshal(&reqData); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(d)
		}
	}
}

func (s *service) allPost(w http.ResponseWriter, r *http.Request) {
	log.Println("allpost called")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Total-Count", "319")
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	// p := post{}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur, err := s.db.QueryAll(ctx, "posts")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer cur.Close(ctx)

	var k []*Post
	for cur.Next(ctx) {
		// elem := &bson.D{}
		// var elem map[string]interface{}
		elem := &Post{}
		if err := cur.Decode(&elem); err != nil {
			log.Println(err)
		}
		// fmt.Println(elem)
		k = append(k, elem)
	}

	err = json.NewEncoder(w).Encode(k)
	// mar, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("allpost call succeed!")
	// w.WriteHeader(http.StatusOK)
	// w.Write(mar)
}

func (s *service) allBcPost(w http.ResponseWriter, r *http.Request) {
	log.Println("allbcpost called")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Total-Count", "319")
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	// p := post{}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur, err := s.db.QueryAll(ctx, "bcposts")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer cur.Close(ctx)

	var k []*BCdataa
	for cur.Next(ctx) {
		// elem := &bson.D{}
		// var elem map[string]interface{}
		elem := &BCdataa{}
		if err := cur.Decode(&elem); err != nil {
			log.Println(err)
		}
		fmt.Println(elem)
		k = append(k, elem)
	}

	err = json.NewEncoder(w).Encode(k)
	// mar, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("allbcpost call succeed!")
	// w.WriteHeader(http.StatusOK)
	// w.Write(mar)
}
func (s *service) post(w http.ResponseWriter, r *http.Request) {
	log.Println("post called")
	args := mux.Vars(r)
	val, err := primitive.ObjectIDFromHex(args["id"])
	fmt.Println(val)
	if err != nil {
		fmt.Println("err converting id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur := s.db.QueryOne(ctx, "posts", "_id", val)

	k := &Post{}
	err = cur.Decode(k)
	if err != nil {
		fmt.Println("err decode cur")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// ret, err := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(k)
	if err != nil {
		fmt.Println("err marshal json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("post call succeed!")
	// w.WriteHeader(http.StatusOK)
	// w.Write(ret)
}
func (s *service) bcPost(w http.ResponseWriter, r *http.Request) {
	log.Println("bcpost called")
	args := mux.Vars(r)
	val, err := primitive.ObjectIDFromHex(args["id"])
	fmt.Println(val)
	if err != nil {
		fmt.Println("err converting id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	k, err := s.findBcPost(ctx, "bcpost", "_id", val)
	if err != nil {
		fmt.Println("err decode cur")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(k)
	if err != nil {
		fmt.Println("err marshal json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("bcpost call succeed!")
	// w.WriteHeader(http.StatusOK)
	// w.Write(ret)
}

func (s *service) findBcPost(ctx context.Context, col string, key string, val interface{}) (*BCdataa, error) {
	cur := s.db.QueryOne(ctx, "bcposts", "_id", val)
	k := &BCdataa{}
	err := cur.Decode(k)
	if err != nil {
		return nil, err
	} else {
		return k, nil
	}
}

func (s *service) updatePost(w http.ResponseWriter, r *http.Request) {
	log.Println("updatepost called")
	args := mux.Vars(r)
	val, err := primitive.ObjectIDFromHex(args["id"])
	if err != nil {
		fmt.Println("err converting id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("err read data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var reqData Post
	err = json.Unmarshal(data, &reqData)
	if err != nil {
		fmt.Println("err marshal data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur := s.db.Update(ctx, "posts", "_id", val, reqData)

	ret := &Post{}
	err = cur.Decode(ret)
	if err != nil {
		fmt.Println("err cur decode")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ret)
	if err != nil {
		fmt.Println("err marshal json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("updatepost call succeed!")
	// w.WriteHeader(http.StatusOK)
	// w.Write(ret)
}

func (s *service) verifyHash(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	imghash := args["imghash"]
	txhash := args["txhash"]
	url := "https://ropsten.etherscan.io/tx/" + txhash
	fmt.Println(imghash, txhash)
	res := Cmpurlhash(imghash, url)
	if res {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	ret := map[string]string{"message": "Failed!!!"}
	_ = json.NewEncoder(w).Encode(ret)
	return
}

func Cmpurlhash(Hash string, Url string) bool {
	fmt.Println("cmpurlhash called")
	var cmp string
	res := false
	c := colly.NewCollector()
	c.OnHTML("textarea[id]", func(e *colly.HTMLElement) {
		if e.Attr("id") == "inputdata" {
			cmp = e.Text
			temp := strings.Split(cmp, "[1]")
			temp = strings.Split(temp[1], "\n")
			temp = strings.Split(temp[0], ":  ")
			fmt.Println(temp[1], Hash)
			if temp[1] == Hash {
				res = true
			}
		}
	})
	err := c.Visit(Url)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("cmpurlhash succeed")
	return res
}

func main() {
	a := NewService("localhost", "8000")
	a.Start("localhost", "27017", "testing")

}
