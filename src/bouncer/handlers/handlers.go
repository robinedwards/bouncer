package handlers

import (
	"bouncer/config"
	"bouncer/experiment"
	"bouncer/feature"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/unrolled/render"
	"io/ioutil"
	"net/http"
	"runtime/debug"
)

type ConfigFactory func () *config.Config

type Context struct {
	Uid string `json:"uid"`
}

type ParticipateResponse struct {
	Experiments map[string]string `json:"experiments,omitempty"`
	Features    map[string]bool   `json:"features,omitempty"`
}

type ParticipateRequest struct {
	Context     Context             `json:"context"`
	Experiments map[string][]string `json:"experiments,omitempty"`
	Features    map[string]float32  `json:"features,omitempty"`
}

func ErrorHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rval := recover(); rval != nil {
				rvalStr := fmt.Sprint(rval)
				packet := raven.NewPacket(rvalStr,
					raven.NewException(errors.New(rvalStr), raven.NewStacktrace(2, 3, nil)),
					raven.NewHttp(r))
				raven.Capture(packet, nil)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("500 Server Error: %s\n\n", rval)))
					w.Write(debug.Stack())
			}
		}()

		handler.ServeHTTP(w, r)
	})
}

func Root(w http.ResponseWriter, req *http.Request) {
	r := render.New()
	r.JSON(w, http.StatusOK,
		map[string][]string{
			"app":   {"bouncer"},
			"paths": {"/experiments/", "/features/", "/groups/", "/participate/", "/stats/"},
		})
}

func ListExperiments(cfg ConfigFactory) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		r := render.New()
		r.JSON(w, http.StatusOK, cfg().Experiments)
	}
}

func ListFeatures(cfg ConfigFactory) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		r := render.New()
		r.JSON(w, http.StatusOK, cfg().Features)
	}
}

func ListGroups(cfg ConfigFactory) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		r := render.New()
		r.JSON(w, http.StatusOK, cfg().Groups)
	}
}

func CheckFeatures(features map[string]float32, uid string, cfg config.Config) map[string]bool {
	r := make(map[string]bool)

	if len(features) == 0 {
		for featureName, f := range cfg.FeatureMap {
			r[featureName] = f.IsEnabled(uid)
		}

		return r
	}

	for featureName, _ := range features {
		if f, ok := cfg.FeatureMap[featureName]; ok {
			r[f.Name] = f.IsEnabled(uid)
		} else {
			// Un-configured feature specified by the client.
			f := feature.NewFeature(featureName, features[featureName], make(map[string]int))
			r[f.Name] = f.IsEnabled(uid)
		}
	}

	return r
}

func CheckExperiments(experiments map[string][]string, uid string, cfg config.Config) map[string]string {
	r := make(map[string]string)
	// if we don't specify experiments to participate in return all
	if len(experiments) == 0 {
		for experimentName, exp := range cfg.ExperimentMap {
			r[experimentName] = exp.GetAlternative(uid)
		}
		return r
	}

	// else only return requested experiments
	for experimentName, _ := range experiments {
		if e, ok := cfg.ExperimentMap[experimentName]; ok {
			r[experimentName] = e.GetAlternative(uid)
		} else {
			// Un-configured feature specified by the client.
			alts := make([]experiment.Alternative, len(experiments[experimentName]))
			for _, alternativeName := range experiments[experimentName] {
				alts = append(alts, experiment.Alternative{alternativeName, 1})
			}

			e := experiment.NewExperiment(experimentName, make(map[string]string), alts...)
			r[experimentName] = e.GetAlternative(uid)
		}
	}
	return r
}

func Participate(config ConfigFactory, logger func (interface{})) func(http.ResponseWriter, *http.Request) {
	cfg := *config()

	return func(w http.ResponseWriter, req *http.Request) {
		r := render.New()
		if req.Method != "POST" {
			r.JSON(w, http.StatusMethodNotAllowed, "Method not supported")
			fmt.Errorf("ERROR: method not supported")
			return
		}

		preq := ParticipateRequest{}

		body, rerr := ioutil.ReadAll(req.Body)
		if rerr != nil {
			r.JSON(w, http.StatusBadRequest, fmt.Sprintf("Error reading body: %s", rerr))
			fmt.Errorf("ERROR: reading body: %v", rerr)
			return
		}

		err := json.Unmarshal(body, &preq)
		if err != nil {
			r.JSON(w, http.StatusBadRequest, fmt.Sprintf("Error decoding json: %s", err))
			fmt.Errorf("ERROR: decoding json: %s received: %s", err, string(body))
			return
		}

		presp := new(ParticipateResponse)
		presp.Experiments = CheckExperiments(preq.Experiments, preq.Context.Uid, cfg)
		presp.Features = CheckFeatures(preq.Features, preq.Context.Uid, cfg)

		r.JSON(w, http.StatusOK, presp)
		go logger(*presp)
	}
}

