package main

import (
	"testing"
	uc "github.com/egymgmbh/mobile-update-check/pb"
	"github.com/stretchr/testify/assert"
)

func TestRuleCompileContraints(t *testing.T) {
{
r := rule{
  				os: uc.OSType_ANDROID,
  				osVersion: ">=8.0.0 && <9.0.0",
  				product: uc.ProductType_FITAPP,
  				productVersion: "<=2.3.0",
  				action: uc.ResponseAction_ADVICE,
}

  err := r.CompileConstraints()
  assert.NotEqual(t, nil, err)
}
{
r := rule{
  				os: uc.OSType_ANDROID,
  				osVersion: ">=8.0.0, <9.0.0",
  				product: uc.ProductType_FITAPP,
  				productVersion: "<=2.3.0",
  				action: uc.ResponseAction_ADVICE,
}

  err := r.CompileConstraints()
  assert.Equal(t, nil, err)
}

}
