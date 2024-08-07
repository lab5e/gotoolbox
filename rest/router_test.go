package rest

//
//Copyright 2018 Telenor Digital AS
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test simple routes
func TestRouting(t *testing.T) {
	matchingRoutes := []string{
		"/foo/first",
		"/foo/second",
		"/bar/first",
		"/bar/second",
		"/baz",
		"/baz/foo/bar",
		"/baz/foo/bar?param=value&value=param",
	}
	nonMatchingRoutes := []string{
		"/foo",
		"/foo/other/some",
		"/baz/bar",
	}

	invocationCount := 0

	routeHandler := func(w http.ResponseWriter, r *http.Request) {
		invocationCount++
	}

	router := ParameterRouter{}
	router.AddRoute("/foo/{arg}", routeHandler)
	router.AddRoute("/bar/{arg}", routeHandler)
	router.AddRoute("/baz", routeHandler)
	router.AddRoute("/baz/foo/bar", routeHandler)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, route := range matchingRoutes {
		handler := router.GetHandler(route)
		if handler != nil {
			handler.ServeHTTP(w, req)
		}
	}

	for _, route := range nonMatchingRoutes {
		handler := router.GetHandler(route)
		if handler != nil {
			handler.ServeHTTP(w, req)
		}
	}

	if invocationCount != len(matchingRoutes) {
		t.Errorf("Did not get the expected number of matches. Got %d expected %d.", invocationCount, len(matchingRoutes))
	}
}

func TestPathKeys(t *testing.T) {
	r := httptest.NewRequest("GET", "/hello", strings.NewReader("{}"))

	if GetPathKey("param",
		r.WithContext(
			context.WithValue(
				r.Context(), PathParameter("param"), "value"))) != "value" {
		t.Fatal("Did not get the value I expected")
	}
	if GetPathKey("param", nil) != "" {
		t.Fatal("Expected empty string for nil request")
	}
	if GetPathKey("param", r) != "" {
		t.Fatal("Expected empty string with no params in context")
	}
	if GetPathKey("param",
		r.WithContext(
			context.WithValue(
				r.Context(), PathParameter("param"), 12))) != "" {
		t.Fatal("Expected empty string for non-string values")
	}
}

// A benchmark that both adds and routes
func BenchmarkRouting(b *testing.B) {
	rng := rand.New(rand.NewSource(42))
	const routeCount int = 50
	router := ParameterRouter{}
	testHandler := func(w http.ResponseWriter, r *http.Request) { /* empty */ }
	for i := 0; i < routeCount; i++ {
		route := fmt.Sprintf("/some/{arg1}/%d/{arg2}", i)
		router.AddRoute(route, testHandler)
	}

	for i := 0; i < b.N; i++ {
		randomRoute := fmt.Sprintf("/some/%d/%d/%d", rng.Intn(routeCount), rng.Intn(routeCount), rng.Intn(routeCount))
		handler := router.GetHandler(randomRoute)
		if handler == nil {
			b.Error("Did not expect nil response")
		}
	}
}

// The number of routes in the router
const routeCount int = 50

// Router for the benchmark test
var brouter ParameterRouter

// Routes to test - using a fixed number
var routesToTest []string

func init() {
	brouter = ParameterRouter{}

	rng := rand.New(rand.NewSource(42))
	testHandler := func(w http.ResponseWriter, r *http.Request) { /* empty */ }
	for i := 0; i < routeCount; i++ {
		route := fmt.Sprintf("/some/{arg1}/%d/{arg2}", i)
		brouter.AddRoute(route, testHandler)
	}
	for i := 0; i < routeCount; i++ {
		routesToTest = append(routesToTest, fmt.Sprintf("/some/%d/%d/%d", rng.Intn(routeCount), i, rng.Intn(routeCount)))
		routesToTest = append(routesToTest, "/not/matching/route")
	}
}

// Test just the routing request; set up isn't very critical
func BenchmarkJustRouting(b *testing.B) {
	for i := 0; i < b.N; i++ {
		handler := brouter.GetHandler(routesToTest[i%routeCount])
		if handler == nil {
			b.Error("Did not expect nil response")
		}
	}
}
