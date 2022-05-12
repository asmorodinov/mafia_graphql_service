package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

var addr = flag.String("addr", "localhost:8090", "GraphQL server address")

func request(bodyStr string) (string, string, int, error) {
	body := bytes.NewBufferString(bodyStr)
	req, err := http.NewRequest(http.MethodPost, "http://"+*addr+"/graphql", body)
	if err != nil {
		return "", "", 0, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", 0, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", "", 0, err
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, respBody, "", "    ")
	if err != nil {
		return string(respBody), "", 0, err
	}

	return pretty.String(), resp.Status, resp.StatusCode, err
}

func jsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	// Trim the beginning and trailing " character
	return string(b[1 : len(b)-1])
}

func requestQuery(query string, variables map[string]interface{}) {
	b, err := json.Marshal(map[string]interface{}{
		"query":     jsonEscape(query),
		"variables": variables,
	})
	if err != nil {
		panic(err)
	}
	resp, status, code, err := request(string(b))
	fmt.Printf("%v\n%v\n%v\n%v\n", resp, status, code, err)
}

var gameId = flag.Int("getGameID", -1, "id of the game that you want to see info about")
var listGames = flag.Bool("getGamesList", false, "set flag to get list of all games")
var commentBody = flag.String("commentBody", "", "body of the comment to add to the game")
var idOfTheGameToCommentAbout = flag.Int("commentGameId", -1, "id of the game that you want to comment about")

func main() {
	flag.Parse()

	// get game info
	if *gameId != -1 {
		fmt.Printf("Getting info about game with id %v\n", *gameId)

		query := `
			query GetGame($id: Int!){
				game(id: $id) {
					id
					roomName
					status
					startedAt
					finishedAt
					isDay
					day
					players {
						login
						isAlive
						deadCause
						role
					}
					comments {
						body
					}
				}
			}
		`
		variables := map[string]interface{}{"id": *gameId}
		requestQuery(query, variables)
	}
	// get games list
	if *listGames {
		fmt.Println("Getting list of all games")
		query := `
			{
				list {
					id
					roomName
					status
					startedAt
					finishedAt
					isDay
					day
					players {
						login
						isAlive
						deadCause
						role
					}
					comments {
						body
					}
				}
			}
		`
		requestQuery(query, nil)
	}
	// add comment
	if *idOfTheGameToCommentAbout != -1 {
		query := `
			mutation AddComment($gameID: Int!, $commentBody: String!) {
				addComment(gameID: $gameID, commentBody: $commentBody) {
					body
				}
			}
		`
		requestQuery(query, map[string]interface{}{"gameID": *idOfTheGameToCommentAbout, "commentBody": *commentBody})
	}
}
