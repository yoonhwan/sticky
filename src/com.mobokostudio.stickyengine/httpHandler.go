package main

import (
	"appengine"
	"appengine/urlfetch"
	"appengine/user"

	"encoding/json"
	"fmt"
	"strings"

	"bytes"
	"net/http"
	"regexp"
)

var DEBUG = false
var postgresHost = "http://54.65.41.121/query"
var errorStatus = false
var validPath = regexp.MustCompile("^/([a-zA-Z0-9]+)$")

func LOGD(w http.ResponseWriter, format string) {
	if DEBUG {
		fmt.Fprintf(w, format)
		fmt.Fprintf(w, "\n")
	}
}

func LOGDF(w http.ResponseWriter, format string, v ...interface{}) {
	if DEBUG {
		fmt.Fprintf(w, format, v)
	}
}

var function_map = map[string]func(http.ResponseWriter, *http.Request){
	"userLogin":     userLogin,
	"saveLocation":  saveLocation,
	"loadLocation":  loadLocation,
	"clearLocation": clearLocation,
}

func startHTTP() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if errorStatus != true {
				LOGD(w, "SUCCESS")

			} else {
				LOGD(w, "ERROR")
			}
		}()

		errorStatus = false
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil || function_map[m[1]] == nil {
			http.NotFound(w, r)
			return
		}
		LOGDF(w, "Hello, world! %s", m[1])

		fn := function_map[m[1]]
		fn(w, r)
	})
}

func userLogin(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	var data string
	if r.Method == "POST" {
		data = r.FormValue("data")
		LOGDF(w, "POST = ", data)

	} else if DEBUG && r.Method == "GET" {
		data = r.FormValue("data")
		LOGDF(w, "GET = ", data)
	}

	if u == nil && len(data) > 0 {
		byt := []byte(data)
		var dat map[string]interface{}
		if err := json.Unmarshal(byt, &dat); err != nil {
			panic(err)
		}

		//data={\"user\":\"%@\"}
		LOGDF(w, "user:%s\n", dat["user"])

		var query = fmt.Sprintf("SELECT * FROM users WHERE name='%s'", dat["user"])
		LOGDF(w, "query = %s\n", query)
		var requestBody = fmt.Sprintf("query=%s", query)

		responseBody := requestHttp(requestBody, w, c)

		type PublicKey struct {
			Id   int
			Name string
		}
		keysBody := []byte(responseBody)
		var keys []PublicKey
		if err := json.Unmarshal(keysBody, &keys); err != nil {
			panic(err)
		}

		if len(keys) == 0 {
			//empty register
			var query = fmt.Sprintf("INSERT INTO users(name) VALUES('%s')", dat["user"])
			LOGDF(w, "query = %s\n", query)
			var requestBody = fmt.Sprintf("query=%s", query)

			responseBody := requestHttp(requestBody, w, c)

			byt := []byte(responseBody)
			var dat map[string]interface{}
			if err := json.Unmarshal(byt, &dat); err != nil {
				panic(err)
			}

			var errcode int
			errcode = int(dat["errcode"].(float64))

			if errcode != 0 {
				errorStatus = true
			} else {
				LOGDF(w, "register success output = %#v \n", dat)
			}
		} else {
			LOGDF(w, "output = %#v \n", keys)
		}
	}
}

func saveLocation(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	var data string
	if r.Method == "POST" {
		data = r.FormValue("data")
		LOGDF(w, "POST = ", data)

	} else if DEBUG && r.Method == "GET" {
		data = r.FormValue("data")
		LOGD(w, "GET = ")
	}

	if u == nil && len(data) > 0 {
		byt := []byte(data)
		var dat map[string]interface{}
		if err := json.Unmarshal(byt, &dat); err != nil {
			panic(err)
		}

		//data={\"lat\":%f,\"lon\":%f,\"timestamp\":%f,\"desc\":\"%@\",\"user\":\"%@\"}
		LOGDF(w, "user:%s location:(%f,%f) timestamp:%f desc:%s\n", dat["user"], dat["lat"], dat["lon"], dat["timestamp"], dat["desc"])
		var query = fmt.Sprintf("INSERT INTO location_trace(user_id, the_geom) SELECT id,'SRID=900913;POINT(%f %f)' FROM users WHERE name='%s'", dat["lat"], dat["lon"], dat["user"])
		LOGDF(w, "query = %s\n", query)
		var requestBody = fmt.Sprintf("query=%s", query)

		responseBody := requestHttp(requestBody, w, c)

		byt = []byte(responseBody)
		if err := json.Unmarshal(byt, &dat); err != nil {
			panic(err)
		}

		var errcode int
		errcode = int(dat["errcode"].(float64))

		if errcode != 0 {
			errorStatus = true
		} else {
			LOGD(w, "register success")
		}
	}

}

