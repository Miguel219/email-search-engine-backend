package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	types "email-search-engine-backend/internal/server/types"
)

const index = "emails"

func IndexEmails(data []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", getUrl("_bulk"), strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	req = setHeaders(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func ListEmails(from string, max_results string) (*types.ZincSearchResponse, error) {
	query := fmt.Sprintf(`{
        "search_type": "alldocuments",
        "query":
        {
			"term": "",
            "start_time": "2023-01-01T19:14:45-06:00"
        },
		"sort_fields": ["-@timestamp"],
        "from": %s,
        "max_results": %s,
        "_source": []
    }`, from, max_results)
	req, err := http.NewRequest("POST", getUrlByIndex(index, "_search"), strings.NewReader(query))
	if err != nil {
		return nil, err
	}
	req = setHeaders(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := new(types.ZincSearchResponse)
	json.NewDecoder(resp.Body).Decode(response)
	return response, nil
}

func SearchEmails(term string, from string, max_results string) (*types.ZincSearchResponse, error) {
	query := fmt.Sprintf(`{
        "search_type": "matchphrase",
        "query":
        {
            "term": "%s",
            "start_time": "2023-01-01T19:14:45-06:00"
        },
		"sort_fields": ["-@timestamp"],
        "from": %s,
        "max_results": %s,
        "_source": []
    }`, term, from, max_results)

	req, err := http.NewRequest("POST", getUrlByIndex(index, "_search"), strings.NewReader(query))
	if err != nil {
		return nil, err
	}
	req = setHeaders(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := new(types.ZincSearchResponse)
	json.NewDecoder(resp.Body).Decode(response)
	return response, nil
}
