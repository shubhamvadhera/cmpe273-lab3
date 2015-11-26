package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

/*-----Debugging variables-----*/

var out io.Writer
var debugModeActivated bool

/*----------Global Variables----------*/

var store map[int]string

/*----------Structs----------*/

//KeyVal struct to hold keyVal pair
type KeyVal struct {
	Key   int    `json:"key"`
	Value string `json:"value"`
}

/*----------Functions----------*/

//construct function acts as program constructor
func construct() {
	debugModeActivated = true //change to true to see all developer messages
	out = ioutil.Discard
	if debugModeActivated {
		out = os.Stdout
	}

	store = make(map[int]string)
}

//PutKey serves PUT request
func PutKey(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	key, _ := strconv.ParseInt(p.ByName("key_id"), 10, 64)
	val := p.ByName("value")
	store[int(key)] = val
	fmt.Println("\nPUT Request Received:", "{", "\"key\":", int(key), "\"value\":", val, "}")
	fmt.Println("\nServer ready... Waiting for requests...")
}

//GetOneKey serves GET request for a given key
func GetOneKey(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	key, _ := strconv.ParseInt(p.ByName("key_id"), 10, 64)
	fmt.Println("\nGET Request Received: key_id:", key)
	resp := KeyVal{Key: int(key), Value: store[int(key)]}
	jsonOut, _ := json.Marshal(resp)
	httpResponse(w, jsonOut, 200)
	fmt.Println("\nResponse:", string(jsonOut), "\n200 OK")
	fmt.Println("\nServer ready... Waiting for requests...")
}

//GetAllKeys serves GET request to return all key, value pairs
func GetAllKeys(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("\nGET Request Received: All")
	size := len(store)
	resp := make([]KeyVal, size)
	i := 0
	for key, val := range store {
		resp[i].Key = key
		resp[i].Value = val
		i++
	}
	jsonOut, _ := json.Marshal(resp)
	httpResponse(w, jsonOut, 200)
	fmt.Println("\nResponse:", string(jsonOut), "\n200 OK")
	fmt.Println("\nServer ready... Waiting for requests...")
}

/*----------Helper Functions----------*/

//write http response
func httpResponse(w http.ResponseWriter, jsonOut []byte, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, "%s", jsonOut)
}

func main() {
	construct()

	mux := httprouter.New()
	fmt.Println("\nServer ready... Waiting for requests...")
	mux.PUT("/keys/:key_id/:value", PutKey)
	mux.GET("/keys/:key_id", GetOneKey)
	mux.GET("/keys", GetAllKeys)

	http.ListenAndServe("localhost:3001", mux)
}
