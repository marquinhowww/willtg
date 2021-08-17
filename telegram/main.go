package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	c "github.com/eptaccio/willtg/commands"
	t "github.com/eptaccio/willtg/types"
	tb "gopkg.in/tucnak/telebot.v2"
)

// Response on action
type Response events.APIGatewayProxyResponse

// Request action
type Request events.APIGatewayProxyRequest

// Handler http requests
func Handler(ctx context.Context, req Request) (Response, error) {
	b, err := tb.NewBot(tb.Settings{
		Token:       os.Getenv("BOT_TOKEN"),
		Synchronous: true,
	})

	if err != nil {
		log.Fatal(err)
	}

	for _, command := range c.CommandMapping {
		func(config t.ImageConfig) {
			b.Handle(config.Command, func(m *tb.Message) {
				filePath, err := c.WriteImage(m.Text, fmt.Sprint(m.ID), config)
				if err != nil {
					b.Send(m.Sender, err.Error())
					return
				}

				photo := &tb.Photo{File: tb.FromDisk(*filePath)}
				_, err = b.SendAlbum(m.Chat, tb.Album{photo})
				if err != nil {
					b.Send(m.Sender, err.Error())
				}
			})
		}(command)
	}

	var u tb.Update
	if err = json.Unmarshal([]byte(req.Body), &u); err == nil {
		b.ProcessUpdate(u)
	}

	resp := Response{
		StatusCode: 200,
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
