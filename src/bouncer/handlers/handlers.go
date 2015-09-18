package handlers

import (
	"net/http"
	"bouncer/experiment"
	"github.com/unrolled/render"
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
		var uid string
		q := req.URL.Query()

		if len(q) == 0 || len(q["uid"]) == 0 || len(q["uid"][0]) == 0 {
			r.JSON(w, http.StatusBadRequest, "Can't parse uid")
			return
		} else {
			uid = q["uid"][0]
		}

		resp := experiment.Participate(db.GetExperiments(), uid)

		r.JSON(w, http.StatusOK, resp)
	}
}
