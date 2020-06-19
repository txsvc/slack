package slack

import (
	"encoding/json"
	e "errors"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/txsvc/commons/pkg/env"
	"github.com/txsvc/commons/pkg/errors"
	"github.com/txsvc/commons/pkg/util"
	"github.com/txsvc/platform/pkg/platform"
)

const (
	// DatastoreAuthorizations collection AUTHORIZATION
	DatastoreAuthorizations string = "AUTHORIZATIONS"
)

type (
	// AuthorizationDS holds basic information about a Slack team/workspace and
	// the OAuth given at installtion time.
	AuthorizationDS struct {
		ID          string
		Name        string
		AccessToken string
		TokenType   string
		AppID       string
		BotUserID   string
		Scope       string
		// internal
		Created int64
		Updated int64
	}
)

// authorizationKey creates a datastore key for a workspace authorization based on the team_id.
func authorizationKey(id string) *datastore.Key {
	return datastore.NameKey(DatastoreAuthorizations, id, nil)
}

func cacheKey(id string) string {
	return "slack.auth." + id
}

// OAuthEndpoint handles the callback from Slack with the temporary access code
// and exchanges it with the real auth token. See https://api.slack.com/docs/oauth
func OAuthEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	// extract parameters
	code := c.Query("code")
	redirectURI := c.Query("redirect_uri")

	// FIXME secure the request by using a state
	// state := c.Query("state")

	if code != "" {
		// exchange the temporary code with a real auth token
		resp, err := getOAuthToken(ctx, code)

		if err != nil {
			platform.Report(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error")
			return
		}

		// get team info
		var teamInfo TeamInfo
		err = Get(ctx, resp.AccessToken, "team.info", "", &teamInfo)
		if err != nil {
			platform.Report(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error")
			return
		}

		if teamInfo.OK == false {
			platform.Report(errors.NewOperationError("team.info", e.New(teamInfo.Error)))
			c.Redirect(http.StatusTemporaryRedirect, "/error")
			return
		}

		err = UpdateAuthorization(ctx, teamInfo.Team.ID, teamInfo.Team.Name, resp.AccessToken, resp.TokenType, resp.Scope, resp.AppID, resp.BotUserID)
		if err != nil {
			platform.Report(err)
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

// GetAuthorization returns the authorization granted to an app
func GetAuthorization(ctx context.Context, id string) (*AuthorizationDS, error) {
	var auth = AuthorizationDS{}

	// just load it, let caching be handled elsewhere ...
	err := platform.DataStore().Get(ctx, authorizationKey(id), &auth)
	if err != nil {
		return nil, err
	}

	return &auth, nil
}

// GetAuthToken returns the oauth token of the workspace integration
func GetAuthToken(ctx context.Context, id string) (string, error) {
	// ENV always overrides anything stored...
	token := env.Getenv("SLACK_AUTH_TOKEN", "")
	if token != "" {
		return token, nil
	}

	// check the in-memory cache
	key := cacheKey(id)
	token, err := platform.Get(ctx, key)
	if token != "" {
		return token, nil
	}

	auth, err := GetAuthorization(ctx, id)
	if err != nil {
		return "", err
	}
	if auth == nil {
		return "", fmt.Errorf("No authorization token for workspace '%s'", id)
	}

	// add the token to the cache
	platform.Set(ctx, key, auth.AccessToken, 1800)

	return auth.AccessToken, nil
}

// UpdateAuthorization updates the authorization, or creates a new one.
func UpdateAuthorization(ctx context.Context, id, name, token, tokenType, scope, appID, botID string) error {
	now := util.Timestamp()
	var auth = AuthorizationDS{}
	key := authorizationKey(id)
	err := platform.DataStore().Get(ctx, key, &auth)

	if err == nil {
		auth.AccessToken = token
		auth.Scope = scope
		auth.Updated = now
	} else {
		auth = AuthorizationDS{
			ID:          id,
			Name:        name,
			AccessToken: token,
			TokenType:   tokenType,
			Scope:       scope,
			AppID:       appID,
			BotUserID:   botID,
			Created:     now,
			Updated:     now,
		}
	}

	// remove the entry from the cache if it is already there ...
	platform.Invalidate(ctx, cacheKey(id))

	_, err = platform.DataStore().Put(ctx, key, &auth)
	return err
}
