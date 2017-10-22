package handlers

import (
	"net/http"
	"encoding/json"
	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
	"context"
	"fmt"
	"os"
	"github.com/golang/protobuf/proto"
	"log"
)

type RequestObject struct {
	Input string `json:"input,omitempty"`
	Url string `json:"url"`
}

type ResponseObject struct {
	Entities []*languagepb.Entity `json:"entities,omitempty"`
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
		if client == nil {
			fmt.Printf("client is NULL@@@@@@@@@@@@@@@@@@@@@")
		}

		entities, err := analyzeEntities(ctx, client, ro.Input)
		if err != nil {
			http.Error(w, "Unable to extract entities", http.StatusInternalServerError)
		}
		encodeErr := json.NewEncoder(w).Encode(entities)
		if encodeErr != nil {
			http.Error(w, "error encoding json", http.StatusInternalServerError)
		}

	}

}

func handleInput(ctx context.Context, client *language.Client, requestObject *RequestObject) ([]*languagepb.Entity, error) {
	input := requestObject.Input
	//url := requestObject.url
	resp, err := analyzeEntities(ctx, client, input)
	if err != nil {
		return nil, err
	}
	entities := resp.Entities
	return entities, nil
}

func usage(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Fprintln(os.Stderr, "usage: analyze [entities|sentiment|syntax|entitysentiment] <text>")
	os.Exit(2)
}

func analyzeEntities(ctx context.Context, client *language.Client, text string) (*languagepb.AnalyzeEntitiesResponse, error) {
	if client == nil || ctx == nil{
		fmt.Println("it's null @@@@@@@@@@@@@@@@@@@@@@")
	}
	return client.AnalyzeEntities(ctx, &languagepb.AnalyzeEntitiesRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
}

func analyzeSentiment(ctx context.Context, client *language.Client, text string) (*languagepb.AnalyzeSentimentResponse, error) {
	return client.AnalyzeSentiment(ctx, &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
	})
}

func analyzeSyntax(ctx context.Context, client *language.Client, text string) (*languagepb.AnnotateTextResponse, error) {
	return client.AnnotateText(ctx, &languagepb.AnnotateTextRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		Features: &languagepb.AnnotateTextRequest_Features{
			ExtractSyntax: true,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
}

func printResp(v proto.Message, err error) {
	if err != nil {
		log.Fatal(err)
	}
	proto.MarshalText(os.Stdout, v)
}
