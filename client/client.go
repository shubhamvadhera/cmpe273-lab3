package main

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

/*-----Debugging variables-----*/

var out io.Writer
var debugModeActivated bool

/*-----Constants-----*/

const server1 string = "http://localhost:3000"
const server2 string = "http://localhost:3001"
const server3 string = "http://localhost:3002"

/*-----Global Variables-----*/
var server1Hash uint
var server2Hash uint
var server3Hash uint

/*----------Structs----------*/

//KeyVal struct to hold keyVal pair
type KeyVal struct {
	Key   int    `json:"key"`
	Value string `json:"value"`
}

/*----------Functions----------*/

//construct function acts as program constructor
func construct() {
	debugModeActivated = false //change to true to see all developer messages
	out = ioutil.Discard
	if debugModeActivated {
		out = os.Stdout
	}
	calServerHashes()
}

//sendPut sends PUT requests
func sendPut(server string, data KeyVal) (statusCode int) {
	url := server + "/" + "keys" + "/" + strconv.Itoa(data.Key) + "/" + data.Value
	client := &http.Client{}
	req, _ := http.NewRequest("PUT", url, nil)
	res, _ := client.Do(req)
	foo, _ := json.Marshal(data)
	fmt.Println("\nPUT Request:", string(foo))
	defer res.Body.Close()
	fmt.Println("\nResponse from ", server, ":", res.StatusCode)
	return res.StatusCode
}

//sendGetAll sends GET request for all keys
func sendGetAll(server string) (keyValArr []KeyVal, statusCode int) {
	url := server + "/" + "keys"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := client.Do(req)
	defer res.Body.Close()
	fmt.Println("\nGET Request: All")

	jsonDataFromHTTP, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Panic@sendGetAll.ioutil.ReadAll")
		panic(err)
	}

	if err := json.Unmarshal([]byte(jsonDataFromHTTP), &keyValArr); err != nil {
		fmt.Println("Panic@sendGetAll.json.Unmarshal")
		panic(err)
	}

	fmt.Println("\nResponse from ", server, ":", res.StatusCode)
	fmt.Println(string(jsonDataFromHTTP))
	return keyValArr, res.StatusCode
}

//sendGetOne sends GET request for a key
func sendGetOne(server string, key int) (value string, statusCode int) {
	url := server + "/" + "keys" + "/" + strconv.Itoa(key)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := client.Do(req)
	defer res.Body.Close()
	fmt.Println("\nGET Request: ", key)

	jsonDataFromHTTP, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Panic@sendGetOne.ioutil.ReadAll")
		panic(err)
	}

	var keyVal KeyVal

	if err := json.Unmarshal([]byte(jsonDataFromHTTP), &keyVal); err != nil {
		fmt.Println("Panic@sendGetOne.json.Unmarshal")
		panic(err)
	}

	fmt.Println("\nResponse from ", server, ":", res.StatusCode)
	fmt.Println(string(jsonDataFromHTTP))
	return keyVal.Value, res.StatusCode
}

func consistentPut(key int, val string) {
	var data KeyVal
	data.Key = key
	data.Value = val
	server := mapper(key)
	sendPut(server, data)
}

func consistentGet(key int) {
	server := mapper(key)
	sendGetOne(server, key)
}

func consistentGetAll() {
	sendGetAll(server1)
	sendGetAll(server2)
	sendGetAll(server3)
}

/*----------Helper Functions----------*/

//maps a key to server via consistent hashing
func mapper(key int) (server string) {
	keyHash := crcHash(string(key))
	for keyHash <= server3Hash {
		switch keyHash {
		case server2Hash:
			return server2
		case server1Hash:
			return server1
		case server3Hash:
			return server3
		default:
			keyHash++
		}
	}
	return server2
}

//initiates server hash value variables
func calServerHashes() {
	server1Hash = crcHash(server1)
	server2Hash = crcHash(server2)
	server3Hash = crcHash(server3)
	fmt.Fprintln(out, "server1Hash:", server1Hash)
	fmt.Fprintln(out, "server2Hash:", server2Hash)
	fmt.Fprintln(out, "server3Hash:", server3Hash)
}

//returns crc32 Hash of a string
func crcHash(plain string) uint {
	hash := crc32.ChecksumIEEE
	return uint(hash([]byte(plain)))
}

func main() {
	construct()
	consistentPut(1, "a")
	consistentPut(2, "b")
	consistentPut(3, "c")
	consistentPut(4, "d")
	consistentPut(5, "e")
	consistentPut(6, "f")
	consistentPut(7, "g")
	consistentPut(8, "h")
	consistentPut(9, "i")
	consistentPut(10, "j")
	consistentGet(1)
	consistentGet(2)
	consistentGet(3)
	consistentGet(4)
	consistentGet(5)
	consistentGet(6)
	consistentGet(7)
	consistentGet(8)
	consistentGet(9)
	consistentGet(10)
	consistentGetAll()
}
