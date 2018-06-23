package resizer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type JobArgs struct {
	ImagerID    string
	CallbackURL string
	Poll        bool
	ImageURL    string
	S3Path      string
	Width       int
	Height      int
}

type Job struct {
	args           JobArgs
	awsBucket      string
	blitlineAPIURL string
	blitlineID     string
}

func (j *Job) Run(ctx context.Context) error {
	payload := BlitlineRequest{
		ApplicationID: j.blitlineID,
		Version:       blitlineVersion,
		ImageURL:      j.args.ImageURL,
		Functions: []BlitlineFunction{
			{
				Name: "resize_to_fill",
				Params: map[string]int{
					"width":  j.args.Width,
					"height": j.args.Height,
				},
				Save: map[string]interface{}{
					"extension":        ".png",
					"image_identifier": j.args.ImagerID,
					"s3_destination": map[string]string{
						"key":    j.args.S3Path,
						"bucket": j.awsBucket,
					},
				},
			},
		},
	}

	if !j.args.Poll {
		payload.CallbackURL = j.args.CallbackURL
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "failed to marshal payload")
	}

	req, err := http.NewRequest("POST", j.blitlineAPIURL, bytes.NewBuffer(b))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send request to Blitline")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("failed to send request to Blitline [%d] %s", resp.StatusCode, body)
	}

	if !j.args.Poll {
		return nil
	}

	var r BlitlineResponse

	err = json.Unmarshal(body, &r)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal Blitline results %s", body)
	}

	jobID := r.Results.JobID

	return poll(j.args.CallbackURL, jobID, 5*time.Second)
}

func poll(callbackURL, jobID string, freq time.Duration) error {
	url := fmt.Sprintf("%s/%s", blitlinePollURL, jobID)

	deadline := time.NewTicker(pollDeadline)
	defer deadline.Stop()

	t := time.NewTicker(freq)
	defer t.Stop()

	var respErr error

Loop:
	for {
		select {
		case <-deadline.C:
			return errors.Errorf("failed to poll Blitline job %q, deadline reached %q", jobID, respErr)
		case <-t.C:
			resp, err := http.Get(url)
			if err != nil {
				return errors.Wrapf(err, "failed to poll Blitline job %q", jobID)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return errors.Wrap(err, "failed to read response body")
			}

			if resp.StatusCode != http.StatusOK {
				respErr = errors.Errorf("failed to send request to Blitline [%d] %s", resp.StatusCode, body)
				resp.Body.Close()
				continue
			}

			req, err := http.NewRequest("POST", callbackURL, bytes.NewReader(body))
			if err != nil {
				return errors.Wrap(err, "failed to create request")
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {
				return errors.Wrap(err, "failed to send request to callback url")
			}
			defer res.Body.Close()

			break Loop
		}
	}

	return nil
}
