package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Masterminds/semver"
)

const (
	actionNone        = "NONE"
	actionUpdate      = "UPDATE"
	actionForceUpdate = "FORCE_UPDATE"
)

var (
	compiledRuleSets = make(map[string][]compiledRule)
)

// RuleSet represents a set of rules for a give os/product combination
type RuleSet struct {
	Key   string
	Rules []Rule
}

// Rule represents a single rule that is part of a RuleSet
type Rule struct {
	OSVersion      string
	ProductVersion string
	Action         string
}

type compiledRule struct {
	osConstraints      *semver.Constraints
	productConstraints *semver.Constraints
	action             string
}

// handler checks if the request path is a valid rule set identifier and applies
// all rules in a set if there is a match
func handler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	if _, ok := compiledRuleSets[key]; !ok {
		w.WriteHeader(http.StatusNoContent)
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
	for _, cr := range compiledRuleSets[key] {
		if cr.osConstraints.Check(osVersion) &&
			cr.productConstraints.Check(productVersion) {
			w.Write([]byte(`{ "action": "` + cr.action + `" }`))
			return
		}
	}
	w.Write([]byte(`{ "action": "` + actionNone + `" }`))
}

func loadRuleSets(fname string) (map[string][]compiledRule, error) {
	var ruleSets []RuleSet

	// read rules file
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("read rules: %v", err)
	}
	// parse rules file
	err = json.Unmarshal(buf, &ruleSets)
	if err != nil {
		return nil, fmt.Errorf("parse rules: %v", err)
	}

	// compile rules contraints
	crs := make(map[string][]compiledRule)
	for _, ruleSet := range ruleSets {
		if _, ok := crs[ruleSet.Key]; !ok {
			crs[ruleSet.Key] = []compiledRule{}
		}
		for _, rule := range ruleSet.Rules {
			var cr compiledRule
			var err error
			cr.osConstraints, err = semver.NewConstraint(rule.OSVersion)
			if err != nil {
				return nil, fmt.Errorf("parse constraint: %v", err)
			}
			cr.productConstraints, err = semver.NewConstraint(rule.ProductVersion)
			if err != nil {
				return nil, fmt.Errorf("parse constraint: %v", err)
			}
			cr.action = rule.Action
			crs[ruleSet.Key] = append(crs[ruleSet.Key], cr)
		}
	}
	return crs, nil
}

func main() {
	var err error
	compiledRuleSets, err = loadRuleSets("rules.json")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
