// @SubApi ibmcloud-terraform-provider [/v1]
// @SubApi Allows you access ibm cloud terraform provider api [/v1]

package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// var httpClient *http.Client
// var sessionMgo *mgo.Session
// var githubToken string
// var githubIBMToken string
// var currentOps = make(map[string]chan StatusResponse)
var planTimeOut = 60 * time.Minute

// ConfigRequest -
type ConfigRequest struct {
	GitURL        string            `json:"git_url,required" description:"The git url of your configuraltion"`
	ConfigName    string            `json:"config_name,omitempty" description:"The configuration repo name"`
	VariableStore *VariablesRequest `json:"variablestore,omitempty" description:"The environments' variable store"`
	LOGLEVEL      string            `json:"log_level,omitempty" description:"The log level defing by user."`
	// Terraformer   string            `json:"terraformer,omitempty" description:"The terraformer."`
	Service string `json:"service,omitempty" description:"The terraformer services."`
}

// ConfigResponse -
type ConfigResponse struct {
	ConfigName string `json:"config_name,required" description:"configuration name"`
}

// StatusResponse -
type StatusResponse struct {
	Status string `json:"status,required" description:"Status of the terraform operation."`
	Error  string `json:"error,omitempty" description:"Error of the terraform operation."`
}

// ActionResponse -
type ActionResponse struct {
	ConfigName string `json:"id,required" description:"Name of the configuration"`
	Action     string `json:"action,required" description:"Action Name"`
	ActionID   string `json:"action_id"`
	Timestamp  string `json:"timestamp"`
	Status     string `json:"status"`
}

// ActionDetails -
type ActionDetails struct {
	ConfigName string `json:"id,required" description:"Name of the configuration"`
	Action     string `json:"action,required" description:"Action Name"`
	ActionID   string `json:"action_id"`
	Output     string `json:"output"`
	Error      string `json:"error"`
}

// VariablesRequest -
type VariablesRequest []EnvironmentVariableRequest

// EnvironmentVariableRequest -
type EnvironmentVariableRequest struct {
	Name  string `json:"name,required" binding:"required" description:"The variable's name"`
	Value string `json:"value,required" binding:"required" description:"The variable's value"`
}

func init() { // todo: @srikar - is this needed

	if currentDir == "" {
		panic("MOUNT_DIR is not set. Please set MOUNT_DIR to continue")
	}

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, os.ModePerm)
	}
	if _, err := os.Stat(stateDir); os.IsNotExist(err) {
		os.MkdirAll(stateDir, os.ModePerm)
	}

}

//ConfHandler handles request to kickoff git clone of the repo.
// @Title ConfHandler
// @Description clone the configuration repo
// @Accept  json
// @Produce  json
// @Param   body     body     ConfigRequest   true "request body"
// @Success 200 {object} ConfigResponse
// @Failure 500 {object} string
// @Failure 400 {object} string
// @Router /v1/configuration [post]
func ConfHandler(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Read body
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Unmarshal
		var msg ConfigRequest
		var response ConfigResponse
		err = json.Unmarshal(b, &msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println(msg.GitURL)
		if msg.GitURL == "" {
			w.WriteHeader(400)
			w.Write([]byte("EMPTY GIT URL"))
			return
		}

		if msg.LOGLEVEL != "" {
			os.Setenv("TF_LOG", msg.LOGLEVEL)
		}

		log.Println("Will clone git repo")

		_, configName, err := CloneRepo(msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("\n", configName)

		response.ConfigName = configName
		log.Println(response)

		output, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return
		}

		confDir := path.Join(currentDir, configName)

		b = make([]byte, 10)
		rand.Read(b)
		randomID := fmt.Sprintf("%x", b)

		err = TerraformInit(confDir, planTimeOut, randomID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(output)
	}
}

//ConfDeleteHandler handles request to kickoff delete for the configuration repo.
// @Title ConfDeleteHandler
// @Description delete the configuration repo
// @Param   repo_name     path    string     true "Some ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Failure 404 {object} string
// @Router /v1/configuration/{repo_name} [delete]
func ConfDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
	}

	vars := mux.Vars(r)
	repoName := vars["repo_name"]

	err := removeRepo(currentDir, repoName)
	if err != nil {
		w.WriteHeader(404)
		log.Println(err)
		w.Write([]byte(fmt.Sprintf("There is no config repo file for this request.")))
		return
	}
	log.Println("Deleted repo....")
}

