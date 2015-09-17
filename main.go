package main

import (
	"abtest"
	"flag"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"net/http"
	"strconv"
	"fmt"
)

var LiveABTests = []abtest.ABTest{}

func ListABTests(w http.ResponseWriter, req *http.Request) {
	r := render.New()
	r.JSON(w, http.StatusOK, LiveABTests)
}

func Participate(w http.ResponseWriter, req *http.Request) {
	r := render.New()
	uid, err := strconv.ParseInt(req.URL.Query()["uid"][0], 10, 64)

	if err != nil {
		r.JSON(w, http.StatusInternalServerError, "Can't parse uid")
	}

	resp := abtest.Participate(LiveABTests, uid)

	r.JSON(w, http.StatusOK, resp)
}

func main() {
	listenPtr := flag.String("listen", "localhost:8000", "host and port to listen on")
	flag.Parse()

	// list of live tests
	LiveABTests = append(LiveABTests, abtest.NewABTest("test1",
		abtest.Alternative{Name: "a", Weight: 1},
		abtest.Alternative{Name: "b", Weight: 1}))

	LiveABTests = append(LiveABTests, abtest.NewABTest("test2",
		abtest.Alternative{Name: "a", Weight: 1},
		abtest.Alternative{Name: "b", Weight: 1}))

	router := mux.NewRouter()
	router.HandleFunc("/", ListABTests)
	router.HandleFunc("/participate", Participate)

	http.Handle("/", router)

	fmt.Println("Listening on " + *listenPtr)
	http.ListenAndServe(*listenPtr, nil)
}
