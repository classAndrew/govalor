package util

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// ReqPerMinute ratelimits
var ReqPerMinute = 60

// Get sends get request for JSON response
func Get(ch chan<- map[string]interface{}, url string, params map[string]string) error {
	queryParamStr := "?"
	for k, v := range params {
		v := strings.ReplaceAll(v, " ", "%20")
		queryParamStr += k + "=" + v + "&"
	}
	// toss out the last character whether it's an '&' or '?' (no params)
	queryParamStr = queryParamStr[:len(queryParamStr)-1]
	res, err := http.Get(url + queryParamStr)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer res.Body.Close()
	var respJSON map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&respJSON); err != nil {
		return err
	}
	ch <- respJSON
	return nil
}
