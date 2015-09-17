package main

import (
	"bouncer/abtest"
	"bouncer/handlers"
	"flag"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
)

type BouncerDB struct {
	ActiveABTests []abtest.ABTest
}

func (db BouncerDB) GetABTests() []abtest.ABTest {
	return db.ActiveABTests
}

var db BouncerDB

func init() {
	db.ActiveABTests = make([]abtest.ABTest, 0)

	// setup abtests for demo
	db.ActiveABTests = append(db.ActiveABTests, abtest.NewABTest("test1",
		abtest.Alternative{Name: "a", Weight: 1},
		abtest.Alternative{Name: "b", Weight: 1}))

	db.ActiveABTests = append(db.ActiveABTests, abtest.NewABTest("test2",
		abtest.Alternative{Name: "a", Weight: 1},
		abtest.Alternative{Name: "b", Weight: 1}))
}

func main() {
	listenPtr := flag.String("listen", "localhost:8000", "host and port to listen on")
	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/", handlers.ListABTests(db))
	router.HandleFunc("/participate/", handlers.Participate(db))

	http.Handle("/", router)

	fmt.Println("Listening on " + *listenPtr)
	http.ListenAndServe(*listenPtr, nil)
}
