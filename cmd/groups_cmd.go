package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/nattvara/dfb/internal/groups"
	"github.com/spf13/cobra"
)

// groupsCmd represents the domains command
var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "Group commands",
	Long: `
A group contians a number of domains, and restic repositories to backup those domains to.
`,
}

var lsGroupsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List groups",
	Run: func(cmd *cobra.Command, args []string) {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)

		header := table.Row{"#", "Group"}
		t.AppendHeader(header)

		groups := groups.FetchGroups()
		var i int
		for _, group := range groups {
			i += 1
			row := table.Row{
				i,
				group.Name,
			}
			t.AppendRow(row)
		}

		t.Render()
	},
}

func init() {
	RootCmd.AddCommand(groupsCmd)

	groupsCmd.AddCommand(lsGroupsCmd)
}
