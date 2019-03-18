Feature: As a developer, I want to be able to setup and verify http routers in
  godog features.

  Scenario: Wrong request method

    When I send "GET" request to "/foo"
    Then the response code should be 405
    And the response should contain following json:
      """
      {"error":"Method Not Allowed"}
      """ 

  Scenario: Missing auth credentials

    When I send "POST" request to "/foo"
    Then the response code should be 401
    And the response should contain following json:
      """
      {"error":"Unauthorized"}
      """ 

  Scenario: Missing request body

    Given I have following request headers:
      | X-Auth | supersecret |
    When I send "POST" request to "/foo"
    Then the response code should be 400
    And the response should contain following json:
      """
      {"error":"Bad Request"}
      """ 

  Scenario: Invalid request body

    Given I have following request headers:
      | X-Auth | supersecret |
    And I have following request body:
      """
      some<garbage
      """
    When I send "POST" request to "/foo"
    Then the response code should be 400
    And the response should contain following json:
      """
      {"error":"Bad Request"}
      """ 

  Scenario: Valid request body, exact response

    Given I have following request headers:
      | X-Auth | supersecret |
    And I have following request body:
      """
      {"foo":{"bar":{"baz":1}},"something":"else"}
      """
    When I send "POST" request to "/foo"
    Then the response code should be 200
    And the response should be:
      """
      {"foo":{"bar":{"baz":1}},"something":"else"}

      """ 

  Scenario: Valid request body, JSON subtree match

    Given I have following request headers:
      | X-Auth | supersecret |
    And I have following request body:
      """
      {"foo":{"bar":{"baz":1}},"something":"else"}
      """
    When I send "POST" request to "/foo"
    Then the response code should be 200
    And the response should contain following json subtree:
      """
      {"foo":{"bar":{"baz":1}}}
      """ 

  Scenario: Valid request body, exact JSON response, ignoring map order

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

  Scenario: Valid request body, response regex matching

    Given I have following request headers:
      | X-Auth | supersecret |
    And I have following request body:
      """
      {"foo":{"bar":{"baz":1}},"something":"else"}
      """
    When I send "POST" request to "/foo"
    Then the response code should be 200
    And the response should match pattern:
      """
      ^\s*{.*}\s*$
      """ 

  Scenario: Response header matching

    Given I have following request headers:
      | X-Auth | supersecret |
    And I have following request body:
      """
      {"foo":1}
      """
    When I send "POST" request to "/foo"
    Then the response code should be 200
    And the response should have following headers:
      | Content-Type | application/json |
