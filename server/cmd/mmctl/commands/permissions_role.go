// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package commands

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"

	"git.biggo.com/Funmula/BigGoChat/server/v8/cmd/mmctl/client"
	"git.biggo.com/Funmula/BigGoChat/server/v8/cmd/mmctl/printer"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
)

var RoleCmd = &cobra.Command{
	Use:   "role",
	Short: "Management of roles",
}

var ShowCmd = &cobra.Command{
	Use:     "show <role_name>",
	Short:   "Show the role information",
	Long:    "Show all the information about a role.",
	Example: `  permissions show system_user`,
	Args:    cobra.ExactArgs(1),
	RunE:    withClient(showRoleCmdF),
}

var AssignCmd = &cobra.Command{
	Use:   "assign <role_name> <username...>",
	Short: "Assign users to role (EE Only)",
	Long:  "Assign users to a role by username (Only works in Enterprise Edition).",
	Example: `  # Assign users with usernames 'john.doe' and 'jane.doe' to the role named 'system_admin'.
  permissions assign system_admin john.doe jane.doe

  # Examples using other system roles
  permissions assign system_manager john.doe jane.doe
  permissions assign system_user_manager john.doe jane.doe
  permissions assign system_read_only_admin john.doe jane.doe`,
	Args: cobra.MinimumNArgs(2),
	RunE: withClient(assignUsersCmdF),
}

var UnassignCmd = &cobra.Command{
	Use:   "unassign <role_name> <username...>",
	Short: "Unassign users from role (EE Only)",
	Long:  "Unassign users from a role by username (Only works in Enterprise Edition).",
	Example: `  # Unassign users with usernames 'john.doe' and 'jane.doe' from the role named 'system_admin'.
  permissions unassign system_admin john.doe jane.doe

  # Examples using other system roles
  permissions unassign system_manager john.doe jane.doe
  permissions unassign system_user_manager john.doe jane.doe
  permissions unassign system_read_only_admin john.doe jane.doe`,
	Args: cobra.MinimumNArgs(2),
	RunE: withClient(unassignUsersCmdF),
}

func init() {
	RoleCmd.AddCommand(
		AssignCmd,
		UnassignCmd,
		ShowCmd,
	)

	PermissionsCmd.AddCommand(
		RoleCmd,
	)
}

func prettyRole(role *model.Role) string {
	sort.Strings(role.Permissions)

	consolePermissionMap := map[string]bool{}
	for _, perm := range role.Permissions {
		if strings.HasPrefix(perm, "sysconsole_") {
			consolePermissionMap[perm] = true
		}
	}

	getUsedBy := func(permissionID string) []string {
		var usedByIDs []string
		if !strings.HasPrefix(permissionID, "sysconsole_") {
			usedBy := map[string]bool{} // map to make a unique set
			for key, vals := range model.SysconsoleAncillaryPermissions {
				for _, val := range vals {
					if val.Id == permissionID {
						if _, ok := consolePermissionMap[key]; ok {
							usedBy[key] = true
						}
					}
				}
			}
			for key := range usedBy {
				usedByIDs = append(usedByIDs, key)
			}
		}
		return usedByIDs
	}

	var b strings.Builder
	w := tabwriter.NewWriter(&b, 0, 0, 1, ' ', 0)

	// Only show the 3-column view if the role has sysconsole permissions
	// sysadmin has every permission, so no point in showing the "Used by"
	// column.
	if len(consolePermissionMap) > 0 && role.Name != "system_admin" {
		fmt.Fprintf(w, "\nProperty\tValue\tUsed by\n")
		fmt.Fprintf(w, "--------\t-----\t-------\n")
		fmt.Fprintf(w, "Name\t%s\t\n", role.Name)
		fmt.Fprintf(w, "DisplayName\t%s\t\n", role.DisplayName)
		fmt.Fprintf(w, "BuiltIn\t%v\t\n", role.BuiltIn)
		fmt.Fprintf(w, "SchemeManaged\t%v\t\n", role.SchemeManaged)
		for i, perm := range role.Permissions {
			if i == 0 {
				fmt.Fprintf(w, "Permissions\t%s\t%v\n", role.Permissions[0], strings.Join(getUsedBy(role.Permissions[0]), ", "))
			} else {
				fmt.Fprintf(w, "\t%s\t%v\n", perm, strings.Join(getUsedBy(perm), ", "))
			}
		}
	} else {
		fmt.Fprintf(w, "\nProperty\tValue\n")
		fmt.Fprintf(w, "--------\t-----\n")
		fmt.Fprintf(w, "Name\t%s\n", role.Name)
		fmt.Fprintf(w, "DisplayName\t%s\n", role.DisplayName)
		fmt.Fprintf(w, "BuiltIn\t%v\n", role.BuiltIn)
		fmt.Fprintf(w, "SchemeManaged\t%v\n", role.SchemeManaged)
		for i, perm := range role.Permissions {
			if i == 0 {
				fmt.Fprintf(w, "Permissions\t%s\n", role.Permissions[0])
			} else {
				fmt.Fprintf(w, "\t%s\n", perm)
			}
		}
	}

	w.Flush()

	return b.String()
}

