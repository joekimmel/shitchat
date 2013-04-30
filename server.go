package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
    "time"
)

//Globals!!!   Dirty, Dirty Globals
//these get initialized in the init() fn.
var POSTS *list.List //list of msgs posted to the server
var MAXPOSTS int //max number of msgs the server keeps in memory
var USERS map[int]User //map from id to user struct

//in an earlier version i was using templates to
//render the page, and that seems like
//it might be a good idea again someday? but for now
//the page template is just an empty struct...
type PageTemplate struct {
}

type Post struct {
	ID   int64
	User string
	Msg  string
}

type SendPosts struct {
	Posts []Post
}

type User struct {
	Name     string
	Id       int
	LastPost int64
}

func page_writer(w http.ResponseWriter) {
	t, err := template.ParseFiles("pagetmp2.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	template_info := PageTemplate{}
	t.Execute(w, template_info)
}

//@ "/" -- http to the root returns either the basic chat page or just "nope":
func page_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		page_writer(w)
	} else {
		fmt.Fprintf(w, "nope")
	}
}

//takes an http request, the body of which is stringified JSON of a msgs post,
//marshals it into the Post struct, inserts the struct into our global POSTS
//list...   this isn't the cleanest separation of concerns.
func request2Post(r *http.Request) {
	// the body of a request is an io.ReadCloser
	// so we should defer closing it before reading it...
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	var post Post
	json.Unmarshal(body, &post)
	post.ID = time.Now().UnixNano()

	POSTS.PushFront(post)
	if POSTS.Len() > MAXPOSTS {
		POSTS.Remove(POSTS.Back())
	}
}

//return the messages with ID > last_ind
func get_msgs_since(last_ind int64) []byte {
	retPosts := list.New()
	for e, i := POSTS.Front(), 0; e != nil; e, i = e.Next(), i+1 {
		if e.Value.(Post).ID > last_ind {
			retPosts.PushFront(e.Value)
		}
	}
	if retPosts.Len() > 0 {
		howmany := retPosts.Len()
		mposts := make([]Post, howmany)
		for e, i := retPosts.Front(), 0; e != nil; e, i = e.Next(), i+1 {
			mposts[i] = e.Value.(Post)
		}
		bb, _ := json.Marshal(mposts)
		return bb
	}
	return nil
}

func flush_msgs() {
	POSTS = list.New()
}

//@.../msgs/
func msgs_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		request2Post(r)
		fmt.Fprintf(w, "")
	} else if r.Method == "GET" {
		last_int, _ := strconv.ParseInt(r.FormValue("last_ind"), 10, 64)
		last_ind := int64(last_int)
		bbytes := get_msgs_since(last_ind)
		w.Write(bbytes)
	} else if r.Method == "DELETE" {
		flush_msgs()
	}else {
		fmt.Fprintf(w, "nope")
	}
}

//@.../user/
func user_handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		w.Write(new_user(r.FormValue("name")))
	case "GET":
		w.Write(get_users())
	case "DELETE":
		{
			id, _ := strconv.Atoi(r.FormValue("ID"))
			delete_user(id)
			fmt.Fprintf(w, "OK")
		}
	case "PUT":
		{
			id, _ := strconv.Atoi(r.FormValue("ID"))
			update_user_name(id, r.FormValue("name"))
		}
	}
}

//makes a new user with specified name,
//returns a marshalled json struct of the new user
func new_user(name string) []byte {
	id := len(USERS)
	user := User{name, id, 0}
	USERS[id] = user
	bb, _ := json.Marshal(user)
	return bb
}

//returns a marshalled json list of all the users.
func get_users() []byte {
	musers := make([]User, len(USERS))
	i := 0
	for _, user := range USERS {
		musers[i] = user
		i++
	}
	bb, _ := json.Marshal(musers)
	return bb
}

func update_user_name(id int, new_name string) {
	USERS[id] = User{new_name, id, 0}
}

func delete_user(id int) {
	delete(USERS, id)
}

func init() {
	POSTS = list.New()
	MAXPOSTS = 10
	USERS = make(map[int]User)
}

func main() {
	http.HandleFunc("/msgs/", msgs_handler)
	http.HandleFunc("/user/", user_handler)
	http.HandleFunc("/", page_handler)
	log.Fatal(http.ListenAndServe(":2999", nil))
}
