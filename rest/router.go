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
	"net/http"
	"regexp"
	"strings"

	"github.com/ExploratoryEngineering/logging"
)

// PathParameter is the type used when storing parameters in the context
type PathParameter string

// The elementwalker function is used when walking through the different parts of the path.
type elementwalker func(pattern string, path string, param map[string]string) bool

// Pull a value (aka parameter) from the path
func getparameter(pattern string, path string, param map[string]string) bool {
	param[pattern] = path
	return true
}

// Do a simple comparison between path elements
func compare(pattern string, path string, param map[string]string) bool {
	return pattern == path
}

// A single route
type route struct {
	elements []string
	process  []elementwalker
	handler  http.HandlerFunc
}

// Match an URI. The URI is split into an array of strings (ie elements of path). Returns
// a list of parameters and a bool
func (r *route) matches(uriElements []string) (map[string]string, bool) {
	if len(uriElements) != len(r.elements) {
		return nil, false
	}
	params := make(map[string]string)
	for i, pathElement := range r.elements {
		if !r.process[i](pathElement, uriElements[i], params) {
			return nil, false
		}
	}
	return params, true
}

// ParameterRouter implements a request URI router that allows for path parameters in
// URIs, ie requests like https://example.com/api/thing/0424242/subthing/ - ie presenting
// resources with IDs in the request URI rather than having to rely on query parameters.
// The router can be plugged in in the standard http package.
type ParameterRouter struct {
	routes []route
}

// NewParameterRouter creates a new router instance
func NewParameterRouter() ParameterRouter {
	return ParameterRouter{routes: make([]route, 0)}
}

// AddRoute adds a new route described by the specified pattern, handled by the supplied ParameterHandler. This method isn't thread safe.
func (r *ParameterRouter) AddRoute(pattern string, handler http.HandlerFunc) {
	patternElements := strings.Split(pattern, "/")
	newRoute := route{
		elements: patternElements,
		handler:  handler,
	}

	const parameterPattern string = "{(.*)}"
	re, err := regexp.Compile(parameterPattern)
	if err != nil {
		panic(fmt.Sprintf("Unable to compile regexp for pattern: %s", err))
	}

	for i, element := range patternElements {
		match := re.FindStringSubmatch(element)
		var procFunc elementwalker
		if len(match) > 0 {
			// Matching element - add to parameter http.ListenAndServe(":8080", nil)
			newRoute.elements[i] = match[1]
			procFunc = getparameter
		} else {
			newRoute.elements[i] = element
			procFunc = compare
		}
		newRoute.process = append(newRoute.process, procFunc)
	}
	r.routes = append(r.routes, newRoute)
}

// GetHandler returns a matching handler (as a closure) for the matching uri. If no matching route
// is found it will return the default http.NotFound handler.
func (r *ParameterRouter) GetHandler(uri string) http.HandlerFunc {
	params := strings.Split(uri, "?")
	pathElements := strings.Split(params[0], "/")
	for _, route := range r.routes {
		params, matches := route.matches(pathElements)
		if matches {
			return func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				for key, value := range params {
					ctx = context.WithValue(ctx, PathParameter(key), value)

				}
				route.handler(w, r.WithContext(ctx))
			}
		}
	}
	logging.Debug("No matching handler for %s", uri)
	return http.NotFound
}

// GetPathKey is a simple utility function that will return the (string) value
// for the specified key in the path. If there's an error it will return an
// empty string.
func GetPathKey(name string, r *http.Request) string {
	if r == nil {
		return ""
	}
	v := r.Context().Value(PathParameter(name))
	if v == nil {
		return ""
	}
	stringVal, ok := v.(string)
	if !ok {
		return ""
	}
	return stringVal
}
