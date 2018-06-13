package worker

import "gitlab.com/luizbranco/cyberbrain/primitives"

type query struct{}

func (q *query) NewRecord() primitives.Record {
	return &Job{}
}

func (q *query) Where() map[string]interface{} {
	return map[string]interface{}{
		"state": scheduled,
	}
}

func (q *query) Raw() string {
	return ""
}

func (q *query) SortBy() map[string]string {
	return nil
}
