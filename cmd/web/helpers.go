package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, tempData *templateData) {
	temp, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	tempData = app.addDefaultData(tempData, r)

	// Write the template to the buffer, instead of straight to the http.ResponseWriter.
	err := temp.Execute(buf, tempData)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Write the contents of the buffer to the http.ResponseWriter.
	buf.WriteTo(w)
}

func (app *application) addDefaultData(tempData *templateData, r *http.Request) *templateData {
	if tempData == nil {
		tempData = &templateData{}
	}
	tempData.CurrentYear = time.Now().Year()
	tempData.Flash = app.session.PopString(r, "flash")
	return tempData
}
