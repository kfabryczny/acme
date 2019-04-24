package acme

import (
	"net/http"
	"testing"

	"reflect"
)

func TestNewClient(t *testing.T) {
	if _, err := NewClient("http://fake"); err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestFetchLink(t *testing.T) {
	linkTests := []struct {
		Name        string
		LinkHeaders []string
		WantedLink  string
		ExpectedURL string
	}{
		{
			Name:        "no links",
			WantedLink:  "fail",
			ExpectedURL: "",
		},
		{Name: "joined links",
			LinkHeaders: []string{`<https://url/path>; rel="next", <http://url/path?query>; rel="up"`},
			WantedLink:  "up",
			ExpectedURL: "http://url/path?query",
		},
		{
			Name:        "separate links",
			LinkHeaders: []string{`<https://url/path>; rel="next"`, `<http://url/path?query>; rel="up"`},
			WantedLink:  "up",
			ExpectedURL: "http://url/path?query",
		},
	}
	for _, currentTest := range linkTests {
		linkURL := fetchLink(&http.Response{Header: http.Header{"Link": currentTest.LinkHeaders}}, currentTest.WantedLink)
		if linkURL != currentTest.ExpectedURL {
			t.Fatalf("%s: links not equal, expected: %s, got: %s", currentTest.Name, currentTest.ExpectedURL, linkURL)
		}
	}
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestFetchLinks(t *testing.T) {
	linkTests := []struct {
		Name         string
		LinkHeaders  []string
		WantedLink   string
		ExpectedURLs []string
	}{
		{
			Name:         "no links",
			WantedLink:   "fail",
			ExpectedURLs: nil,
		},
		{Name: "joined links",
			LinkHeaders:  []string{`<https://url/path>; rel="next", <http://url/path?query>; rel="up"`},
			WantedLink:   "up",
			ExpectedURLs: []string{"http://url/path?query"},
		},
		{
			Name:         "separate links",
			LinkHeaders:  []string{`<https://url/path>; rel="next"`, `<http://url/path?query>; rel="up"`},
			WantedLink:   "up",
			ExpectedURLs: []string{"http://url/path?query"},
		},
		{
			Name:         "multiple links",
			LinkHeaders:  []string{`<https://url/path>; rel="up"`, `<http://url/path?query>; rel="up"`},
			WantedLink:   "up",
			ExpectedURLs: []string{"https://url/path", "http://url/path?query"},
		},
	}
	for _, currentTest := range linkTests {
		linkURLs := fetchLinks(&http.Response{Header: http.Header{"Link": currentTest.LinkHeaders}}, currentTest.WantedLink)
		if !stringSliceEqual(linkURLs, currentTest.ExpectedURLs) {
			t.Fatalf("%s: links not equal, expected: %s, got: %s", currentTest.Name, currentTest.ExpectedURLs, linkURLs)
		}
	}
}

func TestClient_Directory(t *testing.T) {
	if !reflect.DeepEqual(testClient.dir, testClient.Directory()) {
		t.Fatalf("directory mismatch, expected: %+v, got: %+v", testClient.dir, testClient.Directory())
	}
}
