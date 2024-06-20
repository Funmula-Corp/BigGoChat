package app

import (
	"database/sql"
	"errors"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
)

// should define in `model` package
const (
	ChannelReadOnlyRoleId = "biggoryyyyyyyyyyyyyyyyyyyb"
	ChannelReadOnlyRoleName = "channel_readonly"

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
		Id: ChannelReadOnlySchemeId,
		Name: "announcement",
		DisplayName: "announcement",
		Scope: model.SchemeScopeChannel,
		DefaultChannelAdminRole: model.ChannelAdminRoleId,
		DefaultChannelUserRole: ChannelReadOnlyRoleName,
		DefaultChannelVerifiedRole: model.ChannelVerifiedRoleId,
		DefaultChannelGuestRole: model.ChannelGuestRoleId,
	}

	if _, err := s.Store().Scheme().CreateScheme(scheme); err != nil {
		mlog.Fatal("Failed to migrate scheme to database.", mlog.Err(err))
		return
	}

	system := model.System{
		Name: model.CustomChannelReadOnlyRoleCreationMigrationKey,
		Value: "true",
	}
	if err := s.Store().System().Save(&system); err != nil {
		mlog.Fatal("Failed to create channel read-only role migration as completed.", mlog.Err(err))
	}
}

const (
	SystemVerifiedRoleSpecialId = "biggoyyyyyyyyyyyyyyyyyyyyn"
)

func (s *Server) doSystemVerifiedRoleCreationMigration(c *request.Context) {
	if _, err := s.Store().System().GetByName(model.CustomSystemVerifiedRoleCreationMigrationKey); err == nil {
		return
	}

	userRole, err := s.Store().Role().GetByName(c.Context(), model.SystemUserRoleId);
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
	TeamVerifiedRoleId = "biggoryyyyyyyyyyyyyyyyyyyd"
	ChannelVerifiedRoleId = "biggoryyyyyyyyyyyyyyyyyyyr"
)

func (s *Server) doVerifiedTierMigration(c *request.Context) {
	if _, err := s.Store().System().GetByName(model.CustomVerifiedTierMigrationMigrationKey); err == nil {
		return
	}

	teamUserRole, err := s.Store().Role().GetByName(c.Context(), model.TeamUserRoleId);
	if err != nil {
		mlog.Fatal("failed to get role by name", mlog.Err(err))
		return
	}

	// inherit from system_user
	teamVerifiedRole := &model.Role{
		Id: TeamVerifiedRoleId,
		Name: model.TeamVerifiedRoleId,
		DisplayName: "authentication.roles.team_verified.name",
		Description: "authentication.roles.team_verified.description",
		Permissions: teamUserRole.Permissions,
		SchemeManaged: true,
		BuiltIn: true,
	}

	if _, err := s.Store().Role().CreateRole(teamVerifiedRole); err != nil {
		mlog.Fatal("Failed to migrate role to database.", mlog.Err(err))
		return
	}

	channelUserRole, err := s.Store().Role().GetByName(c.Context(), model.ChannelUserRoleId);
	if err != nil {
		mlog.Fatal("failed to get role by name", mlog.Err(err))
		return
	}
	channelVerifiedRole := &model.Role{
		Id: ChannelVerifiedRoleId,
		Name: model.ChannelVerifiedRoleId,
		DisplayName: "authentication.roles.channel_verified.name",
		Description: "authentication.roles.channel_verified.description",
		Permissions: channelUserRole.Permissions,
		SchemeManaged: true,
		BuiltIn: true,
	}

	if _, err := s.Store().Role().CreateRole(channelVerifiedRole); err != nil {
		mlog.Fatal("Failed to migrate role to database.", mlog.Err(err))
		return
	}

	channelUserRole.Permissions = []string{
		"add_reaction", "edit_post", "read_channel", "read_channel_content", "read_private_channel_groups", "read_public_channel_groups", "use_channel_mentions", "use_group_mentions"}
	if _, err := s.Store().Role().Save(channelUserRole); err != nil {
		mlog.Fatal("Failed to migrate role to database.", mlog.Err(err))
		return
	}

	scopes := []string {
		model.SchemeScopeTeam     ,
		model.SchemeScopeChannel  ,
		model.SchemeScopePlaybook ,
		model.SchemeScopeRun      ,
	}
	pageSize := 100
	for _, scope := range(scopes) {
		offset := 0
		for {
			mlog.Info("migrate scheme", mlog.String("scope", scope))
			schemes, err := s.Store().Scheme().GetAllPage(scope, offset, pageSize)
			if errors.Is(err, sql.ErrNoRows){
				break
			}else if err != nil {
				mlog.Fatal("Failed to migrate scheme", mlog.Err(err))
				return
			}
			for _, scheme := range(schemes){
				if scope == model.SchemeScopeChannel || scope == model.SchemeScopeTeam {
					scheme.DefaultChannelVerifiedRole = channelVerifiedRole.Name
				}
				if scope == model.SchemeScopeTeam {
					scheme.DefaultTeamVerifiedRole = teamVerifiedRole.Name
				}
				_, err = s.Store().Scheme().Save(scheme)
				if err != nil {
					mlog.Fatal("Failed to migrate scheme", mlog.Err(err))
					return
				}
			}
			offset += pageSize
			if len(schemes) < pageSize {
				break
			}
		}
	}

	system := model.System{
		Name: model.CustomVerifiedTierMigrationMigrationKey,
		Value: "true",
	}
	if err := s.Store().System().Save(&system); err != nil {
		mlog.Fatal("Failed to create channel read-only role migration as completed.", mlog.Err(err))
	}
}

func (s *Server) doBiggoMigration(c *request.Context) {
	s.doChannelReadOnlyRoleCreationMigration()
	s.doSystemVerifiedRoleCreationMigration(c)
	s.doVerifiedTierMigration(c)
}
