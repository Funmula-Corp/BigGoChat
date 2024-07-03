package app

import (
	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
)

// should define in `model` package
const (
	ChannelReadOnlyRoleId = "biggoryyyyyyyyyyyyyyyyyyyb"
	ChannelReadOnlyRoleName = model.ChannelReadOnlyRoleId

	ChannelReadOnlySchemeId = "biggosyyyyyyyyyyyyyyyyyyyd"
)

func (s *Server) doChannelReadOnlyRoleCreationMigration() {
	if _, err := s.Store().System().GetByName(model.CustomChannelReadOnlyRoleCreationMigrationKey); err == nil {
		return
	}

	role := &model.Role{
		Id: ChannelReadOnlyRoleId,
		Name: ChannelReadOnlyRoleName,
		DisplayName: "authentication.roles.channel_readonly.name",
		Description: "authentication.roles.channel_readonly.description",
		Permissions: []string{
			model.PermissionReadChannel.Id,
			model.PermissionReadChannelContent.Id,
		},
		SchemeManaged: false,
		BuiltIn: true,
	}

	if _, err := s.Store().Role().CreateRole(role); err != nil {
		mlog.Fatal("Failed to migrate role to database.", mlog.Err(err))
		return
	}

	scheme := &model.Scheme{
		Id:                         ChannelReadOnlySchemeId,
		Name:                       model.ChannelReadOnlySchemeName,
		DisplayName:                model.ChannelReadOnlySchemeName,
		Scope:                      model.SchemeScopeChannel,
		DefaultChannelAdminRole:    model.ChannelAdminRoleId,
		DefaultChannelVerifiedRole: ChannelReadOnlyRoleName,
		DefaultChannelUserRole:     ChannelReadOnlyRoleName,
		DefaultChannelGuestRole:    model.ChannelGuestRoleId,
	}

	if _, err := s.Store().Scheme().CreateBuiltInScheme(scheme); err != nil {
		mlog.Fatal("Failed to migrate scheme to database.", mlog.Err(err))
		return
	}

	system := model.System{
		Name:  model.CustomChannelReadOnlyRoleCreationMigrationKey,
		Value: "true",
	}
	if err := s.Store().System().Save(&system); err != nil {
		mlog.Fatal("Failed to create channel read-only role migration as completed.", mlog.Err(err))
	}
}

const (
	SystemVerifiedRoleId   = "biggoyyyyyyyyyyyyyyyyyyyyn"
	SystemVerifiedRoleName =  model.SystemVerifiedRoleId
	SystemVerifiedRoleSpecialId = "biggoyyyyyyyyyyyyyyyyyyyyn"
)

func (s *Server) doSystemVerifiedRoleCreationMigration(c *request.Context) {
	if _, err := s.Store().System().GetByName(model.CustomSystemVerifiedRoleCreationMigrationKey); err == nil {
		return
	}

	userRole, err := s.Store().Role().GetByName(c.Context(), model.SystemUserRoleId)
	if err != nil {
		mlog.Fatal("failed to get role by name", mlog.Err(err))
		return
	}

	// inherit from system_user
	permissions := userRole.Permissions

	role := &model.Role{
		Id: SystemVerifiedRoleSpecialId,
		Name: model.SystemVerifiedRoleId,
		DisplayName: "authentication.roles.system_verified.name",
		Description: "authentication.roles.system_verified.description",
		Permissions: permissions,
		SchemeManaged: true,
		BuiltIn: true,
	}

	if _, err := s.Store().Role().CreateRole(role); err != nil {
		mlog.Fatal("Failed to migrate role to database.", mlog.Err(err))
		return
	}

	system := model.System{
		Name: model.CustomSystemVerifiedRoleCreationMigrationKey,
		Value: "true",
	}
	if err := s.Store().System().Save(&system); err != nil {
		mlog.Fatal("Failed to create channel read-only role migration as completed.", mlog.Err(err))
	}
}

