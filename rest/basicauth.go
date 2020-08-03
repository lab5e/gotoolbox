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
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// CredentialStore is a simple credential store
type CredentialStore interface {
	Valid(username, password string) bool
}

// BasicAuthWrapper wraps all requests with WWW-Authenticate headers and
// validates against the credential store.
func BasicAuthWrapper(realm string, authStore CredentialStore, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=\"%s\"", realm))

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, "unable to decode authorization header", http.StatusUnauthorized)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) < 2 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if !authStore.Valid(pair[0], pair[1]) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		handlerFunc.ServeHTTP(w, r)
	}
}

// MemoryCredentialStore is a memory-backed credential store. Suitable for testing, nothing more.
type MemoryCredentialStore struct {
	username string
	password string
}

// NewMemoryCredentialStore creates a new MemoryCredentialStore instance with a single entry
func NewMemoryCredentialStore(username, password string) *MemoryCredentialStore {
	return &MemoryCredentialStore{username, password}
}

// Valid validates user name and password
func (m *MemoryCredentialStore) Valid(username, password string) bool {
	return (m.username == username && m.password == password)
}