//PlanHandler handles request to run terraform plan.
// @Title PlanHandler
// @Description Execute plan for the configuration.
// @Param   SLACK_WEBHOOK_URL     header    string     false "provide slack webhook url"
// @Param   repo_name     path    string     true "Repo Name"
// @Accept  json
// @Produce  json
// @Success 202 {object} ActionResponse
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /v1/configuration/{repo_name}/plan [post]
func PlanHandler(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		webhook := r.Header.Get("SLACK_WEBHOOK_URL")
		vars := mux.Vars(r)
		repoName := vars["repo_name"]

		var actionResponse ActionResponse
		var statusResponse StatusResponse

		log.Println("Url Param 'repo name' is: " + repoName)
		confDir := path.Join(currentDir, repoName)

		b := make([]byte, 10)
		rand.Read(b)
		randomID := fmt.Sprintf("%x", b)

		outURL := "http://" + r.Host + "/" + r.URL.Path + "/" + randomID + ".out"
		errURL := "http://" + r.Host + "/" + r.URL.Path + "/" + randomID + ".err"

		// Post to slack that the action has started and the link logs
		ResultToSlack(outURL, errURL, "plan", randomID, "In-Progress", webhook)

		go func() {
			PullRepo(repoName)
			err := TerraformPlan(confDir+pathSep+repoName, planTimeOut, randomID)
			if err != nil {
				statusResponse.Error = err.Error()
				statusResponse.Status = "Failed"

				// Update the status in the db in case it is failed
				err = UpdateMongodb(s, randomID, statusResponse.Status)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}
			statusResponse.Status = "Completed"

			// Update the status in the db in case it is completed
			err = UpdateMongodb(s, randomID, statusResponse.Status)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ResultToSlack(outURL, errURL, "plan", randomID, "Completed", webhook)
		}()

		w.WriteHeader(202)

		actionResponse.Action = "plan"
		actionResponse.ConfigName = repoName
		actionResponse.ActionID = randomID
		actionResponse.Timestamp = time.Now().Format("20060102150405")
		actionResponse.Status = "In-Progress"

		// Make an entry in the db
		InsertMongodb(s, actionResponse)

		output, err := json.MarshalIndent(actionResponse, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(output)

	}
}

//ApplyHandler handles request to run terraform apply.
// @Title ApplyHandler
// @Description Execute apply for the configuration.
// @Param   SLACK_WEBHOOK_URL     header    string     false "provide slack webhook url"
// @Param   repo_name     path    string     true "Repo Name"
// @Accept  json
// @Produce  json
// @Success 202 {object} ActionResponse
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /v1/configuration/{repo_name}/apply [post]
func ApplyHandler(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		webhook := r.Header.Get("SLACK_WEBHOOK_URL")

		var actionResponse ActionResponse
		var statusResponse StatusResponse

		vars := mux.Vars(r)
		repoName := vars["repo_name"]

		log.Println("Url Param 'repo name' is: " + repoName)
		confDir := path.Join(currentDir, repoName)

		b := make([]byte, 10)
		rand.Read(b)
		randomID := fmt.Sprintf("%x", b)

		outURL := "http://" + r.Host + "/" + r.URL.Path + "/" + randomID + ".out"
		errURL := "http://" + r.Host + "/" + r.URL.Path + "/" + randomID + ".err"
		ResultToSlack(outURL, errURL, "apply", randomID, "In-Progress", webhook)
		go func() {

			PullRepo(repoName)
			err := TerraformApply(confDir, stateDir, repoName, planTimeOut, randomID)
			if err != nil {
				statusResponse.Error = err.Error()
				statusResponse.Status = "Failed"

				err = UpdateMongodb(s, randomID, statusResponse.Status)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}
			statusResponse.Status = "Completed"
			ResultToSlack(outURL, errURL, "apply", randomID, statusResponse.Status, webhook)
			err = UpdateMongodb(s, randomID, statusResponse.Status)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}()
		w.WriteHeader(202)
		actionResponse.Action = "apply"
		actionResponse.ConfigName = repoName
		actionResponse.ActionID = randomID
		actionResponse.Timestamp = time.Now().Format("20060102150405")
		actionResponse.Status = "In-Progress"

		InsertMongodb(s, actionResponse)
		output, err := json.MarshalIndent(actionResponse, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(output)

	}
}

