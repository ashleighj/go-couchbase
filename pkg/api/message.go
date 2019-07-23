package api

import (
	"context"
	"gocouchbase/pkg/storage/couchbase"
	"net/http"

	"github.com/go-chi/render"
)

type ApiResponse struct {
	Errors []string    `json:"errors"`
	Data   interface{} `json:"data"`
}

func (response *ApiResponse) AddError(err error) {
	response.Errors = append(response.Errors, err.Error())
}

func RespondWithError(w http.ResponseWriter, r *http.Request, httpStatusCode int, errors ...string) {
	response := ApiResponse{Errors: errors}
	response.Respond(w, r, httpStatusCode)
}

func (response *ApiResponse) Respond(w http.ResponseWriter, r *http.Request, responseStatus int) {
	render.Status(r, responseStatus)
	render.JSON(w, r, response)
}

type serviceAuthRequest struct {
	ClientID string `json:"client_id"`
	APIKey   string `json:"api_key"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

type registerResetClientRequest struct {
	ClientID string `json:"client_id"`
	Role     string `json:"role"`
}

type clientInfoResponseData struct {
	ClientID string `json:"client_id"`
	APIKey   string `json:"api_key"`
	Role     string `json:"role"`
}

func (data *clientInfoResponseData) save(ctx context.Context) (err error) {
	// add to couchbse
	if data.Role == "" {
		// reset (existing client) - fetch first
		client, err := Service.GetAPIClient(data.ClientID)
		if err != nil {
			return err
		}
		client.APIKey = data.APIKey
		return Service.StoreAPIClient(client)

	} else {
		// register (new client)
		client := couchbase.APIClient{
			ClientID: data.ClientID,
			APIKey:   data.APIKey,
			Role:     data.Role}

		return Service.StoreAPIClient(client)
	}
	return
}
