package main

import "github.com/graphql-go/graphql"

// go types

type Player struct {
	Login     string
	IsAlive   bool
	DeadCause string
	Role      string
}

type Game struct {
	ID         int
	RoomName   string
	Status     string
	StartedAt  string
	FinishedAt string
	IsDay      bool
	Day        int
	Players    []Player
	Comments   []Comment
}

type Comment struct {
	Body string
}

// copy game helper function
func copyGame(game Game) Game {
	res := game
	res.Players = append(make([]Player, 0, len(game.Players)), game.Players...)
	res.Comments = append(make([]Comment, 0, len(game.Comments)), game.Comments...)
	return res
}

// helper function to get list of games from map
func ListOfGamesFromMap(games map[int]*Game) []Game {
	res := make([]Game, 0)
	for _, v := range games {
		res = append(res, copyGame(*v))
	}
	return res
}

// map of games example
func populate() map[int]*Game {
	game := Game{
		ID:         1,
		RoomName:   "room1",
		Status:     "Active",
		StartedAt:  "2022-05-11 19:41 +03:00",
		FinishedAt: "-",
		IsDay:      true,
		Day:        1,
		Players: []Player{
			{Login: "123", IsAlive: true, DeadCause: "-", Role: "unknown"},
			{Login: "456", IsAlive: true, DeadCause: "-", Role: "unknown"},
			{Login: "789", IsAlive: true, DeadCause: "-", Role: "unknown"},
			{Login: "101", IsAlive: false, DeadCause: "Killed by mafia", Role: "unknown"},
		},
		Comments: []Comment{
			{Body: "First Comment"},
		},
	}
	return map[int]*Game{1: &game}
}

// GraphQL types

var playerType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Player",
		Fields: graphql.Fields{
			"login": &graphql.Field{
				Type: graphql.String,
			},
			"isAlive": &graphql.Field{
				Type: graphql.Boolean,
			},
			"deadCause": &graphql.Field{
				Type: graphql.String,
			},
			"role": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var commentType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Comment",
		Fields: graphql.Fields{
			"body": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var gameType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Game",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"roomName": &graphql.Field{
				Type: graphql.String,
			},
			"status": &graphql.Field{
				Type: graphql.String,
			},
			"startedAt": &graphql.Field{
				Type: graphql.String,
			},
			"finishedAt": &graphql.Field{
				Type: graphql.String,
			},
			"isDay": &graphql.Field{
				Type: graphql.Boolean,
			},
			"day": &graphql.Field{
				Type: graphql.Int,
			},
			"players": &graphql.Field{
				Type: graphql.NewList(playerType),
			},
			"comments": &graphql.Field{
				Type: graphql.NewList(commentType),
			},
		},
	},
)

var gameInputType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "GameInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"id": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"roomName": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"status": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"startedAt": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"finishedAt": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"isDay": &graphql.InputObjectFieldConfig{
				Type: graphql.Boolean,
			},
			"day": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"players": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(playerType),
			},
			"comments": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(commentType),
			},
		},
	},
)

// mutation
func getMutation(s *Server) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"put": &graphql.Field{
				Type:        gameType,
				Description: "Create or update game",
				Args: graphql.FieldConfigArgument{
					"game": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(gameInputType),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: s.putGame,
			},
			"addComment": &graphql.Field{
				Type:        commentType,
				Description: "Create or update game",
				Args: graphql.FieldConfigArgument{
					"gameID": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"commentBody": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: s.addComment,
			},
		},
	})
}

// query
func getQuery(s *Server) *graphql.Object {
	fields := graphql.Fields{
		"game": &graphql.Field{
			Type:        gameType,
			Description: "Get Game By ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: s.getGame,
		},
		"list": &graphql.Field{
			Type:        graphql.NewList(gameType),
			Description: "Get Games List",
			Resolve:     s.getGamesList,
		},
	}
	return graphql.NewObject(graphql.ObjectConfig{Name: "RootQuery", Fields: fields})
}

// schema
func getSchema(s *Server) (graphql.Schema, error) {
	schemaConfig := graphql.SchemaConfig{
		Query:    getQuery(s),
		Mutation: getMutation(s),
	}
	return graphql.NewSchema(schemaConfig)
}
