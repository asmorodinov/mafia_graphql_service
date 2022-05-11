package main

import (
	"bytes"
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
	return string(respBody), resp.Status, resp.StatusCode, err
}

func main() {
	flag.Parse()

	bodyStr := `
	{
		"query": "query GetGame($id: Int!) { game(id: $id) { id roomName } }",
		"variables": {"id": 1}
	}
	`
	resp, status, code, err := request(bodyStr)
	fmt.Printf("%v\n%v\n%v\n%v\n", resp, status, code, err)
}
