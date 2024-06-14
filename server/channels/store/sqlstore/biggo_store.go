package sqlstore

import (
	"fmt"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/store"
	"github.com/pkg/errors"
)

func (s *SqlRoleStore) CreateRole(role *model.Role) (_ *model.Role, err error) {
	transaction, terr := s.GetMasterX().Beginx()
	if terr != nil {
		return nil, errors.Wrap(terr, "begin_transaction")
	}
	defer finalizeTransactionX(transaction, &terr)

	createdRole, terr := s.createRole(role, transaction)
	if terr != nil {
		return nil, errors.Wrap(terr, "unable to create Role")
	} else if terr = transaction.Commit(); terr != nil {
		return nil, errors.Wrap(terr, "commit_transaction")
	}
	return createdRole, nil
}

func (s *SqlSchemeStore) CreateScheme(scheme *model.Scheme) (_ *model.Scheme, err error) {
	transaction, terr := s.GetMasterX().Beginx()
	if terr != nil {
		return nil, errors.Wrap(terr, "begin_transaction")
	}
	defer finalizeTransactionX(transaction, &terr)

	newScheme, terr := s.createSchemeWithoutCreateRoles(scheme, transaction)
	if terr != nil {
		return nil, terr
	}
	if terr = transaction.Commit(); terr != nil {
		return nil, errors.Wrap(terr, "commit_transaction")
	}
	return newScheme, nil
}

func (s *SqlSchemeStore) createSchemeWithoutCreateRoles(scheme *model.Scheme, transaction *sqlxTxWrapper) (*model.Scheme, error) {
	// fetch all exists roles in scheme
	schemeRoleNames := []string{}
	if scheme.DefaultTeamAdminRole != "" {
		schemeRoleNames = append(schemeRoleNames, scheme.DefaultTeamAdminRole)
	}
	if scheme.DefaultTeamUserRole != "" {
		schemeRoleNames = append(schemeRoleNames, scheme.DefaultTeamUserRole)
	}
	if scheme.DefaultChannelAdminRole != "" {
		schemeRoleNames = append(schemeRoleNames, scheme.DefaultChannelAdminRole)
	}
	if scheme.DefaultChannelUserRole != "" {
		schemeRoleNames = append(schemeRoleNames, scheme.DefaultChannelUserRole)
	}
	if scheme.DefaultTeamGuestRole != "" {
		schemeRoleNames = append(schemeRoleNames, scheme.DefaultTeamGuestRole)
	}
	if scheme.DefaultChannelGuestRole != "" {
		schemeRoleNames = append(schemeRoleNames, scheme.DefaultChannelGuestRole)
	}
	if scheme.DefaultPlaybookAdminRole != "" {
		schemeRoleNames = append(schemeRoleNames, scheme.DefaultPlaybookAdminRole)
	}
	if scheme.DefaultPlaybookMemberRole != "" {
		schemeRoleNames = append(schemeRoleNames, scheme.DefaultPlaybookMemberRole)
	}
	if scheme.DefaultRunAdminRole != "" {
		schemeRoleNames = append(schemeRoleNames, scheme.DefaultRunAdminRole)
	}
	if scheme.DefaultRunMemberRole != "" {
		schemeRoleNames = append(schemeRoleNames, scheme.DefaultRunMemberRole)
	}


	schemeRoles := make(map[string]*model.Role)
	roles, err := s.SqlStore.Role().GetByNames(schemeRoleNames)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		schemeRoles[role.Name] = role
	}

	// check all roles are exist
	if len(schemeRoles) != len(schemeRoleNames) {
		return nil, errors.New("createScheme: unable to retrieve default scheme roles")
	}

	if scheme.Id == "" {
		scheme.Id = model.NewId()
	}
	if scheme.Name == "" {
		scheme.Name = model.NewId()
	}
	scheme.CreateAt = model.GetMillis()
	scheme.UpdateAt = scheme.CreateAt

	// Validate the scheme
	if !scheme.IsValidForCreate() {
		return nil, store.NewErrInvalidInput("Scheme", "<any>", fmt.Sprintf("%v", scheme))
	}

	if _, err := transaction.NamedExec(`INSERT INTO Schemes
	(Id, Name, DisplayName, Description, Scope, DefaultTeamAdminRole, DefaultTeamUserRole, DefaultTeamGuestRole, DefaultChannelAdminRole, DefaultChannelUserRole, DefaultChannelGuestRole, CreateAt, UpdateAt, DeleteAt, DefaultPlaybookAdminRole, DefaultPlaybookMemberRole, DefaultRunAdminRole, DefaultRunMemberRole)
		VALUES
		(:Id, :Name, :DisplayName, :Description, :Scope, :DefaultTeamAdminRole, :DefaultTeamUserRole, :DefaultTeamGuestRole, :DefaultChannelAdminRole, :DefaultChannelUserRole, :DefaultChannelGuestRole, :CreateAt, :UpdateAt, :DeleteAt, :DefaultPlaybookAdminRole, :DefaultPlaybookMemberRole, :DefaultRunAdminRole, :DefaultRunMemberRole)`, scheme); err != nil {
		return nil, errors.Wrap(err, "failed to save Scheme")
	}

	return scheme, nil
}
