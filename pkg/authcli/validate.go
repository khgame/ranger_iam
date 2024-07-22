package authcli

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/bagaking/goulp/wlog"
	"github.com/khgame/ranger_iam/internal/utils"
	"github.com/khgame/ranger_iam/pkg/auth"
	"github.com/khicago/irr"
)

// AuthN 从 jwt 中直接获取 user_id 信息
func (cli *Cli) AuthN(ctx context.Context, tokenStr string) (uint64, error) {
	uid, err := cli.ValidateRemote(ctx, tokenStr)
	if err == nil {
		return uid, nil
	}
	if errors.Is(err, ErrValidateRemoteDegraded) || errors.Is(err, ErrValidateRemoteStatusFailed) {
		wlog.ByCtx(ctx, "AuthN").WithError(err).Warn("AuthN failed, fallback to local JWT validation")
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
func (cli *Cli) ValidateRemote(ctx context.Context, token string) (uint64, error) {
	// Try build the full url
	baseURL, err := url.Parse(cli.AuthNSvrURL)
	if err != nil {
		return 0, irr.Wrap(err, "parse base url failed")
	}
	relatedURL, err := url.Parse("api/v1/session/validate")
	if err != nil {
		return 0, irr.Wrap(err, "parse related url failed")
	}

	// Build Request and pass the bearer token
	fullURL := baseURL.ResolveReference(relatedURL).String()
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return 0, irr.Wrap(err, "create request failed")
	}

	// Try to contact the IAM server
	auth.SetTokenStrToHeader(req, token)
	response, err := cli.httpClient.Do(req)
	if err != nil {
		wlog.ByCtx(ctx, "authN").WithError(err).Debugf("remote request failed")
		// Server is down or returned a non-ok status - use local validation
		return 0, irr.Wrap(ErrValidateRemoteStatusFailed, "got err= %s", err)
	} else if response.StatusCode != http.StatusOK {
		wlog.ByCtx(ctx, "authN").Debugf("remote invalid status, status= %s", response.Status)
		return 0, irr.Wrap(ErrValidateRemoteStatusFailed, "failed status= %s", response.Status)
	}

	// Check if the response includes a degradation signal
	// Let's assume the header X-Degraded-Mode indicates the mode
	if response.Header.Get(utils.KEYDegradedMode) == utils.DegradedModeAll {
		wlog.ByCtx(ctx, "authN").Infof("downgrade by remote")
		return 0, ErrValidateRemoteDegraded
	}

	// Server response can be used to get userID
	defer response.Body.Close()
	var res struct {
		UID uint64 `json:"uid"`
	}
	if err = json.NewDecoder(response.Body).Decode(&res); err != nil {
		return 0, err
	}
	return res.UID, nil
}
