package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type book struct {
	ISBN   string `json:"isbn"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// Response ...
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler() (Response, error) {
	bk := &book{
		ISBN:   "978-1420931693",
		Title:  "The Republic",
		Author: "Plato",
	}

	response, _ := json.Marshal(map[string]*book{"data": bk})

	resp := Response{
		StatusCode: 200,
		Body:       string(response),
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
