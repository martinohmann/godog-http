godog-http
==========

[![Build Status](https://travis-ci.org/martinohmann/godog-http.svg)](https://travis-ci.org/martinohmann/godog-http)
[![codecov](https://codecov.io/gh/martinohmann/godog-http/branch/master/graph/badge.svg)](https://codecov.io/gh/martinohmann/godog-http)
[![Go Report Card](https://goreportcard.com/badge/github.com/martinohmann/godog-http)](https://goreportcard.com/report/github.com/martinohmann/godog-http)
[![GoDoc](https://godoc.org/github.com/martinohmann/godog-http?status.svg)](https://godoc.org/github.com/martinohmann/godog-http)

godog-http defines a godog feature context which adds steps to test `http.Handler` implementations.

Installation
------------

```sh
go get -u github.com/martinohmann/godog-http
```

Usage
-----

Example feature:

```
Feature: As a developer, I want to be able to setup and verify http routers in
  godog features.

  Scenario: I make a JSON request

    Given I have following request headers:
      | X-Auth | supersecret |
    And I have following request body:
      """
      {"foo":{"bar":{"baz":1}},"something":"else"}
      """
    When I send "POST" request to "/foo"
    Then the response code should be 200
    And the response should contain following json:
      """
      {"something":"else","foo":{"bar":{"baz":1}}}
      """
    And the response should have following headers:
      | Content-Type | application/json |
```

Check [`feature_context_test.go`](feature_context_test.go) and the
[`features/`](features/) directory for more usage examples.

License
-------

The source code of godog-http is released under the MIT License. See the bundled
LICENSE file for details.
