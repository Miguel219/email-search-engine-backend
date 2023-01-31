package server

import (
	"encoding/json"
	"net/http"

	services "email-search-engine-backend/internal/server/services"
	types "email-search-engine-backend/internal/server/types"

	"github.com/go-chi/render"
)

func ListEmails(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	from := r.Form.Get("from")
	max_results := r.Form.Get("max_results")

	response, err := services.ListEmails(from, max_results)
	if err != nil {
		types.ErrInvalidRequest(err)
		return
	}

	render.JSON(w, r, zincSearchResponseToEmailsResponse(response))
}

// SearchEmails searches the Emails data for a matching email.
// It's just a stub, but you get the idea.
func SearchEmails(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	term := r.Form.Get("term")
	from := r.Form.Get("from")
	max_results := r.Form.Get("max_results")

	response, err := services.SearchEmails(term, from, max_results)
	if err != nil {
		types.ErrInvalidRequest(err)
		return
	}

	render.JSON(w, r, zincSearchResponseToEmailsResponse(response))

}

func zincSearchResponseToEmailsResponse(response *types.ZincSearchResponse) *types.EmailsResponse {

	var emails []types.Email

	for _, hit := range response.Hits.Hits {
		var email types.Email
		emailBytes, _ := json.Marshal(hit.Source)
		json.Unmarshal(emailBytes, &email)
		email.Id = hit.ID
		emails = append(emails, email)
	}

	return &types.EmailsResponse{
		Emails: emails,
		Total:  response.Hits.Total.Value,
	}
}
