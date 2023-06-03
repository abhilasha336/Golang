/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/spf13/cobra"
)

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "get a random joke",
	Long:  `this gives a random joke from a api hit`,
	Run: func(cmd *cobra.Command, args []string) {

		termStr, err := cmd.Flags().GetString("term")
		if err != nil {
			log.Printf("get string error from command:%v", err)
		}
		if termStr != "" {
			jokes(termStr)
		} else {
			getRandomJoke()

		}
	},
}

func init() {
	rootCmd.AddCommand(randomCmd)

	randomCmd.PersistentFlags().String("term", "", "term flag can add joke on command line with this flag --flag=joke")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// randomCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// randomCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Joke struct {
	Id     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"statuscode"`
}

type TermJoke struct {
	Results    json.RawMessage `json:"results"`
	SearchTerm string          `json:"search_term"`
	Status     int             `json:"status"`
	TotalJokes int             `json:"total_jokes"`
}

func getRandomJoke() {
	jokeHitApi := "https://icanhazdadjoke.com/"
	dataByte := getJokeData(jokeHitApi)

	jokeInstanceStruct := Joke{}

	err := json.Unmarshal(dataByte, &jokeInstanceStruct)
	if err != nil {
		log.Printf("unmarshal struct error:%s", err)
	}

	fmt.Printf("data is: %+v", jokeInstanceStruct)

}

func termJoke(termstr string) (totaljokes int, jokeSlice []Joke) {

	url := fmt.Sprintf("https://icanhazdadjoke.com/search?term=%s", termstr)

	responseByte := getJokeData(url)
	termJoke := TermJoke{}
	err := json.Unmarshal(responseByte, &termJoke)
	if err != nil {
		fmt.Printf("unmarshall termjoke :%v", err)
	}

	jokesSlice := []Joke{}
	err = json.Unmarshal(termJoke.Results, &jokesSlice)
	if err != nil {
		fmt.Printf("unmarshall results joke arrays :%v", err)
	}

	return termJoke.TotalJokes, jokesSlice
}
func randomiseJokeList(length int, jokeList []Joke) {

	min := 0
	max := length - 1
	if length <= 0 {
		err := fmt.Errorf("no jokes found with this term")
		fmt.Println(err.Error())

	}
	randomNum := min + rand.Intn(max-min)
	fmt.Println(jokeList[randomNum].Joke)

}

func jokes(termstr string) {
	total, result := termJoke(termstr)
	randomiseJokeList(total, result)
}

// function accepts api and return bytes of response data
func getJokeData(baseApi string) []byte {

	request, err := http.NewRequest(http.MethodGet, baseApi, nil)
	if err != nil {
		log.Printf("request error %s", err)
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-Agent", "Dadjoke CLI (https://github.com/avlashabhi336/cobracliapi)")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("response error by hitting request api with header: %s", err)
	}

	dataByte, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("read error from response.Body: %s\n", err)
	}
	return dataByte
}
