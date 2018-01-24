package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mholt/archiver"
	"github.com/rjz/forager"
	"gopkg.in/go-playground/validator.v9"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var RootDir = "."
var validate *validator.Validate
var templates *template.Template

func jsonError(w http.ResponseWriter, code int, error string) {
	bytes, _ := json.Marshal(struct {
		Error string `json:"error"`
		Code  int    `json:"code"`
	}{error, code})

	w.WriteHeader(code)
	w.Write(bytes)
}

func jsonSimpleError(w http.ResponseWriter, code int) {
	jsonError(w, code, http.StatusText(code))
}

func serve(port string) {

	validate = validator.New()

	router := mux.NewRouter()

	router.HandleFunc("/maps", func(w http.ResponseWriter, r *http.Request) {
		opts := forager.MapOpts{}

		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = json.Unmarshal(bytes, &opts)
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = validate.Struct(opts)
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		dir := os.TempDir() + "/forager"
		ctx := context.Background()

		if err := forager.RenderMap(ctx, opts, dir); err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = archiver.Zip.Write(w.(http.ResponseWriter), []string{dir})
		if err != nil {
			jsonSimpleError(w, http.StatusInternalServerError)
			return
		}
	}).Methods("POST")

	router.PathPrefix("/assets").Handler(http.StripPrefix("/assets", http.FileServer(http.Dir("./client")))).Methods("GET")

	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"GoogleMapsClientKey": os.Getenv("GOOGLE_MAPS_CLIENT_KEY"),
		}

		log.Println(data)

		templates.ExecuteTemplate(w, "index.html.tmpl", data)
	})

	log.Printf("Ready to serve @ %s\n", port)

	log.Fatal(http.ListenAndServe(port, router))
}

func main() {
	portNum := os.Getenv("PORT")
	serve(":" + portNum)
}

func init() {
	pattern := filepath.Join(RootDir, "templates", "*.html.tmpl")
	templates = template.Must(template.ParseGlob(pattern))
}
