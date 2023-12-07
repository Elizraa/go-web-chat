package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"runtime"

	"github.com/elizraa/Globes/data"
)

var logger *log.Logger

/* Convenience function for printing to stdout
func p(a ...interface{}) {
	fmt.Println(a...)
}*/

// Info will log information with "INFO" prefix to logger
func Info(args ...interface{}) {
	logger.SetPrefix("INFO ")
	logger.Println(args...)
}

// Danger will log information with "ERROR" prefix to logger
func Danger(args ...interface{}) {
	// Retrieve caller details
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "??"
		line = 0
	}
	callingFunc := runtime.FuncForPC(pc).Name()

	// Log the error message along with caller information
	logger.Printf("ERROR [%s:%d] %s: %v\n", file, line, callingFunc, args)
}

// Warning will log information with "WARNING" prefix to logger
func Warning(args ...interface{}) {
	logger.SetPrefix("WARNING ")
	// Retrieve caller details
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "??"
		line = 0
	}
	callingFunc := runtime.FuncForPC(pc).Name()

	// Log the error message along with caller information
	logger.Printf("WARNING [%s:%d] %s: %v\n", file, line, callingFunc, args)
}

// ReportStatus is a helper function to return a JSON response indicating outcome success/failure
func ReportStatus(w http.ResponseWriter, success bool, err *data.APIError, ctxId string) {
	var res *data.Outcome
	w.Header().Set("Content-Type", "application/json")
	if success {
		res = &data.Outcome{
			Status:      success,
			API_CALL_ID: ctxId,
			Data:        map[string]interface{}{},
		}
	} else {
		res = &data.Outcome{
			Status:      success,
			Error:       err,
			API_CALL_ID: ctxId,
			Data:        map[string]interface{}{},
		}
	}
	response, _ := json.Marshal(res)
	if _, err := w.Write(response); err != nil {
		Danger("Error writing", response)
	}
}

func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	templates := template.Must(template.ParseFiles(files...))
	if err := templates.ExecuteTemplate(writer, "layout", data); err != nil {
		Danger("Error generating HTML template", data, err.Error())
	}
}

func generateHTMLContent(writer http.ResponseWriter, data interface{}, file string) {
	writer.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles(fmt.Sprintf("templates/content/%s.html", file))
	if err := t.Execute(writer, data); err != nil {
		Danger("Error executing HTML template", data, err.Error())
	}
}

// convenience function to be chained with another HandlerFunc
// that prints to the console which handler was called.
func logConsole(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
		fmt.Println("Handler function called - " + name)
		h(w, r)
	}
}
