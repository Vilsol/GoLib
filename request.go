package GoLib

import (
	"net/http"
	"strconv"
)

func Remap(input string, maps map[string]string) string {
	if val, ok := maps[input]; ok {
		return val
	}

	return input
}

type InvalidParameter struct {
	Exception

	Parameter string
}

func GetString(r *http.Request, name string, def string) string {
	if len(r.URL.Query()[name]) == 0 {
		return def
	}

	return r.URL.Query().Get(name)
}

func GetInt(r *http.Request, name string, def int) int {
	if len(r.URL.Query()[name]) == 0 {
		return def
	}

	i, err := strconv.Atoi(r.URL.Query().Get(name))

	if err != nil {
		Throw(InvalidParameter{nil, name})
	}

	return i
}

func GetIntMinMax(r *http.Request, name string, def int, min int, max int) int {
	i := GetInt(r, name, def)

	if i < min {
		i = min
	}

	if i > max {
		i = max
	}

	return i
}

func GetBool(r *http.Request, name string, def bool) bool {
	if len(r.URL.Query()[name]) == 0 {
		return def
	}

	i, err := strconv.ParseBool(r.URL.Query().Get(name))

	if err != nil {
		Throw(InvalidParameter{nil, name})
	}

	return i
}
