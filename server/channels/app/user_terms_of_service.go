// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"errors"
	"net/http"

	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/store"
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
)

func (a *App) GetUserTermsOfService(userID string) (*model.UserTermsOfService, *model.AppError) {
	u, err := a.Srv().Store().UserTermsOfService().GetByUser(userID)
	if err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return nil, model.NewAppError("GetUserTermsOfService", "app.user_terms_of_service.get_by_user.no_rows.app_error", nil, "", http.StatusNotFound).Wrap(err)
		default:
			return nil, model.NewAppError("GetUserTermsOfService", "app.user_terms_of_service.get_by_user.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		}
	}

	return u, nil
}

func (a *App) SaveUserTermsOfService(userID, termsOfServiceId string, accepted bool) *model.AppError {
	if accepted {
		userTermsOfService := &model.UserTermsOfService{
			UserId:           userID,
			TermsOfServiceId: termsOfServiceId,
		}

		if _, err := a.Srv().Store().UserTermsOfService().Save(userTermsOfService); err != nil {
			var appErr *model.AppError
			switch {
			case errors.As(err, &appErr):
				return appErr
			default:
				return model.NewAppError("SaveUserTermsOfService", "app.user_terms_of_service.save.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
			}
		}
	} else {
		if err := a.Srv().Store().UserTermsOfService().Delete(userID, termsOfServiceId); err != nil {
			return model.NewAppError("SaveUserTermsOfService", "app.user_terms_of_service.delete.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		}
	}

	return nil
}
