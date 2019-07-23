package api

import (
	"context"
	"encoding/json"
	"gocouchbase/pkg/log"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func decodeJsonRequest(r *http.Request, valueObject interface{}) error {
	if valueObject == nil {
		return log.GetErrorf(log.ErrRequiredParamNil, "valueObject")
	}
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(valueObject)
}

func hashAndSalt(str string) (hash string, err error) {
	h, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.MinCost)
	return string(h), err
}

func hashMatchesString(unhashed string, hashed string) (matches bool, err error) {
	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(unhashed))
	if err != nil {
		return false, err
	}

	return true, err
}

func generateAPIKey() string {
	rand.Seed(time.Now().UnixNano())

	length := 32
	runes := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	byteArr := make([]byte, length)
	for i := 0; i < length; i++ {
		byteArr[i] = runes[rand.Intn(len(runes))]
	}

	rand.Shuffle(len(byteArr), func(i, j int) {
		byteArr[i], byteArr[j] = byteArr[j], byteArr[i]
	})

	return string(byteArr)
}

func isRoleValid(ctx context.Context, roleName string) bool {
	roles, err := Service.GetServiceRoles()
	if err != nil {
		log.Error(ctx, err)
	}
	if roles == nil {
		log.Error(ctx, log.ErrStorageNoRoles)
	}

	for _, role := range roles {
		if role.RoleName == roleName {
			return true
		}
	}
	return false
}
