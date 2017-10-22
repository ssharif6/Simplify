package handlers

import (
	"net/http"
	"encoding/json"
	"fmt"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
	"bytes"
	"log"
	"net/url"
)

const eli5BaseUrl = "https://www.reddit.com/r/explainlikeimfive/search.json?q="
const eli5SuffixUrl = "restrict_sr=on&sort=relevance"

type Eli5Response struct {
	Url string
}

type Eli5T3Model struct {
	Id    string `json:"id,omitempty"`
	Score int    `json:"score,omitempty"`
	Url   string `json:"url,omitempty"`
	Ups   int    `json:"ups,omitempty"`
	Downs int    `json:"downs,omitempty"`
}

type Eli5T1Model struct {
	Id    string `json:"id,omitempty"`
	Score int    `json:"score,omitempty"`
	Ups   int    `json:"ups,omitempty"`
	Downs int    `json:"downs,omitempty"`
	Url   string `json:"url,omitempty"`
	Body  string `json:"body,omitempty"`
}

type T1SearchModel struct {
	Kind string `json:"kind,omitempty"`
	Data struct {
		Children []struct {
			Kind string `json:"kind,omitempty"`
			Data struct {
				Score int    `json:"score,omitempty"`
				Body  string `json:"body,omitempty"`
				Ups   int    `json:"ups,omitempty"`
				Downs int    `json:"downs,omitempty"`
				Id string `json:"id,omitempty"`
			} `json:"data,omitempty"`
		} `json:"children,omitempty"`
	} `json:"data,omitempty"`
}

type T3SearchModel struct {
	Kind string `json:"kind,omitempty"`
	Data struct {
		Children []struct {
			Kind string `json:"kind,omitempty"`
			Data struct {
				Score int    `json:"score,omitempty"`
				URL   string `json:"url,omitempty"`
				Ups   int    `json:"ups,omitempty"`
				Downs int    `json:"downs,omitempty"`
				Id string `json:"id,omitempty"`
			} `json:"data,omitempty"`
		} `json:"children,omitempty"`
	} `json:"data,omitempty"`
}

// Get t3
func QueryEli5(entities []*languagepb.Entity) ([]*Eli5T3Model, error) {
	responseObjs := []*T3SearchModel{}
	client := &http.Client{}
	for _, c := range entities {
		fmt.Println(eli5BaseUrl+url.QueryEscape(c.Name)+"&"+eli5SuffixUrl)
		req, err := http.NewRequest("GET", eli5BaseUrl+url.QueryEscape(c.Name)+"&"+eli5SuffixUrl, nil)
		if err != nil {
			log.Fatalln(err)
		}

		req.Header.Set("User-Agent", "your bot 0.1")
		resp, err := client.Do(req)

		if err != nil {
			return nil, err
		}

		eli5Res := &T3SearchModel{}
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)

		//s := buf.String() // Does a complete copy of the bytes in the buffer.
		decodeErr := json.NewDecoder(buf).Decode(eli5Res)
		if decodeErr != nil {
			fmt.Printf("There was an error decoding the json. err = %s", decodeErr)
			fmt.Println("error2")
			return nil, err
		}
		responseObjs = append(responseObjs, eli5Res)
	}

	return getEli5Response(responseObjs), nil
}

func Query(models []*Eli5T3Model) ([]*Eli5T1Model, error) {
	// Iterate through each model, perform http request, extract relevant data
	client := &http.Client{}
	responseObjs := [][]T1SearchModel{}
	for _, c := range models {
		url := c.Url

		req, err := http.NewRequest("GET", url + ".json?sort=confidence", nil)
		req.Header.Set("User-Agent", "your bot 0.1")
		if err != nil {
			return nil, err
		}
		resp, err := client.Do(req)

		if err != nil {
			return nil, err
		}


		buf := new(bytes.Buffer)

		buf.ReadFrom(resp.Body)

		//s := buf.String() // Does a complete copy of the bytes in the buffer.
		keys := make([]T1SearchModel, 0)
		//decodeErr := json.NewDecoder(buf).Decode(eli5Res)
		decodeErr := json.Unmarshal(buf.Bytes(), &keys)
		if decodeErr != nil {
			fmt.Printf("There was an error decoding the json. err = %s", decodeErr)
			fmt.Errorf("error decoding json")
		}
		responseObjs = append(responseObjs, keys)
	}

	t1Objs := extractT1Objects(responseObjs)
	return t1Objs, nil
}

func extractT1Objects(models [][]T1SearchModel) ([]*Eli5T1Model) {
	responseObjs := []*Eli5T1Model{}
	for _, c := range models {
		for _, x := range c {
			children := x.Data.Children
			for _, d := range children {
				if d.Kind == "t1" {
					data := d.Data
					t1 := Eli5T1Model{
						Ups:   data.Ups,
						Score: data.Score,
						Id:    data.Id,
						Downs: data.Downs,
						Body:  data.Body,
					}
					responseObjs = append(responseObjs, &t1)
				}
			}

		}
	}
	return responseObjs
}

// TODO Error handling
func getEli5Response(eli5SearchObjects []*T3SearchModel) ([]*Eli5T3Model) {
	responseObjs := []*Eli5T3Model{}
	for _, c := range eli5SearchObjects {
		children := c.Data.Children
		if len(children) > 0 {
			for i := 0; i < 5; i++ {
				obj := children[i].Data
				t3 := Eli5T3Model{
					Downs: obj.Downs,
					Id:    obj.Id,
					Score: obj.Score,
					Ups:   obj.Ups,
					Url:   obj.URL,
				}
				responseObjs = append(responseObjs, &t3)
			}

		} else {
			fmt.Println("LENGTH IS 0")
		}
	}
	return responseObjs
}