const (
	ChannelVerifiedRoleId  = "biggoryyyyyyyyyyyyyyyyyyyr"
	ChannelVerifiedRoleName = "channel_verified"
	TeamVerifiedRoleId  = "biggoryyyyyyyyyyyyyyyyyyyd"
	TeamVerifiedRoleName = "team_verified"
	TeamModeratorRoleId  = "biggoryyyyyyyyyyyyyyyyyyyf"
	TeamModeratorRoleName = "team_moderator"

	ChannelAllowUnverifiedSchemeId = "biggosyyyyyyyyyyyyyyyyyyyf"
	ChannelAllowUnverifiedSchemeName = model.ChannelAllowUnverifiedSchemeName
)
func (s *Server) doMigrationKeySchemesRolesCreation(c *request.Context) {

	if _, err := s.Store().System().GetByName(model.MigrationBigGoSchemeRolesCreation); err == nil {
		return
	}

	channelUserRole, err := s.Store().Role().GetByName(c.Context(), model.ChannelUserRoleId)
	if err != nil {
		mlog.Fatal("failed to get role by name", mlog.Err(err))
		return
	}
	// inherit from channel_user
	permissions := channelUserRole.Permissions
	channelVerifiedRole := &model.Role{
		Id:            ChannelVerifiedRoleId,
		Name:          model.ChannelVerifiedRoleId,
		DisplayName:   "authentication.roles.channel_verified.name",
		Description:   "authentication.roles.channel_verified.description",
		Permissions:   permissions,
		SchemeManaged: true,
		BuiltIn:       true,
	}
	if _, err := s.Store().Role().CreateRole(channelVerifiedRole); err != nil {
		mlog.Fatal("Failed to create role to database.", mlog.Err(err))
		return
	}

	teamUserRole, err := s.Store().Role().GetByName(c.Context(), model.TeamUserRoleId)
	if err != nil {
		mlog.Fatal("failed to get role by name", mlog.Err(err))
		return
	}
	// inherit from team_user
	permissions = teamUserRole.Permissions
	teamVerifiedRole := &model.Role{
		Id:            TeamVerifiedRoleId,
		Name:          TeamVerifiedRoleName,
		DisplayName:   "authentication.roles.team_verified.name",
		Description:   "authentication.roles.team_verified.description",
		Permissions:   permissions,
		SchemeManaged: true,
		BuiltIn:       true,
	}
	if _, err := s.Store().Role().CreateRole(teamVerifiedRole); err != nil {
		mlog.Fatal("Failed to migrate role to database.", mlog.Err(err))
		return
	}

	teamAdminRole, err := s.Store().Role().GetByName(c.Context(), model.TeamAdminRoleId)
	if err != nil {
		mlog.Fatal("failed to get role by name", mlog.Err(err))
		return
	}
	// inherit from team_admin
	permissions = []string{}
	for _, p := range(teamAdminRole.Permissions){
		if p != model.PermissionManageTeamRoles.Id {
			permissions = append(permissions, p)
		}
	}
	teamModeratorRole := &model.Role{
		Id:            TeamModeratorRoleId,
		Name:          TeamModeratorRoleName,
		DisplayName:   "authentication.roles.team_moderator.name",
		Description:   "authentication.roles.team_moderator.description",
		Permissions:   permissions,
		SchemeManaged: true,
		BuiltIn:       true,
	}
	if _, err := s.Store().Role().CreateRole(teamModeratorRole); err != nil {
		mlog.Fatal("Failed to migrate role to database.", mlog.Err(err))
		panic("failed")
	}

	// migrate schemes
	offset := 0
	pageSize := 100
	for _, scope := range []string{model.SchemeScopeTeam, model.SchemeScopeChannel} {
		for {
			schemes, err := s.Store().Scheme().GetAllPage(scope, offset, pageSize)
			if err != nil {
				mlog.Fatal("Failed to get schemes", mlog.String("scope", scope), mlog.Err(err))
				return
			}
			for _, scheme := range schemes {
				if scheme.Id == ChannelReadOnlySchemeId {
					scheme.DefaultChannelVerifiedRole = ChannelReadOnlyRoleName
				} else {
					if scheme.Scope == model.SchemeScopeTeam {
						scheme.DefaultTeamModeratorRole = teamModeratorRole.Id
						scheme.DefaultTeamVerifiedRole = teamVerifiedRole.Id
					}
					scheme.DefaultChannelVerifiedRole = channelVerifiedRole.Id
				}
				if _, err := s.Store().Scheme().Save(scheme); err != nil {
					mlog.Fatal("Failed to save schemes", mlog.String("scope", scope), mlog.String("scheme", scheme.Id), mlog.Err(err))
				}
			}
			if len(schemes) < pageSize {
				break
			}else{
				offset += pageSize
			}
		}
	}

	scheme := &model.Scheme{
		Id:                         ChannelAllowUnverifiedSchemeId,
		Name:                       ChannelAllowUnverifiedSchemeName,
		DisplayName:                ChannelAllowUnverifiedSchemeName,
		Scope:                      model.SchemeScopeChannel,
		DefaultChannelAdminRole:    model.ChannelAdminRoleId,
		DefaultChannelVerifiedRole: model.ChannelVerifiedRoleId,
		DefaultChannelUserRole:     model.ChannelVerifiedRoleId,
		DefaultChannelGuestRole:    model.ChannelGuestRoleId,
	}

	if _, err := s.Store().Scheme().CreateBuiltInScheme(scheme); err != nil {
		mlog.Fatal("Failed to create scheme to database.", mlog.Err(err))
		return
	}

	// migrate TeamMembers and ChannelMembers
	users, err := s.Store().User().GetAll()
	if err != nil {
		mlog.Fatal("Failed to get user", mlog.Err(err))
		return
	}
	for _, user := range(users){
		if s.Store().User().UpdateMemberVerifiedStatus(c, user) != nil {
			mlog.Fatal("Failed to update MemberVerifiedStatus", mlog.String("userId", user.Id))
			panic("")
		}
	}

	// PluginAPI 要有可以更新 TeamMember 和 ChannelMember 的地方
	// TODO: a scheme for channel that unverified user can post

	system := model.System{
		Name:  model.MigrationBigGoSchemeRolesCreation,
		Value: "true",
	}
	if err := s.Store().System().Save(&system); err != nil {
		mlog.Fatal("Failed to create verified-tier roles migration as completed.", mlog.Err(err))
	}
}

