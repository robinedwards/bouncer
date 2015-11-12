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
	"log"
)

func Root(w http.ResponseWriter, req *http.Request) {
	r := render.New()
	r.JSON(w, http.StatusOK,
		map[string][]string{
			"App": {"bouncer"},
			"Paths": {"/experiments/", "/features/", "/groups/", "/participate/", "/stats/"},
		})
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

	if len(features) == 0 {
		for featureName, f := range config.FeatureMap {
			r[featureName] = f.IsEnabled(uid)
		}

		return r
	}

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

	// if we don't specify experiments to participate in return all
	if len(experiments) == 0 {
		for experimentName, exp := range config.ExperimentMap {
			r[experimentName] = exp.GetAlternative(uid)
		}
		return r
	}

	// else only return requested experiments
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


func Participate(cfg config.Config, ) func(http.ResponseWriter, *http.Request) {
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
			fmt.Println("ERROR: reading body: %v", rerr)
			return
		}

		err := json.Unmarshal(body, &preq)
		if err != nil {
			r.JSON(w, http.StatusBadRequest, fmt.Sprintf("Error decoding json: %v", err))
			fmt.Println("ERROR: decoding json: %v", err)
			return
		}

		presp := new(ParticipateResponse)

		presp.Experiments = CheckExperiments(preq.Experiments, preq.Uid, cfg)
		presp.Features = CheckFeatures(preq.Features, preq.Uid, cfg)

		r.JSON(w, http.StatusOK, presp)
		go logParticipation(*presp)
	}
}

func logParticipation(presp ParticipateResponse) {
	body, err := json.Marshal(presp)
	if err != nil {
		fmt.Println("ERROR: couldn't encode json response: %v", err)
	}
	log.Println(body)
}