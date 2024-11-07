package commands

import (
	"context"

	"git.biggo.com/Funmula/BigGoChat/server/v8/cmd/mmctl/client"
	"git.biggo.com/Funmula/BigGoChat/server/v8/cmd/mmctl/printer"
	"github.com/spf13/cobra"
)

var CacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "cache observation",
}

var CacheUserCmd = &cobra.Command{
	Use:   "user",
	Short: "cache observation for user",
}

var CacheUserSessionsCmd = &cobra.Command{
	Use:   "sessions [user]",
	Short: "cache observation for user sessions",
	Args:  cobra.ExactArgs(1),
	RunE:  withClient(cacheUserSessionsCmdF),
}

var CacheUserChannelMembersCmd = &cobra.Command{
	Use:   "channel-members [user]",
	Short: "cache observation for user channel members",
	Args:  cobra.ExactArgs(1),
	RunE:  withClient(cacheUserChannelMembersCmdF),
}

func init() {
	CacheUserChannelMembersCmd.Flags().BoolP("include-deleted", "d", false, "include deleted channel members")

	CacheUserCmd.AddCommand(
		CacheUserSessionsCmd,
		CacheUserChannelMembersCmd,
	)

	CacheCmd.AddCommand(CacheUserCmd)

	RootCmd.AddCommand(CacheCmd)
}

func cacheUserSessionsCmdF(c client.Client, cmd *cobra.Command, args []string) error {
	user, err := getUserFromArg(c, args[0])
	if err != nil {
		return err
	}

	sessions, _, err := c.GetCachedSessions(context.TODO(), user.Id)
	if err != nil {
		return err
	}

	if len(sessions) == 0 {
		printer.Print("No cached sessions found")
		return nil
	}

	printer.PrintT("{{len .}} sessions cached for user {{(index . 0).UserId}}", sessions)

	for _, session := range sessions {
		printer.PrintT("  {{.Id}} - {{.DeviceId}} - {{.CreateAt}} - {{.ExpiresAt}}", session)
	}

	return nil
}

func cacheUserChannelMembersCmdF(c client.Client, cmd *cobra.Command, args []string) error {
	user, err := getUserFromArg(c, args[0])
	if err != nil {
		return err
	}

	includeDeleted, err := cmd.Flags().GetBool("include-deleted")
	if err != nil {
		return err
	}

	members, _, err := c.GetCachedAllChannelMembersForUser(context.TODO(), user.Id, includeDeleted)
	if err != nil {
		return err
	}

	if len(members) == 0 {
		printer.Print("No cached channel members found")
		return nil
	}

	printer.PrintT("{{len .}} channel memberships cached for user {{.UserId}}", map[string]interface{}{
		"len":    len(members),
		"UserId": user.Id,
	})

	for channelId, member := range members {
		printer.PrintT("  {{.ChannelId}} - Roles: {{.Roles}} - ExcludePermissions: {{.ExcludePermissions}} - IgnoreExclude: {{.IgnoreExclude}}", map[string]interface{}{
			"ChannelId":          channelId,
			"Roles":              member.Roles,
			"ExcludePermissions": member.ExcludePermissions,
			"IgnoreExclude":      member.IgnoreExclude,
		})
	}

	return nil
}
