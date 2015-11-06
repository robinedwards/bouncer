package handlers

import (
	"bouncer/experiment"
	"bouncer/config"
	"bouncer/feature"
	"github.com/unrolled/render"
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func Root(w http.ResponseWriter, req *http.Request) {
	r := render.New()
	r.JSON(w, http.StatusOK, map[string]string{"App": "bouncer"})
}

func ListExperiments(cfg config.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		r := render.New()
		r.JSON(w, http.StatusOK, cfg.Experiments)
	}
}

func ListFeatures(cfg config.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		r := render.New()
		r.JSON(w, http.StatusOK, cfg.Features)
	}
}

func ListGroups(cfg config.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		r := render.New()
		r.JSON(w, http.StatusOK, cfg.Groups)
	}
}

type ParticipateResponse struct {
	Experiments map[string]string
	Features 	map[string]bool
}

type ParticipateRequest struct {
	Uid 		string
	Experiments map[string][]string
	Features	map[string]float32
}


func CheckFeatures(features map[string]float32, uid string, config config.Config) map[string]bool {
	r := make(map[string]bool)

	for featureName, _ := range features {
		if f, ok := config.FeatureMap[featureName]; ok {
			r[f.Name] = f.IsEnabled(uid)
		} else {
			// Un-configured feature specified by the client.
			f := feature.NewFeature(featureName, features[featureName], make(map[string]int))
			r[f.Name] = f.IsEnabled(uid)
		}
	}

	return r
}


func CheckExperiments(experiments map[string][]string, uid string, config config.Config) map[string]string {
	r := make(map[string]string)

	for experimentName, _ := range experiments {
		if e, ok := config.ExperimentMap[experimentName]; ok {
			r[experimentName] = e.GetAlternative(uid)
		} else {
			// Un-configured feature specified by the client.
			alts := make([]experiment.Alternative, len(experiments[experimentName]))
			for _, alternativeName := range(experiments[experimentName]) {
				alts = append(alts, experiment.Alternative{alternativeName, 1})
			}

			e := experiment.NewExperiment(experimentName, make(map[string]string), alts...)
			r[experimentName] = e.GetAlternative(uid)
		}
	}
	return r
}


func Participate(cfg config.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		r := render.New()

		if req.Method != "POST" {
			r.JSON(w, http.StatusMethodNotAllowed, "Method not supported")
			return
		}

		preq := ParticipateRequest{}

		body, rerr := ioutil.ReadAll(req.Body)
		if rerr != nil {
			r.JSON(w, http.StatusBadRequest, fmt.Sprintf("Error reading body: %v", rerr))
			return
		}

		err := json.Unmarshal(body, &preq)
		if err != nil {
			r.JSON(w, http.StatusBadRequest, fmt.Sprintf("Error decoding json: %c", err))
			return
		}

		presp := new(ParticipateResponse)

		presp.Experiments = CheckExperiments(preq.Experiments, preq.Uid, cfg)
		presp.Features = CheckFeatures(preq.Features, preq.Uid, cfg)

		// TODO log this to file as json.

		r.JSON(w, http.StatusOK, presp)
	}
}
