package resizer

import "time"

const (
	blitlineVersion = 1.21
	blitlineAPIURL  = "https://api.blitline.com/job"
	blitlinePollURL = "http://cache.blitline.com/listen"
	pollDeadline    = 5 * time.Minute
)

// BlitlineRequest is the json request send to Blitline
type BlitlineRequest struct {
	ApplicationID string             `json:"application_id"`
	ImageURL      string             `json:"src"`
	Version       float64            `json:"v"`
	CallbackURL   string             `json:"postback_url,omitempty"`
	Functions     []BlitlineFunction `json:"functions"`
}

// BlitlineFunction represents a Blitline function
type BlitlineFunction struct {
	Name   string                 `json:"name"`
	Params map[string]int         `json:"params"`
	Save   map[string]interface{} `json:"save"`
}

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
