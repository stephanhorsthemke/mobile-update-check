package mobileupdatecheck

import (
	"log"
	"net/http"

	"github.com/Masterminds/semver"
)

const (
	actionNone        = "NONE"
	actionUpdate      = "UPDATE"
	actionForceUpdate = "FORCE_UPDATE"
)

var (
	rules = make(map[string][]*rule)
)

type rule struct {
	osVersion          string
	osConstraints      *semver.Constraints
	productVersion     string
	productConstraints *semver.Constraints
	action             string
}

// compileRules pre-computes version contraints in the update rules to reduce
// runtime latency and avoid duplicate work
func compileRules() error {
	for key := range rules {
		for _, rule := range rules[key] {
			var err error
			rule.osConstraints, err = semver.NewConstraint(rule.osVersion)
			if err != nil {
				return err
			}
			rule.productConstraints, err = semver.NewConstraint(rule.productVersion)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// handler checks if the request path is a valid rule set identifier and applies
// all rules in a set if there is a match
func handler(w http.ResponseWriter, r *http.Request) {
	ruleSet := r.URL.Path[1:]
	if _, ok := rules[ruleSet]; !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	osVersion, err := semver.NewVersion(r.URL.Query().Get("osVersion"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	productVersion, err := semver.NewVersion(r.URL.Query().Get("productVersion"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, rule := range rules[ruleSet] {
		if rule.osConstraints.Check(osVersion) &&
			rule.productConstraints.Check(productVersion) {
			w.Write([]byte(`{ "action": "` + rule.action + `" }`))
			return
		}
	}
	w.Write([]byte(`{ "action": "` + actionNone + `" }`))
}

func init() {
	rules["ios/fitapp"] = []*rule{} // empty rule set
	rules["ios/trainerapp"] = []*rule{
		{
			osVersion:      "9.0.0",
			productVersion: "2.3.0",
			action:         actionForceUpdate,
		},
	}
	rules["android/fitapp"] = []*rule{{
		osVersion:      ">=8.0.0, <9.0.0",
		productVersion: "<=2.3.0",
		action:         actionUpdate,
	}}

	if err := compileRules(); err != nil {
		log.Fatalf("compile rules: %v", err)
	}

	http.HandleFunc("/", handler)
}
