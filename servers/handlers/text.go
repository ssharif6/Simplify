package handlers

import (
	"net/http"
	"encoding/json"
	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
	"context"
)

type RequestObject struct {
	Input string `json:"input,omitempty"`
	Url string `json:"url"`
}

type ResponseObject struct {
	Entities []*languagepb.Entity `json:"entities,omitempty"`
	T1Objects []*Eli5T1Model `json:"t1objects"`
}

func TextHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Add headers
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")

		// request checking
		ro := &RequestObject{}
		err := json.NewDecoder(r.Body).Decode(ro)
		if len(ro.Input) == 0 {
			http.Error(w, "input cannot be empty", http.StatusBadRequest)
		}
		if err != nil {
			http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
		}

		// Setup context
		ctx := context.Background()
		client, err := language.NewClient(ctx)
		if err != nil {
			http.Error(w, "Unable to connect to google cog services", http.StatusInternalServerError)
		}


		entities, err := handleInput(ctx, client, ro)
		if err != nil {
			http.Error(w, "Unable to extract entities", http.StatusInternalServerError)
		}

		eli5, err := QueryEli5(entities)
		if err != nil {
			http.Error(w, "unable to get eli5 responses", http.StatusInternalServerError)
		}

		t1, err := Query(eli5)
		if err != nil {
			http.Error(w, "unable to get eli5 responses", http.StatusInternalServerError)
		}


		r := ResponseObject{
			Entities: entities,
			T1Objects: t1,
		}

		encodeErr := json.NewEncoder(w).Encode(r)
		if encodeErr != nil {
			http.Error(w, "error encoding json", http.StatusInternalServerError)
		}
	}

}

func handleInput(ctx context.Context, client *language.Client, requestObject *RequestObject) ([]*languagepb.Entity, error) {
	input := requestObject.Input

	//url := requestObject.url
	resp, err := AnalyzeEntities(ctx, client, input)
	if err != nil {
		return nil, err
	}
	entities := resp.Entities
	return entities, nil
}

