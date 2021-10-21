package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

const queueApi = "https://digitalmakerspace.uncw.edu/api/v1/queue"
const imageApi = "https://digitalmakerspace.uncw.edu/api/v1/image_rotation"

//go:embed *.gohtml
var tpls embed.FS

var templates = template.Must(template.ParseFS(tpls, "*"))

type image struct {
	Filename string `json:"printed_picture_file"`
	Title    string `json:"printed_picture_file_original"`
}

type rotationResponse struct {
	Status string  `json:"status"`
	Items  []image `json:"data"`
}

func rotatonHandler(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get(imageApi)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var data rotationResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}
	if data.Status != "good!" {
		log.Fatal(fmt.Sprintf("API gave bad response: %v", data.Status))
	}
	templates.ExecuteTemplate(w, "rotation.gohtml", data)
}

type detail struct {
	Title     string `json:"title"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Status    string `json:"status"`
	Filename  string `json:"printed_picture_file"`
}

type queueResponse struct {
	Status string   `json:"status"`
	Items  []detail `json:"data"`
}

func queueHandler(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get(queueApi)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var data queueResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}
	if data.Status != "good!" {
		log.Fatal(fmt.Sprintf("API gave bad response: %v", data.Status))
	}
	statusMap := map[string]string{
		"Just Arrived":          "Submitted",
		"Reviewed and Approved": "Reviewed and Approved",
		"Started Printing":      "Printing",
		"Ready for Pickup":      "Ready for Pickup",
	}
	for index, i := range data.Items {
		data.Items[index].Status = statusMap[i.Status]
	}
	templates.ExecuteTemplate(w, "queue.gohtml", data)
}

func main() {
	fmt.Println("A Golang port of 3dprintqueueview by garrett")
	fmt.Println("Serving pages at http://localhost:8000/ and http://localhost:8000/image_rotation")
	fmt.Println("nodjs version code at https://libapps-admin.uncw.edu/randall-dev/3d-prin-queue-view")
	fmt.Println("go version code at https://github.com/GarrettArm/go-3d-print-queue-view")

	http.HandleFunc("/", queueHandler)
	http.HandleFunc("/image_rotation", rotatonHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
