package handlers

import (
	"net/http"
	"bouncer/experiment"
	"github.com/unrolled/render"
	"strconv"
)

type BouncerDB interface {
	GetExperiments() []experiment.Experiment
}

func ListExperiments(db BouncerDB) func(http.ResponseWriter, *http.Request) {
	return func (w http.ResponseWriter, req *http.Request) {
		r := render.New()
		r.JSON(w, http.StatusOK, db.GetExperiments())
	}
}

func Participate(db BouncerDB) func(http.ResponseWriter, *http.Request) {
	return func (w http.ResponseWriter, req *http.Request) {
		r := render.New()
		uid, err := strconv.ParseInt(req.URL.Query()["uid"][0], 10, 64)

		if err != nil {
			r.JSON(w, http.StatusBadRequest, "Can't parse uid")
		}

		resp := experiment.Participate(db.GetExperiments(), uid)

		r.JSON(w, http.StatusOK, resp)
	}
}
