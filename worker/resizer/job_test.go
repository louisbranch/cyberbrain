package resizer

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/luizbranco/cyberbrain/test"
)

func TestJob_Run(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()

			test.OK(t, err)

			payload := `{"application_id":"fake-app-id","src":"https://placeimg.com/400/300","v":1.21,"postback_url":"http://www.example.com/blitline/cards/AB34","functions":[{"name":"resize_to_fill","params":{"height":300,"width":400},"save":{"extension":".png","image_identifier":"AB34","s3_destination":{"bucket":"example-bucket","key":"cards/AB34.png"}}}]}`

			test.Equal(t, "payload", payload, string(b))

			w.WriteHeader(200)
			w.Write([]byte("received"))
		}))

		args := JobArgs{
			ImagerID:    "AB34",
			ImageURL:    "https://placeimg.com/400/300",
			S3Path:      "cards/AB34.png",
			CallbackURL: "http://www.example.com/blitline/cards/AB34",
			Width:       400,
			Height:      300,
		}

		job := &Job{
			args:           args,
			awsBucket:      "example-bucket",
			blitlineAPIURL: srv.URL,
			blitlineID:     "fake-app-id",
		}

		ctx := context.Background()

		err := job.Run(ctx)
		test.OK(t, err)
	})

	t.Run("failed blitline response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(400)
			w.Write([]byte("invalid arguments"))
		}))

		job := &Job{
			blitlineAPIURL: srv.URL,
		}

		ctx := context.Background()

		err := job.Run(ctx)
		test.Error(t, err, "failed to send request to Blitline [400] invalid arguments")
	})
}
