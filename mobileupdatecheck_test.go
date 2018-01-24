package mobileupdatecheck

import (
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

	}
}

/*
func TestForbiddenQueries(t *testing.T) {
	{
		req, err := http.NewRequest("GET", "/", nil)
		assert.Equal(t, nil, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(handler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
		assert.Equal(t, "", rr.Body.String())
	}
	{
		req, err := http.NewRequest("GET",
			"/foo/barapp?osVersion=1.0.0&productVersion=1.2.3", nil)
		assert.Equal(t, nil, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(handler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
		assert.Equal(t, "", rr.Body.String())
	}
}

func TestHandlerInvalidQueries(t *testing.T) {
	{
		req, err := http.NewRequest("GET", "/ios/fitapp", nil)
		assert.Equal(t, nil, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(handler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "", rr.Body.String())
	}
	{
		req, err := http.NewRequest("GET", "/ios/fitapp?osVersion=1.0.0", nil)
		assert.Equal(t, nil, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(handler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "", rr.Body.String())
	}
}

*/

/*
func TestHandlerGoodQueries(t *testing.T) {
	for _, tc := range testCases {
		{
			req, err := http.NewRequest("GET", "/"+tc.os+"/"+tc.product+
				"?osVersion="+tc.osVersion+"&productVersion="+tc.productVersion, nil)
			assert.Equal(t, nil, err)

			rr := httptest.NewRecorder()
			http.HandlerFunc(handler).ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
			assert.Equal(t, `{ "action": "`+tc.expectedAction+`" }`, rr.Body.String())
		}
	}
}
*/
