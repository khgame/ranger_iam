package authcli

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/khgame/ranger_iam/internal/util"
)

// AuthN 从 jwt 中直接获取 user_id 信息
func (cli *Cli) AuthN(ctx context.Context, tokenStr string) (uint, error) {
	uid, err := cli.ValidateRemote(ctx, tokenStr)
	if err == nil {
		return uid, nil
	}
	if errors.Is(err, ErrValidateRemoteDegraded) || errors.Is(err, ErrValidateRemoteStatusFailed) {
		claims, e := cli.localJWT.ValidateClaims(tokenStr)
		if e != nil {
			return 0, e
		}
		return claims.UID, nil
	}
	return 0, err
}

// ValidateRemote tries to validate the token with the IAM server first.
// If the IAM server is unavailable or gives a degraded response,
// it performs local JWT validation (long term ticket).
// Otherwise, it opts for the short ticket.
func (cli *Cli) ValidateRemote(ctx context.Context, token string) (uint, error) {
	// Try to contact the IAM server
	response, err := cli.httpClient.Get(cli.AuthNSvrURL + "api/v1/session/validate")
	if err != nil || response.StatusCode != http.StatusOK {
		// Server is down or returned a non-ok status - use local validation
		return 0, ErrValidateRemoteStatusFailed
	}

	// Check if the response includes a degradation signal
	// Let's assume the header X-Degraded-Mode indicates the mode
	if response.Header.Get(util.KEYDegradedMode) == util.DegradedModeAll {
		return 0, ErrValidateRemoteDegraded
	}

	// Server response can be used to get userID
	defer response.Body.Close()
	var res struct {
		UID uint `json:"uid"`
	}
	if err = json.NewDecoder(response.Body).Decode(&res); err != nil {
		return 0, err
	}
	return res.UID, nil
}
