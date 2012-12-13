// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

//
// The Neo4j Manual section numbers quoted herein refer to the manual for 
// milestone release 1.8.M06.  http://docs.neo4j.org/chunked/milestone/

package restclient

import (
	"github.com/bmizerany/assert"
	"github.com/bmizerany/pat"
	"log"
	// "sort"
	"encoding/json"
	"net/http"
	// "io/ioutil"
	//	"net/url"
	"testing"
)

const (
	port = "9000"
)

type structType struct {
	Foo map[string]string
	Bar int
	Baz string
}

var (
	fooMap    = map[string]string{"foo": "bar"}
	fooStruct = structType{
		Foo: fooMap,
		Bar: 222,
		Baz: "baz",
	}
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	// 
	// Routing
	//
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(HandleGET))
	//
	// Start webserver
	//
	http.Handle("/", mux)
	go func() {
		log.Println("Starting webserver on port " + port + "...")
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Panicln(err)
		}
	}()
}

func HandleGET(w http.ResponseWriter, req *http.Request) {
	u := req.URL
	q := u.Query()
	log.Println(q)
	for k, _ := range fooMap {
		if fooMap[k] != q.Get(k) {
			msg := "Bad query params: " + u.Query().Encode()
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	}
	//
	// Generate response
	//
	blob, err := json.Marshal(fooStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Add("content-type", "application/json")
	w.Write(blob)
}

/*
func HandlePOST(w http.ResponseWriter, req *http.Request) {
	//
	// Parse Payload
	//
	if req.ContentLength <= 0 {
		msg := "Content-Length must be greater than 0."
		http.Error(w, msg, http.StatusLengthRequired)
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var greq GetRequest
	err = json.Unmarshal(body, &greq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	valid := true
	switch true {
	case greq.Foo != "foo":
		valid = false
	case greq.Bar != fooBarMap:
		valid = false
	case greq.Baz != 111:
		valid = false
	}
	if valid == false {
		msg := "Bad request body"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	//
	// Compose Response
	//
	//
	// Generate response
	//
	blob, err := json.Marshal(gresp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Add("content-type", "application/json")
	w.Write(blob)
}
*/

func TestGET(t *testing.T) {
	c := New()
	r := RestRequest{
		Url:    "http://localhost:" + port,
		Method: GET,
		Params: fooMap,
		// Params: map[string]string{"bad": "value"},
		Result: new(structType),
	}
	status, err := c.Do(&r)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, status, 200)
	assert.Equal(t, r.Result, &fooStruct)
}