package app

import (
	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
)

const (
	CustomChannelReadOnlyRoleCreationMigrationKey = "CustomChannelReadOnlyRoleCreationMigrationComplete"
	CustomSystemVerifiedRoleCreationMigrationKey = "CustomSystemVerifiedRoleCreationMigrationComplete"
)

// should define in `model` package
const (
	ChannelReadOnlyRoleId = "biggoryyyyyyyyyyyyyyyyyyyb"
	ChannelReadOnlyRoleName = "channel_readonly"

	ChannelReadOnlySchemeId = "biggosyyyyyyyyyyyyyyyyyyyd"
)

func (s *Server) doChannelReadOnlyRoleCreationMigration() {
	if _, err := s.Store().System().GetByName(CustomChannelReadOnlyRoleCreationMigrationKey); err == nil {
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
		SchemeManaged: true,
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
		DefaultChannelGuestRole: model.ChannelGuestRoleId,
	}

	if _, err := s.Store().Scheme().CreateScheme(scheme); err != nil {
		mlog.Fatal("Failed to migrate scheme to database.", mlog.Err(err))
		return
	}

	system := model.System{
		Name: CustomChannelReadOnlyRoleCreationMigrationKey,
		Value: "true",
	}
	if err := s.Store().System().Save(&system); err != nil {
		mlog.Fatal("Failed to create channel read-only role migration as completed.", mlog.Err(err))
	}
}

const (
	SystemVerifiedRoleId = "biggoyyyyyyyyyyyyyyyyyyyyn"
	SystemVerifiedRoleName = "system_verified"
)

func (s *Server) doSystemVerifiedRoleCreationMigration(c *request.Context) {
	if _, err := s.Store().System().GetByName(CustomSystemVerifiedRoleCreationMigrationKey); err == nil {
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
		Id: SystemVerifiedRoleId,
		Name: SystemVerifiedRoleName,
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
		Name: CustomSystemVerifiedRoleCreationMigrationKey,
		Value: "true",
	}
	if err := s.Store().System().Save(&system); err != nil {
		mlog.Fatal("Failed to create channel read-only role migration as completed.", mlog.Err(err))
	}
}

func (s *Server) doBiggoMigration(c *request.Context) {
	s.doChannelReadOnlyRoleCreationMigration()
	s.doSystemVerifiedRoleCreationMigration(c)
}