package resizer

import (
	"encoding/json"
	"testing"

	"gitlab.com/luizbranco/cyberbrain/primitives"
	"gitlab.com/luizbranco/cyberbrain/test"
	"gitlab.com/luizbranco/cyberbrain/test/mocks"
)

func TestWorker_Register(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		pool := &mocks.WorkerPool{}

		w := &Worker{
			WorkerPool: pool,
			BlitlineID: "123",
			AWSBucket:  "example-bucket",
		}

		pool.RegisterFunc = func(name string, worker primitives.Worker) error {
			test.Equal(t, "worker name", workerName, name)
			test.Equal(t, "worker", worker, w)
			return nil
		}

		err := w.Register()
		test.OK(t, err)
	})

	t.Run("worker pool is not defined", func(t *testing.T) {
		w := &Worker{}

		err := w.Register()
		test.Error(t, err)
	})

	t.Run("worker fails to enqueue", func(t *testing.T) {
		pool := &mocks.WorkerPool{}

		w := &Worker{
			WorkerPool: pool,
		}

		err := w.Register()
		test.Error(t, err)
	})
}

func TestWorker_Resize(t *testing.T) {
	imager := primitives.Card{
		ImageURL: "https://placeimg.com/400/300",
	}

	t.Run("ok", func(t *testing.T) {
		pool := &mocks.WorkerPool{}

		w := &Worker{
			WorkerPool:  pool,
			CallbackURL: "http://www.example.com/blitline",
			BlitlineID:  "fake-app-id",
			AWSBucket:   "example-bucket",
		}

		pool.EnqueueFunc = func(name string, v interface{}) error {
			test.Equal(t, "worker name", workerName, name)

			args, ok := v.(JobArgs)
			if !ok {
				t.Fatalf("invalid job args %v", v)
			}

			exp := JobArgs{
				ImagerID:    "AB34",
				ImageURL:    "https://placeimg.com/400/300",
				S3Path:      "cards/AB34.png",
				CallbackURL: "http://www.example.com/blitline/cards/AB34",
				Width:       400,
				Height:      300,
			}

			test.Equal(t, "job args", exp, args)

			return nil
		}

		err := w.Resize(&imager, "AB34", 400, 300)
		test.OK(t, err)
	})

	t.Run("worker pool not defined", func(t *testing.T) {
		w := &Worker{}

		err := w.Resize(&imager, "AB34", 400, 300)
		test.Error(t, err)
	})

	t.Run("worker fails to enqueue", func(t *testing.T) {
		pool := &mocks.WorkerPool{}

		w := &Worker{
			WorkerPool: pool,
		}

		err := w.Resize(&imager, "AB34", 400, 300)
		test.Error(t, err)
	})
}

func TestWorker_Spawn(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		args := JobArgs{
			ImagerID:    "AB34",
			ImageURL:    "https://placeimg.com/400/300",
			S3Path:      "cards/AB34.png",
			CallbackURL: "http://www.example.com/blitline/cards/AB34",
			Width:       400,
			Height:      300,
		}

		b, err := json.Marshal(args)
		test.OK(t, err)

		w := &Worker{
			BlitlineID: "fake-app-id",
			AWSBucket:  "example-bucket",
		}

		job, err := w.Spawn(b)

		exp := &Job{
			args:           args,
			awsBucket:      "example-bucket",
			blitlineAPIURL: blitlineAPIURL,
			blitlineID:     "fake-app-id",
		}

		test.Equal(t, "job with args", exp, job)
		test.OK(t, err)
	})

	t.Run("invalid args", func(t *testing.T) {
		args := []byte("")

		w := &Worker{}

		job, err := w.Spawn(args)

		test.Equal(t, "no job returned", nil, job)
		test.Error(t, err)
	})
}
