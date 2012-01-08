package goweb

import (
	"net/http"
	"regexp"
)

type Any interface{}
type Responder func(http.ResponseWriter, *http.Request, Result, ...string) Result
type Result struct {
	Final bool
	State map[string]Any
}

// Or composes Routers by returning the first final Result of the Routers.
func Or(rs ...Responder) Responder {
	return func(w http.ResponseWriter, c *http.Request, s Result, _ ...string) Result {
		for _, r := range rs {
			s = r(w, c, s)
			if s.Final {
				return s
			}
		}
		return s
	}
}

// Route returns a new router that matches the pat against
// the url that is passed in.
func Route(pat string, f Responder) Responder {
	r, err := regexp.Compile("^" + pat + "$")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, c *http.Request, s Result, _ ...string) Result {
		path := MatchUrl(r, c.URL.Path)
		if path != nil {
			s.Final = true
			f(w, c, s, path...)
			return s
		}
		s.Final = false
		return s
	}
}

// MatchUrl returns the list of submatches of r
func MatchUrl(r *regexp.Regexp, against string) []string {
	if x := r.FindStringSubmatch(against); x != nil {
		return x[1:]
	}
	return nil
}

// Pass sets s.Final to false
func Pass(s *Result) {
	s.Final = false
}
func Handler(route Responder) http.HandlerFunc {
	return func(w http.ResponseWriter, c *http.Request) {
		s := Result{Final: false, State: make(map[string]Any)}
		route(w, c, s)
	}
}
