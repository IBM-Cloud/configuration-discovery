// @APIVersion 1.0.0
// @APITitle Swagger IBM Cloud Provider API
// @APIDescription Swagger IBM Cloud Provider API
// @Contact sakshiag@in.ibm.com

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/IBM-Cloud/configuration-discovery/service"
	"github.com/fvbock/endless"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
)

var staticContent = flag.String("staticPath", "./swagger/swagger-ui", "Path to folder with Swagger UI")

// IndexHandler ..
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	isJsonRequest := false

	if acceptHeaders, ok := r.Header["Accept"]; ok {
		for _, acceptHeader := range acceptHeaders {
			if strings.Contains(acceptHeader, "json") {
				isJsonRequest = true
				break
			}
		}
	}

	if isJsonRequest {
		w.Write([]byte(resourceListingJson))
	} else {
		http.Redirect(w, r, "/swagger-ui/", http.StatusFound)
	}
}

// APIDescriptionHandler ..
func APIDescriptionHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := strings.Trim(r.RequestURI, "/")

	if apiKey == "v1" {
		if json, ok := apiDescriptionsJson[apiKey]; ok {
			w.Write([]byte(json))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		if json, ok := apiDescriptionsJsonV2[apiKey]; ok {
			w.Write([]byte(json))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config := service.GetConfiguration()

	//session, err := mgo.Dial(fmt.Sprintf("%s:%s@%s:%d", config.Mongo.UserName, config.Mongo.Password, config.Mongo.Host, config.Mongo.Port))
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatalln("Could not create mongo db session", err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	ensureIndex(session)

	var port int
	flag.IntVar(&port, "p", 8080, "Port on which this server listens")
	flag.Parse()
	r := mux.NewRouter()

	r.HandleFunc("/", IndexHandler)

	r.PathPrefix("/swagger-ui").Handler(http.StripPrefix("/swagger-ui", http.FileServer(http.Dir(*staticContent))))

	for apiKey := range apiDescriptionsJson {
		log.Println("API :", apiKey)
		r.HandleFunc("/"+apiKey, APIDescriptionHandler)
	}

	for apiKey := range apiDescriptionsJsonV2 {
		log.Println("API :", apiKey)
		r.HandleFunc("/"+apiKey, APIDescriptionHandler)
	}

	r.HandleFunc("/v1/configuration", service.ConfHandler(session)).Methods("POST")

	r.HandleFunc("/v1/configuration/{repo_name}", service.ConfDeleteHandler).Methods("DELETE")

	r.HandleFunc("/v1/configuration/{repo_name}/plan", service.PlanHandler(session)).Methods("POST")

	r.HandleFunc("/v1/configuration/{repo_name}/show", service.ShowHandler(session)).Methods("POST")

	r.HandleFunc("/v1/configuration/{repo_name}/apply", service.ApplyHandler(session)).Methods("POST")

	r.HandleFunc("/v1/configuration/{repo_name}/destroy", service.DestroyHandler(session)).Methods("POST")

	r.HandleFunc("/v1/configuration/{repo_name}/{action}/{actionID}/log", service.LogHandler).Methods("GET")

	r.HandleFunc("/v1/configuration/{repo_name}/{action}/{actionID}/status", service.StatusHandler(session)).Methods("GET")

	r.HandleFunc("/v1/configuration/{repo_name}/{action}/{log_file}", service.ViewLogHandler)

	r.HandleFunc("/v1/configuration/{repo_name}/{action}", service.GetActionDetailsHandler(session)).Methods("GET")

	r.HandleFunc("/v2/configuration", service.ConfHandler(session)).Methods("POST")

	r.HandleFunc("/v2/configuration/{repo_name}/import", service.TerraformerImportHandler(session)).Methods("GET")

	r.HandleFunc("/v2/configuration/{repo_name}/{action}/{actionID}/log", service.LogHandler).Methods("GET")

	r.HandleFunc("/v2/configuration/{repo_name}/{action}/{actionID}/status", service.StatusHandler(session)).Methods("GET")

	r.HandleFunc("/v2/configuration/{repo_name}/statefile", service.TerraformerStateHandler(session)).Methods("GET")

	r.HandleFunc("/v2/configuration/{repo_name}/statefile", service.TerraformerStateHandler(session)).Methods("POST")

	log.Println("Server will listen at port", config.Server.HTTPAddr, config.Server.HTTPPort)
	muxWithMiddlewares := http.TimeoutHandler(r, time.Second*60, "Timeout!")
	err = endless.ListenAndServe(fmt.Sprintf("%s:%d", config.Server.HTTPAddr, config.Server.HTTPPort), muxWithMiddlewares)
	if err != nil {
		log.Println("Couldn't start the server", err)
	}
}

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()
	c := session.DB("action").C("actionDetails")

	index := mgo.Index{
		Key:        []string{"actionid"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}