func loadLocation(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	var data string
	if r.Method == "POST" {
		data = r.FormValue("data")
		LOGDF(w, "POST = ", data)

	} else if DEBUG && r.Method == "GET" {
		data = r.FormValue("data")
		LOGD(w, "GET = ")
	}

	if u == nil && len(data) > 0 {
		byt := []byte(data)
		var dat map[string]interface{}
		if err := json.Unmarshal(byt, &dat); err != nil {
			panic(err)
		}

		//data={\"lat\":%f,\"lon\":%f,\"timestamp\":%f,\"desc\":\"%@\",\"user\":\"%@\"}
		LOGDF(w, "user:%s location:(%f,%f) timestamp:%f desc:%s\n", dat["user"], dat["lat"], dat["lon"], dat["timestamp"], dat["desc"], dat["limit"])
		
		var query = fmt.Sprintf("SELECT location_trace.id, ST_X(the_geom) as lat, ST_Y(the_geom) as lon, st_astext(the_geom) as Point FROM users, location_trace WHERE users.name='%s' AND users.id=location_trace.user_id ORDER BY id DESC LIMIT %d", dat["user"], int(dat["limit"].(float64)))
		LOGDF(w, "query = %s\n", query)
		var requestBody = fmt.Sprintf("query=%s", query)

		responseBody := requestHttp(requestBody, w, c)

		type PublicKey struct {

			//			Name string
			//			User_id   int
			//			Point string
			Lat float64
			Lon float64
		}

		keysBody := []byte(responseBody)
		var keys []PublicKey
		if err := json.Unmarshal(keysBody, &keys); err != nil {
			panic(err)
		}
		//LOGDF(w, "output = %#v \n", keys)

		var output string = "["
		for index, element := range keys {
			LOGDF(w, "%v\n", element)
			output = fmt.Sprintf("%s{\"lat\":%f,\"lon\":%f}", output, element.Lat, element.Lon)

			if index == len(keys)-1 {
				break
			}
			output = fmt.Sprintf("%s,", output)
		}
		output = fmt.Sprintf("%s]", output)
		fmt.Fprintln(w, string(output))

	}

}

func clearLocation(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	var data string
	if r.Method == "POST" {
		data = r.FormValue("data")
		LOGDF(w, "POST = ", data)

	} else if DEBUG && r.Method == "GET" {
		data = r.FormValue("data")
		LOGD(w, "GET = ")
	}

	if u == nil && len(data) > 0 {
		byt := []byte(data)
		var dat map[string]interface{}
		if err := json.Unmarshal(byt, &dat); err != nil {
			panic(err)
		}

		//data={\"lat\":%f,\"lon\":%f,\"timestamp\":%f,\"desc\":\"%@\",\"user\":\"%@\"}

		var query = fmt.Sprintf("SELECT * FROM users WHERE name='%s'", dat["user"])
		LOGDF(w, "query = %s\n", query)
		var requestBody = fmt.Sprintf("query=%s", query)

		responseBody := requestHttp(requestBody, w, c)

		type PublicKey struct {
			Id   int
			Name string
		}
		keysBody := []byte(responseBody)
		var keys []PublicKey
		if err := json.Unmarshal(keysBody, &keys); err != nil {
			panic(err)
		}

		if len(keys) > 0 {
			var query = fmt.Sprintf("DELETE FROM location_trace WHERE user_id=%d", keys[0].Id)
			LOGDF(w, "query = %s\n", query)
			var requestBody = fmt.Sprintf("query=%s", query)

			responseBody := requestHttp(requestBody, w, c)

			byt = []byte(responseBody)
			if err := json.Unmarshal(byt, &dat); err != nil {
				panic(err)
			}

			var errcode int
			errcode = int(dat["errcode"].(float64))

			if errcode != 0 {
				errorStatus = true
			} else {
				LOGD(w, "register success")
			}
		}
	}

}

func requestHttp(reqBody string, w http.ResponseWriter, c appengine.Context) (respBody string) {
	client := urlfetch.Client(c)

	req, _ := http.NewRequest("POST", postgresHost, strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	LOGDF(w, "HTTP GET returned status %v %s\n", resp.Status, buf)

	return buf.String()
}