type rolesSort []*model.Role

func (rl rolesSort)Len() int {
	return len(rl)
}

func (rl rolesSort)Swap(i, j int) {
	rl[i], rl[j] = rl[j], rl[i]
}

func (rl rolesSort) Less(i, j int) bool {
	scoreMap := map[string]int {
		"system": 0,
		"team": 1,
		"channel": 2,
	}
	getScore := func (name string) int {
		toks := strings.SplitN(name, "_", 2)
		if len(toks) == 1 {
			return 9999
		}
		score, ok := scoreMap[toks[0]]
		if ok {
			return score
		}else{
			return 9998
		}
	}
	si := getScore(rl[i].Name)
	sj := getScore(rl[j].Name)
	if si != sj {
		return si < sj
	}
	if rl[i].BuiltIn != rl[j].BuiltIn{
		return rl[i].BuiltIn
	}
	if rl[i].SchemeManaged != rl[j].SchemeManaged {
		return rl[i].SchemeManaged
	}
	return rl[i].Name < rl[j].Name
}

func showRoleCmdF(c client.Client, cmd *cobra.Command, args []string) error {
	role, _, err := c.GetRoleByName(context.TODO(), args[0])
	if err != nil {
		return err
	}

	printer.PrintT(prettyRole(role), nil)

	return nil
}

func showRolesCmdF(c client.Client, cmd *cobra.Command, args []string) error {
	roles, _, err := c.GetAllRoles(context.TODO())
	if err != nil {
		return err
	}
	sort.Sort(rolesSort(roles))
	rolePermMap := make(map[string](map[string]bool), len(roles))

	line := "Scope,Permission"
	for _, role := range roles {
		line += fmt.Sprintf(",%s", role.Name)
		permSet := make(map[string]bool, len(role.Permissions))
		for _, perm := range role.Permissions {
			permSet[perm] = true
		}
		rolePermMap[role.Name] = permSet
	}
	printer.Print(line)
	line = "SchemeManaged,"
	for _, role := range roles {
		line += fmt.Sprintf(",%v", role.SchemeManaged)
	}
	printer.Print(line)
	line = "BuiltIn,"
	for _, role := range roles {
		line += fmt.Sprintf(",%v", role.BuiltIn)
	}
	printer.Print(line)
	for _, permission := range model.AllPermissions {
		line := fmt.Sprintf("%s,%s", permission.Scope, permission.Id)
		for _, role := range roles {
			permSet := rolePermMap[role.Name]
			if e := permSet[permission.Id]; e {
				line += ",1"
			} else {
				line += ",0"
			}
		}
		printer.Print(line)
	}
	return nil
}

func assignUsersCmdF(c client.Client, cmd *cobra.Command, args []string) error {
	role, _, err := c.GetRoleByName(context.TODO(), args[0])
	if err != nil {
		return err
	}

	users := getUsersFromUserArgs(c, args[1:])

	var errs *multierror.Error
	for i, user := range users {
		if user == nil {
			printer.PrintError("Couldn't find user '" + args[i+1] + "'.")
			errs = multierror.Append(errs, fmt.Errorf("couldn't find user '%s'", args[i+1]))
			continue
		}

		var userHasRequestedRole bool
		startingRoles := strings.Fields(user.Roles)
		for _, roleName := range startingRoles {
			if roleName == role.Name {
				userHasRequestedRole = true
			}
		}

		if userHasRequestedRole {
			continue
		}

		userRoles := startingRoles
		userRoles = append(userRoles, role.Name)
		_, err = c.UpdateUserRoles(context.TODO(), user.Id, strings.Join(userRoles, " "))
		if err != nil {
			return err
		}
	}

	return errs.ErrorOrNil()
}

func unassignUsersCmdF(c client.Client, cmd *cobra.Command, args []string) error {
	users := getUsersFromUserArgs(c, args[1:])

	for i, user := range users {
		if user == nil {
			printer.PrintError("Couldn't find user '" + args[i+1] + "'.")
			continue
		}

		userRoles := strings.Fields(user.Roles)
		originalCount := len(userRoles)

		for j := 0; j < len(userRoles); j++ {
			if userRoles[j] == args[0] {
				userRoles = append(userRoles[:j], userRoles[j+1:]...)
				j--
			}
		}

		if originalCount > len(userRoles) {
			_, err := c.UpdateUserRoles(context.TODO(), user.Id, strings.Join(userRoles, " "))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