//DestroyHandler handles request to run terraform delete.
// @Title DestroyHandler
// @Description Execute destroy for the configuration.
// @Param   SLACK_WEBHOOK_URL     header    string     false "provide slack webhook url"
// @Param   repo_name     path    string     true "Repo Name"
// @Accept  json
// @Produce  json
// @Success 202 {object} ActionResponse
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /v1/configuration/{repo_name}/destroy [post]
func DestroyHandler(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		var actionResponse ActionResponse
		var statusResponse StatusResponse

		webhook := r.Header.Get("SLACK_WEBHOOK_URL")
		vars := mux.Vars(r)
		repoName := vars["repo_name"]

		log.Println("Url Param 'repo name' is: " + repoName)
		confDir := path.Join(currentDir, repoName)

		b := make([]byte, 10)
		rand.Read(b)
		randomID := fmt.Sprintf("%x", b)

		outURL := "http://" + r.Host + "/" + r.URL.Path + "/" + randomID + ".out"
		errURL := "http://" + r.Host + "/" + r.URL.Path + "/" + randomID + ".err"
		ResultToSlack(outURL, errURL, "destroy", randomID, "In-Progress", webhook)
		go func() {
			err := TerraformDestroy(confDir, stateDir, repoName, planTimeOut, randomID)
			if err != nil {
				statusResponse.Error = err.Error()
				statusResponse.Status = "Failed"

				err = UpdateMongodb(s, randomID, statusResponse.Status)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}
			statusResponse.Status = "Completed"
			ResultToSlack(outURL, errURL, "destroy", randomID, "Completed", webhook)
			err = UpdateMongodb(s, randomID, statusResponse.Status)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}()

		w.WriteHeader(202)
		actionResponse.Action = "destroy"
		actionResponse.ConfigName = repoName
		actionResponse.ActionID = randomID
		actionResponse.Timestamp = time.Now().Format("20060102150405")
		actionResponse.Status = "In-Progress"

		InsertMongodb(s, actionResponse)

		output, err := json.MarshalIndent(actionResponse, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(output)

	}
}

//ShowHandler handles request to run terraform show.
// @Title ShowHandler
// @Description Execute show for the configuration.
// @Param   SLACK_WEBHOOK_URL     header    string     false "provide slack webhook url"
// @Param   repo_name     path    string     true "Repo Name"
// @Accept  json
// @Produce  json
// @Success 202 {object} ActionResponse
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /v1/configuration/{repo_name}/show [post]
func ShowHandler(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		session := s.Copy()
		defer session.Close()

		var actionResponse ActionResponse
		var statusResponse StatusResponse

		webhook := r.Header.Get("SLACK_WEBHOOK_URL")
		vars := mux.Vars(r)

		repoName := vars["repo_name"]

		log.Println("Url Param 'repo name' is: " + repoName)
		confDir := path.Join(currentDir, repoName)

		b := make([]byte, 10)
		rand.Read(b)
		randomID := fmt.Sprintf("%x", b)

		outURL := "http://" + r.Host + "/" + r.URL.Path + "/" + randomID + ".out"
		errURL := "http://" + r.Host + "/" + r.URL.Path + "/" + randomID + ".err"
		ResultToSlack(outURL, errURL, "show", randomID, "In-Progress", webhook)
		go func() {
			err := TerraformShow(confDir, stateDir, repoName, planTimeOut, randomID)
			if err != nil {
				statusResponse.Error = err.Error()
				statusResponse.Status = "Failed"

				err = UpdateMongodb(s, randomID, statusResponse.Status)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}
			outURL := "http://" + r.Host + "/" + r.URL.Path + "/" + randomID + ".out"
			errURL := "http://" + r.Host + "/" + r.URL.Path + "/" + randomID + ".err"

			statusResponse.Status = "Completed"
			ResultToSlack(outURL, errURL, "show", randomID, statusResponse.Status, webhook)
			err = UpdateMongodb(s, randomID, statusResponse.Status)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}()
		w.WriteHeader(202)

		actionResponse.Action = "show"
		actionResponse.ConfigName = repoName
		actionResponse.ActionID = randomID
		actionResponse.Timestamp = time.Now().Format("20060102150405")
		actionResponse.Status = "In-Progress"

		InsertMongodb(s, actionResponse)

		output, err := json.MarshalIndent(actionResponse, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(output)

	}
}

