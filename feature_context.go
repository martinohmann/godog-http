// Package http defines a godog feature context which adds steps to test
// http.Handler implementations.
package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/martinohmann/jsoncompare"
)

// FeatureContext adds steps to setup and verify http handlers in godog tests.
type FeatureContext struct {
	handler http.Handler
	resp    *httptest.ResponseRecorder
	body    io.Reader
	header  http.Header
}

// NewFeatureContext creates a new FeatureContext. It expects a http.Handler
// which the test suite can make requests against.
func NewFeatureContext(handler http.Handler) *FeatureContext {
	return &FeatureContext{handler: handler}
}

// beforeScenario is called before each scenario and resets the http.Request
// parameters and the httptest.ResponseRecorder.
func (c *FeatureContext) beforeScenario(interface{}) {
	c.resp = httptest.NewRecorder()
	c.body = nil
	c.header = nil
}

// iHaveFollowingRequestHeaders sets headers for the next http.Request.
func (c *FeatureContext) iHaveFollowingRequestHeaders(header *gherkin.DataTable) error {
	for _, h := range header.Rows {
		if len(h.Cells) != 2 {
			return fmt.Errorf(
				"expected header table row to have two columns (name, value), got %d",
				len(h.Cells),
			)
		}

		if c.header == nil {
			c.header = make(http.Header)
		}

		name, value := h.Cells[0].Value, h.Cells[1].Value

		c.header[name] = []string{value}
	}

	return nil
}

// iHaveFollowingRequestBody sets the body of the next http.Request.
func (c *FeatureContext) iHaveFollowingRequestBody(body *gherkin.DocString) error {
	c.body = bytes.NewBuffer([]byte(body.Content))

	return nil
}

// iSendRequestTo sends a http.Request to the handler's ServeHTTP method. It
// also passes a httptest.ResponseRecorder to the handler which is used to make
// assertions about the response of the request in later steps.
func (c *FeatureContext) iSendRequestTo(method string, url string) error {
	req, err := http.NewRequest(method, url, c.body)
	if err != nil {
		return err
	}

	if c.header != nil {
		req.Header = c.header
	}

	c.handler.ServeHTTP(c.resp, req)

	return nil
}

// theResponseCodeShouldBe asserts that the last response has given status code.
func (c *FeatureContext) theResponseCodeShouldBe(code int) error {
	if c.resp.Code != code {
		return fmt.Errorf("expected response code to be %d, got %d", code, c.resp.Code)
	}

	return nil
}

// theResponseShouldBe asserts that the response body exactly matches the
// expectations.
func (c *FeatureContext) theResponseShouldBe(body *gherkin.DocString) error {
	expected := body.Content
	actual := c.resp.Body.String()

	if actual != expected {
		return fmt.Errorf("expected response body to be %q, got %q", expected, actual)
	}

	return nil
}

// theResponseShouldMatchPattern asserts that the response body matches a
// certain regexp pattern.
func (c *FeatureContext) theResponseShouldMatchPattern(pattern *gherkin.DocString) error {
	r := regexp.MustCompile(pattern.Content)
	body := c.resp.Body.Bytes()

	if !r.Match(body) {
		return fmt.Errorf(
			"expected response body %q to match pattern %q, but it did not",
			string(body),
			pattern.Content,
		)
	}

	return nil
}

// theResponseShouldContainFollowingJSON asserts that the response matches the
// expected JSON. Use this if you also want to validate the JSON and you also
// do not care about the order of map keys.
func (c *FeatureContext) theResponseShouldContainFollowingJSON(body *gherkin.DocString) error {
	return c.compareResponseJSON([]byte(body.Content), jsoncompare.MatchStrict)
}

// theResponseShouldContainFollowingJSONSubtree asserts that the response
// contains the given JSON. Use this if there can be assitional keys in the
// response JSON you do not care about.
func (c *FeatureContext) theResponseShouldContainFollowingJSONSubtree(body *gherkin.DocString) error {
	return c.compareResponseJSON([]byte(body.Content), jsoncompare.MatchSubtree)
}

// compareResponseJSON compares an expected byte slice with the response body
// using a jsoncompare.Comparator with given matchMode.
func (c *FeatureContext) compareResponseJSON(expected []byte, matchMode jsoncompare.MatchMode) error {
	comparator := jsoncompare.NewComparator(matchMode)

	return comparator.Compare(c.resp.Body.Bytes(), expected)
}

// theResponseShouldHaveFollowingHeaders asserts that given headers are present
// in the response and that their values match the expectations.
func (c *FeatureContext) theResponseShouldHaveFollowingHeaders(headers *gherkin.DataTable) error {
	for _, h := range headers.Rows {
		if len(h.Cells) != 2 {
			return fmt.Errorf(
				"expected header table row to have two columns (name, value), got %d",
				len(h.Cells),
			)
		}

		name, value := h.Cells[0].Value, h.Cells[1].Value

		if header, ok := c.resp.Header()[name]; !ok || len(header) == 0 {
			return fmt.Errorf("header %q missing", name)
		} else if header[0] != value {
			return fmt.Errorf("expected value %q for header %q, got %q", value, name, header[0])
		}
	}

	return nil
}

// Register registers the feature context to the godog suite.
func (c *FeatureContext) Register(s *godog.Suite) {
	s.BeforeScenario(c.beforeScenario)

	// Given
	s.Step(`^I have following request headers:$`, c.iHaveFollowingRequestHeaders)
	s.Step(`^I have following request body:$`, c.iHaveFollowingRequestBody)

	// When
	s.Step(`^I send "(OPTIONS|GET|HEAD|POST|PUT|DELETE|TRACE|CONNECT)" request to "([^"]*)"$`, c.iSendRequestTo)

	// Then
	s.Step(`^the response code should be (\d+)$`, c.theResponseCodeShouldBe)
	s.Step(`^the response should be:$`, c.theResponseShouldBe)
	s.Step(`^the response should match pattern:$`, c.theResponseShouldMatchPattern)
	s.Step(`^the response should contain following json:$`, c.theResponseShouldContainFollowingJSON)
	s.Step(`^the response should contain following json subtree:$`, c.theResponseShouldContainFollowingJSONSubtree)
	s.Step(`^the response should have following headers:$`, c.theResponseShouldHaveFollowingHeaders)
}
