package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func example() {
	s := NewServer()

	// Query
	query := `
		mutation Create($game: GameInput!) {
			put(game: $game, password: "1234") {
				id
				roomName
			}
		}
	`
	res, err := s.doQuery(query, map[string]interface{}{"game": map[string]interface{}{"id": 0, "roomName": "Test room", "status": "TestStatus"}})
	fmt.Printf("%s %v\n", res, err)

	// Query
	query = `
		mutation AddComment($gameID: Int!, $commentBody: String!) {
			addComment(gameID: $gameID, commentBody: $commentBody) {
				body
			}
		}
	`
	res, err = s.doQuery(query, map[string]interface{}{"gameID": 0, "commentBody": "Test comment"})
	fmt.Printf("%s %v\n", res, err)

	// Query
	query = `
		{
			list {
				id
				roomName
				status
				startedAt
				isDay
				day
				comments {
					body
				}
			}
		}
	`
	res, err = s.doQuery(query, nil)
	fmt.Printf("%s %v\n", res, err)

	// Query
	query = `
		query GetGame($id: Int!){
			game(id: $id) {
				id
				roomName
				status
				startedAt
				isDay
				day
				players {
					login
					isAlive
					deadCause
				}
				comments {
					body
				}
			}
		}
	`
	res, err = s.doQuery(query, map[string]interface{}{"id": 1})
	fmt.Printf("%s %v\n", res, err)

	query = `
	{
		"query": "query GetGame($id: Int!) { game(id: $id) { id roomName } }",
		"variables": {"id": 0}
	}
	`
	var v map[string]interface{}
	err = json.Unmarshal([]byte(query), &v)
	if err != nil {
		log.Fatalf("unmarshal err: %v", err)
	}
	res, err = s.doQuery(v["query"].(string), v["variables"].(map[string]interface{}))
	fmt.Printf("%s %v\n", res, err)

	res, err = s.doJSONQueryWithVariables(query)
	fmt.Printf("%s %v\n", res, err)
}