//LogHandler handles request to get the log.
// @Title LogHandler
// @Description Get logs for the configuration.
// @Param   repo_name     path    string     true "repo name"
// @Param   action_name     path    string     true "action name"
// @Param   action_id     path    string     true "action id"
// @Accept  json
// @Produce  json
// @Success 200 {object} ActionDetails
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /v1/configuration/{repo_name}/{action_name}/{action_id}/log [get]
func LogHandler(w http.ResponseWriter, r *http.Request) {

	var response ActionDetails

	vars := mux.Vars(r)
	repoName := vars["repo_name"]
	action := vars["action"]
	actionID := vars["actionID"]

	log.Println("Url Param 'repo name' is: " + repoName)
	log.Println("Url Param 'action' is: " + action)
	log.Println("Url Param 'actionID' is: " + actionID)

	outFile, errFile, err := readLogFile(actionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response.ConfigName = repoName
	response.Output = outFile
	response.Error = errFile
	response.Action = action
	response.ActionID = actionID

	output, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)

}

//StatusHandler handles request to get the action status.
// @Title StatusHandler
// @Description Get status of the action.
// @Param   repo_name     path    string     true "repo name"
// @Param   action_name     path    string     true "action name"
// @Param   action_id     path    string     true "action id"
// @Accept  json
// @Produce  json
// @Success 200 {object} StatusResponse
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /v1/configuration/{repo_name}/{action_name}/{action_id}/status [get]
func StatusHandler(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Status call...")
		session := s.Copy()
		defer session.Close()

		var response StatusResponse
		var actionResponse ActionResponse

		vars := mux.Vars(r)
		repoName := vars["repo_name"]
		action := vars["action"]
		actionID := vars["actionID"]

		log.Println("Url Param 'repo name' is: " + repoName)
		log.Println("Url Param 'action' is: " + action)
		log.Println("Url Param 'actionID' is: " + actionID)

		c := session.DB("action").C("actionDetails")
		err := c.Find(bson.M{"actionid": actionID}).One(&actionResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response.Status = actionResponse.Status
		output, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(output)
	}
}

//ViewLogHandler handles request to retrieve the log file
func ViewLogHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	logFile := vars["log_file"]

	body, err := ioutil.ReadFile(path.Join(logDir, logFile))
	if err != nil {
		w.WriteHeader(404)
		log.Println(err)
		w.Write([]byte(fmt.Sprintf("There is no log file for this request")))
		return
	}
	w.WriteHeader(200)
	w.Write(body)
}

//GetActionDetailsHandler handles request to get all the information for a particular action.
// @Title GetActionDetailsHandler
// @Description Get all the information for a particular action
// @Param   repo_name     path    string     true "repo name"
// @Param   action_name     path    string     true "action name"
// @Accept  json
// @Produce  json
// @Success 200 {object} ActionResponse
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /v1/configuration/{repo_name}/{action_name} [get]
func GetActionDetailsHandler(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("action details handler..")
		session := s.Copy()
		defer session.Close()

		vars := mux.Vars(r)
		action := vars["action"]

		var actionResponse []ActionResponse
		c := session.DB("action").C("actionDetails")

		err := c.Find(bson.M{"action": action}).All(&actionResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(200)

		output, err := json.MarshalIndent(actionResponse, "", "  ")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(output)

	}
}
