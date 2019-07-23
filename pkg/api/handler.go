package api

import (
	"fmt"
	"net/http"
	"time"

	"gocouchbase/pkg/auth"
	"gocouchbase/pkg/config"
	"gocouchbase/pkg/log"
	"gocouchbase/pkg/service"
	"gocouchbase/pkg/storage/couchbase"

	"github.com/go-chi/render"
)

var Service service.Service

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer log.TimeExecution(ctx, "api.handleHealthCheck", time.Now())

	response := make(map[string]string)
	response["status"] = "Alive and kickin'"

	render.JSON(w, r, response)
}

func handleServiceAuthenticate(w http.ResponseWriter, r *http.Request) {
	var request serviceAuthRequest
	var responseData tokenResponse
	response := ApiResponse{Data: &responseData}

	ctx := r.Context()
	defer log.TimeExecution(ctx, "api.handleServiceAuthenticate", time.Now())

	err := decodeJsonRequest(r, &request)
	if err != nil {
		log.Error(ctx, err)
		RespondWithError(w, r, http.StatusBadRequest, log.ErrBadRequest)
		return
	}

	client, err := Service.GetAPIClient(request.ClientID)
	if err != nil {
		log.Error(ctx, err)
		RespondWithError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if request.APIKey != client.APIKey {
		log.Error(ctx, log.ErrInvalidAPICredentials)
		RespondWithError(w, r, http.StatusUnauthorized, log.ErrInvalidAPICredentials)
		return
	}

	token, err := auth.GetEncodedToken(
		request.ClientID,
		time.Now().Add(time.Minute*time.Duration(config.Get().TokenExpiryMinutes)))

	if err != nil {
		log.Error(ctx, err)
		RespondWithError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	responseData.Token = token
	response.Respond(w, r, http.StatusOK)
}

func handleRegisterReset(w http.ResponseWriter, r *http.Request) {
	var request registerResetClientRequest
	var responseData clientInfoResponseData
	response := ApiResponse{Data: &responseData}

	ctx := r.Context()
	defer log.TimeExecution(ctx, "api.registerResetHandler", time.Now())

	err := decodeJsonRequest(r, &request)
	if err != nil {
		log.Error(ctx, err)
		RespondWithError(w, r, http.StatusBadRequest, log.ErrBadRequest)
		return
	}
	if request.Role != "" && !isRoleValid(ctx, request.Role) {
		log.Error(ctx, log.ErrRequestInvalidRole)
		RespondWithError(w, r, http.StatusBadRequest, log.ErrRequestInvalidRole)
		return
	}

	responseData.ClientID = request.ClientID
	responseData.APIKey = generateAPIKey()
	responseData.Role = request.Role

	err = responseData.save(ctx)
	if err != nil {
		log.Error(ctx, err)
		response.AddError(err)
	}

	response.Respond(w, r, http.StatusOK)
}

func handleRolesGet(w http.ResponseWriter, r *http.Request) {
	responseData := make(map[string][]couchbase.ServiceRole)
	response := ApiResponse{Data: &responseData}

	ctx := r.Context()
	defer log.TimeExecution(ctx, "api.handleRolesGet", time.Now())

	roles, err := Service.GetServiceRoles()
	if err != nil {
		log.Error(ctx, err)
		RespondWithError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	responseData["roles"] = roles
	response.Respond(w, r, http.StatusOK)
}

func handleRoleCreate(w http.ResponseWriter, r *http.Request) {
	var newRole couchbase.ServiceRole
	responseData := make(map[string]couchbase.ServiceRole)
	response := ApiResponse{Data: &responseData}

	ctx := r.Context()
	defer log.TimeExecution(ctx, "api.handleRoleCreate", time.Now())

	err := decodeJsonRequest(r, &newRole)
	if err != nil {
		log.Error(ctx, err)
		RespondWithError(w, r, http.StatusBadRequest, log.ErrBadRequest)
		return
	}

	if newRole.RoleName == "" {
		log.Error(ctx, log.ErrBadRequest)
		RespondWithError(w, r, http.StatusBadRequest, fmt.Sprintf(log.ErrRequiredParamNil, "role_name"))
		return
	}

	if len(newRole.AllowedPaths) < 1 {
		log.Error(ctx, log.ErrBadRequest)
		RespondWithError(w, r, http.StatusBadRequest, fmt.Sprintf(log.ErrRequiredParamNil, "role_paths"))
		return
	}

	err = Service.StoreServiceRole(newRole)
	if err != nil {
		log.Error(ctx, err)
		RespondWithError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	responseData["role"] = newRole
	response.Respond(w, r, http.StatusOK)
}
