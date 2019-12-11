package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-redis/redis"
)

var HTML []byte
var mux *http.ServeMux

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func addToRedisList(key, value string) {
	client.RPush(key, value)
}

func main() {

	// Mux for static files
	mux = http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./public")))

	HTML, _ = ioutil.ReadFile("./public/index.html")

	http.HandleFunc("/", static(indexHandler))
	http.HandleFunc("/post", static(postHandler))

	http.ListenAndServe(":8888", nil)
}

// There are better ways to check for static files,
// this is just a lazy way which checks if the url
// contains a dot, if so it will assume it's static
func static(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.ContainsRune(r.URL.Path, '.') {
			mux.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(HTML)
}

type Response struct {
	RedisKey string `json:"redis_key"`
	RedisVal string `json:"redis_val"`
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	data := new(Response)

	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	addToRedisList(data.RedisKey, data.RedisVal)

	body, _ := ioutil.ReadAll(r.Body)
	bodyString := string(body)
	fmt.Println("bodystring: ", bodyString)
	fmt.Println("data: ", data)

	value := client.LIndex(data.RedisKey, -1)
	resp := &Response{RedisKey: data.RedisKey, RedisVal: value.Val()}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
}

type ResponseSql struct {
	Tag1 string `json:"tag1"`
	Tag2 string `json:"tag2"`
	Tag3 string `json:"tag3"`
}