func (a *App) doMigrationKeyBigGoRolesPermissions() (permissionsMap, error) {
	return permissionsMap{
		permissionTransformation{
			On:     permissionAnd(isRole(model.ChannelUserRoleId)),
			Remove: []string{
				PermissionManagePublicChannelMembers,
				PermissionManagePrivateChannelMembers,
				PermissionManagePublicChannelProperties,
				PermissionManagePrivateChannelProperties,
				PermissionDeletePublicChannel,
				PermissionDeletePrivateChannel,
				PermissionCreatePost,
				PermissionAddReaction,
				model.PermissionCreatePostEphemeral.Id,
				model.PermissionUploadFile.Id,
				// PermissionUseChannelMentions,
			},
		},
		permissionTransformation{
			On:     permissionAnd(isRole(model.ChannelVerifiedRoleId)),
			Remove: []string{
				PermissionManagePublicChannelMembers,
				PermissionManagePrivateChannelMembers,
				PermissionManagePublicChannelProperties,
				PermissionManagePrivateChannelProperties,
			},
		},
		permissionTransformation{
			On:     permissionAnd(isRole(model.ChannelAdminRoleId)),
			Remove: []string{model.PermissionConvertPrivateChannelToPublic.Id},
			Add: []string{
				PermissionManagePublicChannelMembers,
				PermissionManagePrivateChannelMembers,
				PermissionConvertPublicChannelToPrivate,
				PermissionManagePublicChannelProperties,
				PermissionManagePrivateChannelProperties,
			},
		},
		permissionTransformation{
			On:     permissionAnd(isRole(model.TeamUserRoleId)),
			Remove: []string{
				model.PermissionCreatePublicChannel.Id,
				model.PermissionCreatePrivateChannel.Id,
				model.PermissionPrivatePlaybookCreate.Id},
		},
		permissionTransformation{
			On:     permissionAnd(isRole(model.TeamModeratorRoleId)),
			Add: []string{
				PermissionManagePublicChannelProperties,
				PermissionManagePrivateChannelProperties,
			},
			Remove: []string{
				model.PermissionConvertPrivateChannelToPublic.Id,
			},
		},
		permissionTransformation{
			On:     permissionAnd(isRole(model.TeamAdminRoleId)),
			Add: []string{
				PermissionManagePublicChannelProperties,
				PermissionManagePrivateChannelProperties,
			},
			Remove: []string{
				model.PermissionConvertPrivateChannelToPublic.Id,
			},
		},
		permissionTransformation{
			On:     permissionAnd(isRole(model.SystemUserRoleId)),
			Remove: []string{
				model.PermissionCreateTeam.Id,
				model.PermissionCreateGroupChannel.Id,
			},
		},
		permissionTransformation{
			On:     permissionAnd(isRole(model.SystemVerifiedRoleId)),
			Add: []string{model.PermissionCreateTeam.Id},
		},
	}, nil
}

func (s *Server) doBiggoPermissionMigration() error {
	a := New(ServerConnector(s.Channels()))
	PermissionsMigrations := []struct {
		Key       string
		Migration func() (permissionsMap, error)
	}{
		{Key: model.MigrationKeyBigGoRolesPermissions, Migration: a.doMigrationKeyBigGoRolesPermissions},
	}

	roles, err := s.Store().Role().GetAll()
	if err != nil {
		return err
	}

	for _, migration := range PermissionsMigrations {
		migMap, err := migration.Migration()
		if err != nil {
			return err
		}
		if err := s.doPermissionsMigration(migration.Key, migMap, roles); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) doBiggoMigration(c *request.Context) {
	s.doChannelReadOnlyRoleCreationMigration()
	s.doSystemVerifiedRoleCreationMigration(c)
	s.doMigrationKeySchemesRolesCreation(c)

	// must be the last, make sure all roles are created
	if err := s.doBiggoPermissionMigration(); err != nil {
		mlog.Fatal("(app.App).doBiggoPermissionMigration failed", mlog.Err(err))
	}
}
