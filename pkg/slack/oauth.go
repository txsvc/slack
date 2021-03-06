package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/context"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/txsvc/commons/pkg/util"
	"github.com/txsvc/platform/pkg/platform"
	"github.com/txsvc/service/pkg/auth"
)

type (
	// OAuthResponse is used to give a simple reponse to the user as feedback to a custom action reuqest
	OAuthResponse struct {
		OK              bool            `json:"ok,omitempty"`
		AccessToken     string          `json:"access_token,omitempty"`
		TokenType       string          `json:"token_type,omitempty"`
		Scope           string          `json:"scope,omitempty"`
		AppID           string          `json:"app_id,omitempty"`
		BotUserID       string          `json:"bot_user_id,omitempty"`
		IncomingWebhook *WebhookElement `json:"incoming_webhook,omitempty"`
	}
)

// OAuthEndpoint handles the callback from Slack with the temporary access code
// and exchanges it with the real auth token. See https://api.slack.com/docs/oauth
func OAuthEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	// extract parameters
	code := c.Query("code")
	redirectURI := c.Query("redirect_uri")

	// FIXME: secure the request by using a state
	// state := c.Query("state")

	if code != "" {
		// exchange the temporary code with a real auth token
		resp, err := getOAuthToken(ctx, code)

		if err != nil {
			platform.ReportError(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error")
			return
		}

		// get team info
		var teamInfo TeamInfo
		err = Get(ctx, resp.AccessToken, "team.info", "", &teamInfo)
		if err != nil {
			platform.ReportError(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error")
			return
		}

		if teamInfo.OK == false {
			platform.ReportError(fmt.Errorf("error: %s", teamInfo.Error))
			c.Redirect(http.StatusTemporaryRedirect, "/error")
			return
		}

		err = UpdateAuthorization(ctx, teamInfo.Team.ID, teamInfo.Team.Name, resp.AccessToken, resp.TokenType, resp.Scope, resp.AppID, resp.BotUserID)
		if err != nil {
			platform.ReportError(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error")
			return
		}
	}

	// back to the sign-up process on the main website
	if redirectURI == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	} else {
		c.Redirect(http.StatusTemporaryRedirect, redirectURI)
	}
}

// getOAuthToken exchanges a temporary OAuth verifier code for an access token
func getOAuthToken(ctx context.Context, code string) (*OAuthResponse, error) {

	url := SlackEndpoint + "oauth.v2.access?code=" + code

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(os.Getenv(SlackClientID), os.Getenv(SlackClientSecret))

	// post the request to Slack
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// unmarshal the response
	var response OAuthResponse
	err = json.NewDecoder(resp.Body).Decode(&response)

	return &response, err
}

// UpdateAuthorization updates the authorization, or creates a new one.
func UpdateAuthorization(ctx context.Context, clientID, teamName, token, tokenType, scope, appID, botID string) error {

	// find the authorization first
	a, err := auth.GetAuthorization(ctx, clientID, auth.AuthTypeSlack)

	now := util.Timestamp()
	if err == nil {
		a.Token = token
		a.Scope = scope
		a.Updated = now
	} else {
		a = &auth.Authorization{
			ClientID:  clientID, // TeamID
			Name:      teamName, // Team name
			Token:     token,
			TokenType: tokenType,
			UserID:    botID,
			Scope:     scope,
			Expires:   0,
			// internal
			AuthType: auth.AuthTypeSlack,
			Created:  now,
			Updated:  now,
		}
	}

	return auth.CreateAuthorization(ctx, a)
}
