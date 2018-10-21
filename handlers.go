package GoLib

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"os"
	"regexp"
	"runtime"

	"fmt"
	"time"

	"reflect"

	"github.com/Vilsol/GoLib/logger"
	"github.com/gorilla/mux"
)

type GenericHandle func(http.ResponseWriter, *http.Request) GenericResponse

func ProcessResponse(handle GenericHandle) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		response := handle(w, r)

		format := r.URL.Query().Get("format")
		pretty := len(r.URL.Query()["pretty"]) > 0

		switch format {
		default:
			w.Header().Set("Content-Type", "application/json")
		case "xml":
			w.Header().Set("Content-Type", "application/xml")
		}

		if response.Error != nil {
			w.WriteHeader(response.Error.Status)
		}

		switch format {
		default:
			encoder := json.NewEncoder(w)

			if pretty {
				encoder.SetIndent("", "    ")
			}

			encoder.Encode(response)
		case "xml":
			encoder := xml.NewEncoder(w)

			if pretty {
				encoder.Indent("", "    ")
			}

			encoder.Encode(response)
		}
	})
}

var panicRegex = regexp.MustCompile("(?m)^panic\\(.+?\\)")

type DataHandle func(*http.Request) (interface{}, *ErrorResponse)

func DataHandler(handle DataHandle) GenericHandle {
	return func(w http.ResponseWriter, r *http.Request) GenericResponse {
		var data interface{}
		var err *ErrorResponse

		Block{
			Try: func() {
				data, err = handle(r)
			},
			Catch: func(e Exception) {
				switch reflect.TypeOf(e).Name() {
				case "InvalidParameter":
					err = &ErrorInvalidParameter
					err.Message = "Invalid Parameter: " + (e.(InvalidParameter).Parameter)
					break
				case "string":
					err = &ErrorInternalError
					err.Message += e.(string)
				case "errorString":
					err = &ErrorInternalError
					err.Message += e.(error).Error()
				default:
					err = &ErrorInternalError
					break
				}

				fmt.Fprintln(os.Stderr, err.Message)

				buf := make([]byte, 1<<10)
				runtime.Stack(buf, false)
				panicIndex := panicRegex.FindIndex(buf)
				if len(panicIndex) == 0 {
					fmt.Fprintln(os.Stderr, string(buf))
				} else {
					fmt.Fprintln(os.Stderr, string(buf)[0:bytes.IndexByte(buf, '\n') + 1] + string(buf)[panicIndex[0]:])
				}
			},
		}.Do()

		return GenericResponse{
			Success: err == nil,
			Error:   err,
			Data:    data,
		}
	}
}

type RegisterRoute func(method string, path string, handle DataHandle)

func RouteHandler(router *mux.Router, prefix string) RegisterRoute {
	return func(method string, path string, handle DataHandle) {
		route := router.NewRoute()
		route.PathPrefix(prefix)
		route.Path(path)
		route.Methods(method)
		route.Handler(ProcessResponse(DataHandler(handle)))
	}
}

func LoggerHandler(h http.Handler) http.Handler {
	l := logger.New("API")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hijackedLogger := makeLogger(w)
		start := time.Now()

		h.ServeHTTP(hijackedLogger, r)

		end := time.Now()
		latency := end.Sub(start)

		clientIP := r.RemoteAddr
		method := r.Method

		statusCode := hijackedLogger.Status()

		statusColor := colorForStatus(statusCode)
		methodColor := colorForMethod(method)

		path := r.URL.Path

		l.Log("%s %3d %s| %12v | %21s |%s %-7s %s %s",
			statusColor, statusCode, logger.Reset,
			latency,
			clientIP,
			methodColor, method, logger.Reset,
			path,
		)
	})
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return logger.Green
	case code >= 300 && code < 400:
		return logger.Magenta
	case code >= 400 && code < 500:
		return logger.Yellow
	default:
		return logger.Red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return logger.Blue
	case "POST":
		return logger.Cyan
	case "PUT":
		return logger.Yellow
	case "DELETE":
		return logger.Red
	case "PATCH":
		return logger.Green
	case "HEAD":
		return logger.Magenta
	case "OPTIONS":
		return logger.White
	default:
		return logger.Reset
	}
}

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprint(w, "404 page not found")
	})
}
