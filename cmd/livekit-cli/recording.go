package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/ggwhite/go-masker"
	livekit "github.com/livekit/protocol/proto"
	lksdk "github.com/livekit/server-sdk-go"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	RecordCommands = []*cli.Command{
		{
			Name:   "start-recording",
			Usage:  "starts a recording with a deployed recorder service",
			Before: createRecordingClient,
			Action: startRecording,
			Flags: []cli.Flag{
				urlFlag,
				apiKeyFlag,
				secretFlag,
				&cli.StringFlag{
					Name:     "request",
					Usage:    "StartRecordingRequest as json file (see https://github.com/livekit/protocol/blob/main/livekit_recording.proto#L16)",
					Required: true,
				},
			},
		},
		{
			Name:   "end-recording",
			Before: createRecordingClient,
			Action: endRecording,
			Flags: []cli.Flag{
				urlFlag,
				apiKeyFlag,
				secretFlag,
				&cli.StringFlag{
					Name:     "id",
					Usage:    "id of the recording",
					Required: true,
				},
			},
		},
	}

	recordingClient *lksdk.RecordingServiceClient
)

func createRecordingClient(c *cli.Context) error {
	url := c.String("url")
	apiKey := c.String("api-key")
	apiSecret := c.String("api-secret")

	if c.Bool("verbose") {
		fmt.Printf("creating client to %s, with api-key: %s, secret: %s\n",
			url,
			masker.ID(apiKey),
			masker.ID(apiSecret))
	}

	recordingClient = lksdk.NewRecordingServiceClient(url, apiKey, apiSecret)
	return nil
}

func startRecording(c *cli.Context) error {
	reqFile := c.String("request")
	reqBytes, err := ioutil.ReadFile(reqFile)
	if err != nil {
		return err
	}
	req := &livekit.StartRecordingRequest{}
	err = protojson.Unmarshal(reqBytes, req)
	if err != nil {
		return err
	}

	if c.Bool("verbose") {
		PrintJSON(req)
	}

	resp, err := recordingClient.StartRecording(context.Background(), req)
	if err != nil {
		return err
	}

	fmt.Printf("Recording started. Recording ID: %s\n", resp.RecordingId)
	return nil
}

func endRecording(c *cli.Context) error {
	_, err := recordingClient.EndRecording(context.Background(), &livekit.EndRecordingRequest{
		RecordingId: c.String("id"),
	})
	return err
}