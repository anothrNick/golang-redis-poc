/*
Simple HTTP server that stores and returns request count for a window of time, as a POC for rate limiting with redis.

This example keeps track of a global rate limit per minute.
*/
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis"
)

// RateLimit is a global rate limit
const RateLimit int = 30

// Cache is a simple struct that will be rendered as the response JSON
type Cache struct {
	Requests int `json:"requests"`
}

func main() {
	// create a new redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "cache:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// ping the redis server
	pong, err := client.Ping().Result()

	// fail quickly if ping fails
	if err != nil {
		fmt.Println(pong, err)
		return
	}

	// single http handler for demo
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// get current minute for rate limit
		t := time.Now().Minute()
		currentKey := fmt.Sprintf("requestCount:%d", t)

		// get the request count key for the current window
		val, err := client.Get(currentKey).Int()

		if err == redis.Nil {
			fmt.Println(currentKey + " does not exist")
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if val > RateLimit {
			http.Error(w, "Rate limit reached.", http.StatusTooManyRequests)
			return
		}

		// set response as current value
		cache := Cache{Requests: val}

		// increment and set request count in redis
		val++
		// expire key after 1 minute
		err = client.Set(currentKey, val, time.Minute*1).Err()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// marshal json response
		js, err := json.Marshal(cache)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// return response json
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})

	http.ListenAndServe(":5001", nil)
}
