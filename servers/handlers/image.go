package handlers

import (
	"net/http"
	"encoding/json"
	"os"
	"math/rand"
	"io"
	"fmt"
)

type ImageResponseObject struct {
	Texts []string `json:"texts,omitempty"`
	Labels []string `json:"labels,omitempty"`
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	// do something
	fmt.Println("2")

	if r.Method == "POST" {
		// Add headers
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")

		// request checking
		ro := &RequestObject{}
		err := json.NewDecoder(r.Body).Decode(ro)
		if len(ro.Url) == 0 {
			http.Error(w, "input cannot be empty", http.StatusBadRequest)
		}
		if err != nil {
			http.Error(w, "Unable to parse JSON", http.StatusBadRequest)

		}
		url := ro.Url
		response, e := http.Get(url)
		if e != nil {
			http.Error(w, "bad image", http.StatusBadRequest)
		}
		defer response.Body.Close()
		file, err := os.Create("/tmp/asdf.jpg")
		if err != nil {
			http.Error(w, "error creating file", http.StatusInternalServerError)
			return
		}
		// Use io.Copy to just dump the response body to the file. This supports huge files
		_, err = io.Copy(file, response.Body)
		if err != nil {
			http.Error(w, "error copying", http.StatusInternalServerError)
			return
		}

		r := &ImageResponseObject{}
		detectTextResp, err := DetectText("/tmp/asdf.jpg")
		if err != nil {
			http.Error(w, "error with DETECT TEXTY", http.StatusInternalServerError)
		}

		if len(detectTextResp) > 5 {
			r.Texts = detectTextResp
		} else {
			detectLabelsResp, err := DetectLabels("/tmp/asdf.jpg")
			if err != nil {
				http.Error(w, "error with DETECT TEXTY", http.StatusInternalServerError)
			}
			r.Labels = detectLabelsResp
		}

		jsonErr := json.NewEncoder(w).Encode(r)
		if jsonErr != nil {
			http.Error(w, "error with my life", http.StatusInternalServerError)

		}

		//// Setup context
		//ctx := context.Background()
		//client, err := vision.NewImageAnnotatorClient(ctx)
		//if err != nil {
		//	http.Error(w, "Unable to connect to google cog services", http.StatusInternalServerError)
		//}
		//
		//DetectText()
		//image, err := vision.NewImageFromReader(file)
		//if err != nil {
		//	http.Error(w, "error creating new imageReader from file", http.StatusInternalServerError)
		//	return
		//}
		//
		//texts, err := client.DetectTexts(ctx, image, nil, 10)
		//if err != nil {
		//	http.Error(w, "error with client", http.StatusInternalServerError)
		//	return
		//}
		//stringSlice := make([]string, 0)
		//
		//for _, c := range texts {
		//	stringSlice = append(stringSlice, c.Description)
		//}
		//
		//r := &ImageResponseObject{
		//	Texts: stringSlice,
		//}
		//jsonErr := json.NewEncoder(w).Encode(r)
		//if jsonErr != nil {
		//	http.Error(w, "error with json", http.StatusInternalServerError)
		//	return
		//}
		//
		//defer file.Close()
	}
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
