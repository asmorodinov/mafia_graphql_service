package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

type Server struct {
	mutex  sync.Mutex
	games  map[int]*Game
	schema graphql.Schema
}

func NewServer() *Server {
	s := &Server{sync.Mutex{}, make(map[int]*Game), graphql.Schema{}}

	// get schema
	schema, err := getSchema(s)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	s.schema = schema

	return s
}

func (s *Server) putGame(params graphql.ResolveParams) (interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	password := params.Args["password"].(string)
	if password != "1234" {
		return Game{}, fmt.Errorf("Unauthorized")
	}

	gameMap := params.Args["game"].(map[string]interface{})
	id := gameMap["id"].(int)

	var game *Game
	if v, ok := s.games[id]; ok {
		game = v
	} else {
		game = &Game{}
	}
	updateGame(game, gameMap)
	s.games[id] = game

	return copyGame(*game), nil
}

func (s *Server) addComment(params graphql.ResolveParams) (interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id := params.Args["gameID"].(int)
	comment := Comment{params.Args["commentBody"].(string)}

	if _, ok := s.games[id]; !ok {
		return Comment{}, fmt.Errorf("game with id %v not found", id)
	}

	game := s.games[id]
	game.Comments = append(game.Comments, comment)

	return comment, nil
}

func (s *Server) getGame(p graphql.ResolveParams) (interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id, ok := p.Args["id"].(int)

	if !ok {
		return Game{}, fmt.Errorf("id is not an int")
	}
	game, ok := s.games[id]
	if !ok {
		return Game{}, fmt.Errorf("game with id %v not found", id)
	}
	return game, nil
}

func (s *Server) getGamesList(p graphql.ResolveParams) (interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return ListOfGamesFromMap(s.games), nil
}

func (s *Server) doQuery(query string, variables map[string]interface{}) (string, error) {
	params := graphql.Params{Schema: s.schema, RequestString: query, VariableValues: variables}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		return "", fmt.Errorf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	return string(rJSON), nil
}

func (s *Server) doJSONQueryWithVariables(query string) (string, error) {
	var v map[string]interface{}
	err := json.Unmarshal([]byte(query), &v)
	if err != nil {
		return "", fmt.Errorf("unmarshal err: %v", err)
	}
	res, err := s.doQuery(v["query"].(string), v["variables"].(map[string]interface{}))
	return res, err
}

func (s *Server) graphqlHandler(c *gin.Context) {
	// parse request json body
	type Request struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}
	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	// unescape query
	newstr, err := strconv.Unquote("\"" + request.Query + "\"")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// do query
	res, err := s.doQuery(newstr, request.Variables)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	// return result
	c.Data(http.StatusOK, gin.MIMEJSON, []byte(res))
}

var addr = flag.String("addr", "localhost:8090", "GraphQL server address")

func startServer() {
	flag.Parse()

	server := NewServer()
	router := gin.Default()
	router.POST("/graphql", server.graphqlHandler)
	router.Run(*addr)
}

func main() {
	// example()
	startServer()
}
