package blitline

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gitlab.com/luizbranco/cyberbrain/primitives"
	"gitlab.com/luizbranco/cyberbrain/web"
	"gitlab.com/luizbranco/cyberbrain/web/server/finder"
	"gitlab.com/luizbranco/cyberbrain/web/server/response"
)

// BlitlineResponse is the json response from Blitline postback or polling
type BlitlineResponse struct {
	Results struct {
		JobID  string `json:"job_id"`
		Images []struct {
			ImageIdentifier string `json:"image_identifier"`
			S3URL           string `json:"s3_url"`
		} `json:"images"`
	} `json:"results"`
}

func PatchCard(conn primitives.Database, ub web.URLBuilder, hash string) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		card, _, err := finder.Card(conn, ub, hash, finder.NoOption)
		if err != nil {
			return err.(response.Error)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return response.WrapError(err, http.StatusBadRequest, "failed to read body")
		}

		var res BlitlineResponse

		err = json.Unmarshal(body, &res)
		if err != nil {
			return response.WrapError(err, http.StatusBadRequest, "failed to unmarshal json")
		}

		if len(res.Results.Images) == 0 {
			return response.NewError(http.StatusBadRequest, "no images result")
		}

		card.ImageURL = res.Results.Images[0].S3URL

		err = conn.Update(card)
		if err != nil {
			return response.WrapError(err, http.StatusBadRequest, "failed to update card")
		}

		page := web.Page{
			Title:    "Card Updated",
			Partials: []string{"200"},
		}

		return response.NewContent(page)
	}
}
