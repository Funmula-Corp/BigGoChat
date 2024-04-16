package oauthbiggo

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/public/shared/request"
	"github.com/mattermost/mattermost/server/v8/einterfaces"
)

const UserAuthServiceBiggo = "biggo"

type BiggoProvider struct {
}

type BiggoUser struct {
	Id     string `json:"userid"`
	Email  string `json:"email"`
	UserId string `json:"at_userid"`
	Name   string `json:"name"`
	Image  string `json:"origin_profileimg"`
}

func init() {
	provider := &BiggoProvider{}
	einterfaces.RegisterOAuthProvider(UserAuthServiceBiggo, provider)
}

func userFromBiggoUser(logger mlog.LoggerIFace, glu *BiggoUser) *model.User {
	user := &model.User{}
	username := glu.UserId
	user.Username = model.CleanUsername(logger, username)
	splitName := strings.Split(glu.Name, " ")
	if len(splitName) == 2 {
		user.FirstName = splitName[0]
		user.LastName = splitName[1]
	} else {
		user.FirstName = glu.Name
	}
	user.Email = glu.Email
	user.Email = strings.ToLower(user.Email)
	user.AuthData = &glu.Id
	user.AuthService = UserAuthServiceBiggo

	return user
}

func BiggoUserFromJSON(data io.Reader) (*BiggoUser, error) {
	decoder := json.NewDecoder(data)
	var glu BiggoUser
	err := decoder.Decode(&glu)
	if err != nil {
		return nil, err
	}
	return &glu, nil
}

func (glu *BiggoUser) IsValid() error {
	if glu.Id == "" {
		return errors.New("user id can't be empty")
	}

	if glu.Email == "" {
		return errors.New("user e-mail should not be empty")
	}

	return nil
}

func (gp *BiggoProvider) GetUserFromJSON(c request.CTX, data io.Reader, tokenUser *model.User) (*model.User, error) {
	glu, err := BiggoUserFromJSON(data)
	if err != nil {
		return nil, err
	}
	if err = glu.IsValid(); err != nil {
		return nil, err
	}

	return userFromBiggoUser(c.Logger(), glu), nil
}

func (gp *BiggoProvider) GetSSOSettings(_ request.CTX, config *model.Config, service string) (*model.SSOSettings, error) {
	return &config.BiggoSettings, nil
}

func (gp *BiggoProvider) GetUserFromIdToken(_ request.CTX, idToken string) (*model.User, error) {
	return nil, nil
}

func (gp *BiggoProvider) IsSameUser(_ request.CTX, dbUser, oauthUser *model.User) bool {
	return dbUser.AuthData == oauthUser.AuthData
}
