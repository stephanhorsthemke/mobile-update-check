package main

import (
	"context"
	"github.com/Masterminds/semver"
	uc "github.com/egymgmbh/mobile-update-check/pb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRuleCompileConstraints(t *testing.T) {
	{
		r := rule{
			os:             uc.OSType_ANDROID,
			osVersion:      ">=8.0.0 && <9.0.0",
			product:        uc.ProductType_FITAPP,
			productVersion: "<=2.3.0",
			action:         uc.ResponseAction_ADVICE,
		}

		err := r.CompileConstraints()
		assert.NotEqual(t, nil, err)
	}
	{
		r := rule{
			os:             uc.OSType_ANDROID,
			osVersion:      ">=8.0.0, <9.0.0",
			product:        uc.ProductType_FITAPP,
			productVersion: "<=2.3.0",
			action:         uc.ResponseAction_ADVICE,
		}

		err := r.CompileConstraints()
		assert.Equal(t, nil, err)
	}
	{
		r := rule{
			os:             uc.OSType_ANDROID,
			osVersion:      ">=8.0.0, <9.0.0",
			product:        uc.ProductType_FITAPP,
			productVersion: "<=2..0",
			action:         uc.ResponseAction_ADVICE,
		}

		err := r.CompileConstraints()
		assert.NotEqual(t, nil, err)
	}

}

func TestRuleApply(t *testing.T) {
	{
		r := rule{
			os:             uc.OSType_ANDROID,
			osVersion:      ">=8.0.0, <9.0.0",
			product:        uc.ProductType_FITAPP,
			productVersion: "<=2.3.0",
			action:         uc.ResponseAction_ADVICE,
		}

		err := r.CompileConstraints()
		assert.Equal(t, nil, err)

		{
			osTestVersion, _ := semver.NewVersion("8.3.0")
			productTestVersion, _ := semver.NewVersion("2.3.0")
			testUpdateCheckRequest := &uc.UpdateCheckRequest{
				uc.OSType_IOS,
				"8.3.0",
				uc.ProductType_FITAPP,
				"2.3.0",
			}

			action, valid := r.Apply(
				testUpdateCheckRequest,
				osTestVersion,
				productTestVersion)

			assert.Equal(t, uc.ResponseAction_NONE, action)
			assert.False(t, valid)
		}

		{
			osTestVersion, _ := semver.NewVersion("8.3.0")
			productTestVersion, _ := semver.NewVersion("2.3.0")
			testUpdateCheckRequest := &uc.UpdateCheckRequest{
				uc.OSType_ANDROID,
				"8.3.0",
				uc.ProductType_FITAPP,
				"2.3.0",
			}

			action, valid := r.Apply(
				testUpdateCheckRequest,
				osTestVersion,
				productTestVersion)

			assert.Equal(t, uc.ResponseAction_ADVICE, action)
			assert.True(t, valid)
		}

		{
			osTestVersion, _ := semver.NewVersion("8.3.0")
			productTestVersion, _ := semver.NewVersion("2.4.0")
			testUpdateCheckRequest := &uc.UpdateCheckRequest{
				uc.OSType_ANDROID,
				"8.3.0",
				uc.ProductType_FITAPP,
				"2.4.0",
			}

			action, valid := r.Apply(
				testUpdateCheckRequest,
				osTestVersion,
				productTestVersion)

			assert.Equal(t, uc.ResponseAction_NONE, action)
			assert.False(t, valid)
		}
	}
}

func TestUcsQuery(t *testing.T) {

	testUpdateCheckRequest := &uc.UpdateCheckRequest{
		uc.OSType_ANDROID,
		"8.3.0",
		uc.ProductType_FITAPP,
		"2.4.0",
	}

	ucs := ucs{}
	testUpdateCheckResponse, err := ucs.Query(context.Background(), testUpdateCheckRequest)
	assert.Equal(t, nil, err)
	assert.Equal(t, uc.ResponseAction_NONE, testUpdateCheckResponse.GetAction())
}
