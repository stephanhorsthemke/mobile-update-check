# mobile-update-check

[![Build Status](https://travis-ci.org/egymgmbh/mobile-update-check.svg?branch=master)](https://travis-ci.org/egymgmbh/mobile-update-check)
[![codecov](https://codecov.io/gh/egymgmbh/mobile-update-check/branch/master/graph/badge.svg)](https://codecov.io/gh/egymgmbh/mobile-update-check)
[![Go Report Card](https://goreportcard.com/badge/github.com/egymgmbh/mobile-update-check)](https://goreportcard.com/report/github.com/egymgmbh/mobile-update-check)
[![eGym Core Team](https://img.shields.io/badge/eGym-Core%20Team-orange.svg)](https://www.egym.com/)

The Mobile-Update-Check offers a simple way to check an apps version and
potentially request an update from the user if necessary.

Currently only product and os-version numbers of the following format are supported:
```
MAJOR.MINOR.PATCH
```

# Adding Update Rules
In order to add a rule for an application, one needs to add the rules to the

```
rules.json
```
The rules.json file contains a list of rules-sets. Each ruleset contains a list of rules for a key.
A key is the `os/product` combination e.g. `ios/fitapp` which also is the path used then for the GET request.

An example for the rules.json could look as follows:

```
[
  {
    "key": "ios/someapp",
    "rules": [
      {
        "osVersion": "=<5.0.1",
        "productVersion": "=<2.0.10",
        "action": "ADVICE"
      },
      {
        "osVersion": "=<5.0.1",
        "productVersion": "=<3.0.10",
        "action": "NONE"
      }
    ]
  },
  {
    "key": "android/someapp",
    "rules": [
      {
        "osVersion": "=<5.0.1",
        "productVersion": "=<2.0.10",
        "action": "ADVICE"
      }
    ]
  },
  {
    "key": "newos/newapp",
    "rules": [
      {
        "osVersion": "=<5.0.1",
        "productVersion": "=<2.0.10",
        "action": "FORCE"
      }
    ]
  }
]
```

## Version Numbers within the Update Rules
The version number parsing / validating is done with the following library:
[Masterminds SemVer project on Github](https://github.com/Masterminds/semver)

Basic Comparisons (taken from the project on GitHub):

```
There are two elements to the comparisons. First, a comparison string is a list of comma separated and comparisons.
These are then separated by || separated or comparisons. For example, ">= 1.2, < 3.0.0 || >= 4.2.3" is looking for
a comparison that's greater than or equal to 1.2 and less than 3.0.0 or is greater than or equal to 4.2.3.

=: equal (aliased to no operator)
!=: not equal
>: greater than
<: less than
>=: greater than or equal to
<=: less than or equal to
```
For further information and more complex cases, please check the documentation on Github:
[Masterminds SemVer Github](https://github.com/Masterminds/semver)

# Testing
Please write tests for new functionality as well as tests if edge cases occur
that haven't been tested yet. Add these tests to:
```
mobileupdatecheck_test.go
```

The following will run the application as is:
```
go run mobileupdatecheck.go
```
And you can test the locally running version e.g. with the following query in your browser:
```
localhost:8008/android/fitapp?osVersion=1.0.0&productVersion=1.2.3
```

# Deployment
The mobile-update-check application is deployed via our current egym internal
deployment tools that should be accessible by egym devs. If not please request
access from the core team.