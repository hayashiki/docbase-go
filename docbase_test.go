package docbase

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	url, _ := url.Parse(server.URL)

	client = NewClient(nil, "dummyTeam", "dummyToken", OptionDocbaseURL(url))
}

func teardown() {
	defer server.Close()
}

func TestNewClient(t *testing.T) {
	cli := NewClient(nil, "fakeTeam", "fakeToken")

	if got, want := cli.BaseURL.String(), "https://api.docbase.io/teams/fakeTeam"; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
}

func TestClient_NewRequest(t *testing.T) {
	cli := NewClient(nil, "fakeTeam", "fakeToken")

	method := "GET"
	inURL, outURL := "foo", "https://api.docbase.io/teams/fakeTeam/foo"
	inBody := struct{ Foo string }{Foo: "Bar"}

	outBody := `{"Foo":"Bar"}`

	req, err := cli.NewRequest(method, inURL, inBody)

	if err != nil {
		t.Errorf("err")
	}

	if got, want := req.Method, method; got != want {
		t.Errorf("NewRequest(%q) Method is %v, want %v", method, got, want)
	}

	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	body, _ := ioutil.ReadAll(req.Body)

	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%q) Body is %v, want %v", inBody, got, want)
	}
}

func TestDo_httpError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req, nil)

	_, ok := err.(*ErrorResponse)
	if !ok {
		t.Errorf("Error should be of type ErrorResponse but is %v: %+v", reflect.TypeOf(err), err)
	}

	if got, want := err.Error(), "Bad Request"; got != want {
		t.Errorf("Request Do %v, want %v", got, want)
	}
}

