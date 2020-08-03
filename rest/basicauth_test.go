package rest

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func tempHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
func TestBasicAuth(t *testing.T) {
	credStore := NewMemoryCredentialStore("john", "doe")

	server := httptest.NewServer(BasicAuthWrapper("foo", credStore, tempHandler))

	defer server.Close()

	c := http.Client{}

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal("Error creating request: ", err)
	}
	//	req.SetBasicAuth("", "")
	resp, _ := c.Do(req)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatal("Got ", resp.StatusCode, " with not username or password")
	}
	req.Header.Set("Authorization", "xxx&! !!!!")
	resp, _ = c.Do(req)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatal("Got ", resp.StatusCode, " with garbled chars")
	}

	req.Header.Set("Authorization", fmt.Sprintf("something %s", base64.StdEncoding.EncodeToString([]byte("justusername"))))
	resp, _ = c.Do(req)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatal("Got ", resp.StatusCode, " with garbled chars")
	}

	req.SetBasicAuth("john", "incorrect")
	resp, _ = c.Do(req)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatal("Got ", resp.StatusCode, " with incorrect password")
	}

	req.SetBasicAuth("john", "doe")
	resp, _ = c.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Did not expect error")
	}
}
