package docbase

import (
	"fmt"
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

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	if got := r.Header.Get(header); got != want {
		t.Errorf("NewRequest() %s header is %v, want %v", header, got, want)
	}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func TestNewClient(t *testing.T) {
	cli := NewClient(nil, "fakeTeam", "fakeToken")

	if got, want := cli.BaseURL.String(), "https://api.docbase.io/teams/fakeTeam"; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}

	if got, want := cli.Client, http.DefaultClient; got != want {
		t.Errorf("NewClient Client is %v, want %v", got, want)
	}
}

func TestClient_Do(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/sample", func(w http.ResponseWriter, r *http.Request) {

		testMethod(t, r, "POST")
		http.Error(w, `{
		    "code": 400,
		    "status": "bad request"
		}`, 400)
	})

	req, _ := client.NewRequest(http.MethodPost, "/sample", nil)

	_, err := client.Do(req, nil)
	if err == nil {
		t.Error("Do returns with expected error")
	}
}

func TestClient_NewRequest(t *testing.T) {
	cli := NewClient(nil, "fakeTeam", "fakeToken")

	method := "POST"
	inURL, outURL := "/foo", "https://api.docbase.io/teams/fakeTeam/foo"
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

	testHeader(t, req, "Content-Type", "application/json")
	testHeader(t, req, "Accept", "application/json")
	testHeader(t, req, "X-DocBaseToken", cli.AccessToken)
	testHeader(t, req, "X-Api-Version", apiVersion)
	testHeader(t, req, "USER_AGENT", userAgent)
}

func TestDo_httpError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{
    "error": "bad_request",
    "messages": [
        "Nameを入力してください"
    ]
}`)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req, nil)

	errResp, ok := err.(*ErrorResponse)
	if !ok {
		t.Errorf("Error should be of type ErrorResponse but is %v: %+v", reflect.TypeOf(err), err)
	}

	if got, want := errResp.ErrorStr, "bad_request"; got != want {
		t.Errorf("Request Do %v, want %v", got, want)
	}
}
