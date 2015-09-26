package handlers

import (
	"bouncer/experiment"
	"bouncer/config"
	"github.com/unrolled/render"
	"net/http"
)

func ListExperiments(cfg config.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		r := render.New()
		r.JSON(w, http.StatusOK, cfg.Experiments)
	}
}

func Participate(cfg config.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		r := render.New()
		var uid string
		q := req.URL.Query()

		if len(q) == 0 || len(q["uid"]) == 0 || len(q["uid"][0]) == 0 {
			r.JSON(w, http.StatusBadRequest, "Can't parse uid")
			return
		} else {
			uid = q["uid"][0]
		}

		resp := experiment.Participate(cfg.Experiments, uid)

		r.JSON(w, http.StatusOK, resp)
	}
}
