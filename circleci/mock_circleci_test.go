package circleci_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/onsi/ginkgo"
)

var (
	server *httptest.Server
)

func Host(server *httptest.Server) string {
	return server.URL
}

type MockRoute struct {
	Method      string
	Endpoint    string
	Output      string
	Status      int
	QueryString string
	PostForm    *string
}

func testQueryString(QueryString string, QueryStringExp string) {
	value, _ := url.QueryUnescape(QueryString)

	if QueryStringExp != value {
		defer ginkgo.GinkgoRecover()
		ginkgo.Fail(fmt.Sprintf("Error: Query string '%s' should be equal to '%s'", QueryStringExp, value))
	}
}

func testPostQuery(req *http.Request, postFormBody *string) {
	if postFormBody != nil {
		if body, err := ioutil.ReadAll(req.Body); err != nil {
			defer ginkgo.GinkgoRecover()
			ginkgo.Fail("No request body but expected one")
		} else {
			defer req.Body.Close()
			if strings.TrimSpace(string(body)) != strings.TrimSpace(*postFormBody) {
				defer ginkgo.GinkgoRecover()
				ginkgo.Fail(fmt.Sprintf("Expected POST body (%s) does not equal POST body (%s)", *postFormBody, body))
			}
		}
	}
}

func setup(mock MockRoute) {
	setupMultiple([]MockRoute{mock})
}

func setupMultiple(mockEndpoints []MockRoute) {
	router := mux.NewRouter()

	querySplitter := func(c rune) bool {
		return c == '&' || c == '='
	}

	for _, mock := range mockEndpoints {
		method := mock.Method
		endpoint := mock.Endpoint
		output := mock.Output
		status := mock.Status
		queryString := mock.QueryString
		postFormBody := mock.PostForm

		queries := strings.FieldsFunc(queryString, querySplitter)

		if method == "POST" {
			router.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
				testQueryString(r.URL.RawQuery, queryString)
				testPostQuery(r, postFormBody)
				w.WriteHeader(status)
				fmt.Fprintf(w, output)
			}).Methods(method)
		} else {
			router.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
				testQueryString(r.URL.RawQuery, queryString)
				w.WriteHeader(status)
				fmt.Fprintf(w, output)
			}).Methods(method).Queries(queries...)
		}
	}

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer ginkgo.GinkgoRecover()
		ginkgo.Fail(fmt.Sprintf("Route requested but not mocked: %s", r.URL))
	})

	server = httptest.NewServer(router)
}

func teardown() {
	if server != nil {
		server.Close()
		server = nil
	}
}
