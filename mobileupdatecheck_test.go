package main

import (
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuleCompile(t *testing.T) {
	// non-existent rules file
	{
		_, err := loadRuleSets(path.Join("testdata", "non-existent.json"))
		assert.NotEqual(t, nil, err)
	}
	// empty rule is invalid
	{
		_, err := loadRuleSets(path.Join("testdata", "empty.json"))
		assert.NotEqual(t, nil, err)
	}
	// invalid os version
	{
		_, err := loadRuleSets(path.Join("testdata", "bad-os.json"))
		assert.NotEqual(t, nil, err)
	}
	// invalid product version
	{
		_, err := loadRuleSets(path.Join("testdata", "bad-product.json"))
		assert.NotEqual(t, nil, err)
	}
	// all good
	{
		r, err := loadRuleSets(path.Join("testdata", "all-good.json"))
		assert.Equal(t, nil, err)
		assert.Equal(t, actionUpdate, r["android/fitapp"][0].action)
	}
}

func TestBadQueries(t *testing.T) {
	compiledRuleSets, _ = loadRuleSets(path.Join("testdata", "rules.json"))
	// missing os/product combination
	{
		req, err := http.NewRequest("GET", "/", nil)
		assert.Equal(t, nil, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(handler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "", rr.Body.String())
	}
	// missing queryparams
	{
		req, err := http.NewRequest("GET", "/android/fitapp", nil)
		assert.Equal(t, nil, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(handler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "", rr.Body.String())
	}
	// invalid os/product combination
	{
		req, err := http.NewRequest("GET",
			"/foo/barapp?osVersion=1.0.0&productVersion=1.2.3", nil)
		assert.Equal(t, nil, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(handler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		assert.Equal(t, "", rr.Body.String())
	}
	// missing os version
	{
		req, err := http.NewRequest("GET",
			"/android/fitapp?productVersion=1.2.3", nil)
		assert.Equal(t, nil, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(handler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "", rr.Body.String())
	}
	// missing product version
	{
		req, err := http.NewRequest("GET",
			"/android/fitapp?osVersion=1.2.3", nil)
		assert.Equal(t, nil, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(handler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "", rr.Body.String())
	}
}

func TestGoodQueries(t *testing.T) {
	compiledRuleSets, _ = loadRuleSets(path.Join("testdata", "rules.json"))
	// rule match
	{
		req, err := http.NewRequest("GET",
			"/android/fitapp?osVersion=1.0.0&productVersion=1.0.0", nil)
		assert.Equal(t, nil, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(handler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "{ \"action\": \"ADVICE\" }", rr.Body.String())
	}
	// no rule match
	{
		req, err := http.NewRequest("GET",
			"/android/fitapp?osVersion=999.0.0&productVersion=999.0.0", nil)
		assert.Equal(t, nil, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(handler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "{ \"action\": \"NONE\" }", rr.Body.String())
	}
}
