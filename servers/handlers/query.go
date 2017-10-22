package handlers

import (
	"net/http"
	"encoding/json"
	"fmt"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
	"bytes"
	"log"
)

const eli5BaseUrl = "https://www.reddit.com/r/explainlikeimfive/search.json?q="
const eli5SuffixUrl = "restrict_sr=on&sort=relevance&t=all"

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

		req, err := http.NewRequest("GET", eli5BaseUrl+c.Name+"&"+eli5SuffixUrl, nil)
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

	fmt.Println("QueryEli5")
	fmt.Println(responseObjs)

	return getEli5Response(responseObjs), nil
}

func Query(models []*Eli5T3Model) ([]*Eli5T1Model, error) {
	// Iterate through each model, perform http request, extract relevant data
	fmt.Println("GOT TO QUERY")
	client := &http.Client{}
	responseObjs := [][]T1SearchModel{}
	for _, c := range models {
		url := c.Url
		fmt.Println(url)

		req, err := http.NewRequest("GET", url + ".json?sort=top", nil)
		//req, err := http.NewRequest("GET","https://www.reddit.com/r/explainlikeimfive/comments/5u2nkx/eli5_what_is_a_linux_kernel/.json?sort=top", nil)
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
		eli5Res := &T1SearchModel{}
		//decodeErr := json.NewDecoder(buf).Decode(eli5Res)
		decodeErr := json.Unmarshal(buf.Bytes(), &keys)
		if decodeErr != nil {
			fmt.Printf("There was an error decoding the json. err = %s", decodeErr)
			fmt.Errorf("error decoding json")
		}
		fmt.Println("kanye")
		fmt.Println(*eli5Res)
		responseObjs = append(responseObjs, keys)
	}


	fmt.Println("LENGTH OF RESPONSE OBJCTS")
	fmt.Println(len(responseObjs))
	t1Objs := extractT1Objects(responseObjs)
	return t1Objs, nil
}

func extractT1Objects(models [][]T1SearchModel) ([]*Eli5T1Model) {
	fmt.Println("GOT TO EXTRACT T1 OBJECTS")
	responseObjs := []*Eli5T1Model{}
	for _, c := range models {
		fmt.Println("C")
		for _, x := range c {
			children := x.Data.Children
			fmt.Println("LENGTH OF CHILDREN")
			fmt.Println(len(children))
			for _, d := range children {
				fmt.Println(d.Kind)
				if d.Kind == "t1" {
					data := d.Data
					fmt.Println("DATA@#@@@@@@")
					fmt.Println(data)
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
	fmt.Println("GOT TO END OF EXTRACT T1 OBJECTS")
	fmt.Println("LENGTH OF RESPONSEOBJS")
	fmt.Println(len(responseObjs))

	return responseObjs
}

// TODO Error handling
func getEli5Response(eli5SearchObjects []*T3SearchModel) ([]*Eli5T3Model) {
	fmt.Println("Get eli5 response")
	responseObjs := []*Eli5T3Model{}
	for _, c := range eli5SearchObjects {
		fmt.Println("CO OBJECT")
		fmt.Println(c.Data.Children)
		children := c.Data.Children
		if len(children) > 0 {

			fmt.Println("THIS IS THE LENGHT OF ELI5RESPONSE")
			fmt.Println(len(children))
			for i := 0; i < len(c.Data.Children); i++ {
				obj := children[i].Data
				fmt.Println("CHILDREN OBJECT")
				//fmt.Println(obj)
				fmt.Println(obj.Id)
				fmt.Println(obj.Score)
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
	fmt.Println("LENGTH OF ELI5RESPONSE")
	fmt.Println(len(responseObjs))
	fmt.Println("GOT TO END OF ELI5RESPONSE")
	return responseObjs
}
