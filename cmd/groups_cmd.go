package cmd

import (
	"fmt"
	"log"
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

var addGroupsCmd = &cobra.Command{
	Use:   "add [group]",
	Short: "Add new group",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := groups.GetGroupFromString(args[0])
		if err == nil {
			log.Fatalf("group %s already exists", args[0])
		}

		fmt.Println(args[0])

		group := groups.New(args[0])
		group.Create()

		// domainName := args[1]

		// if group.DomainExists(domainName) {
		// 	fmt.Println("Domain already exists.")
		// 	return
		// }

		// if CreatSymlinkPath != "" {
		// 	group.AddDomainWithNameAndSymlink(domainName, CreatSymlinkPath)
		// } else {
		// 	group.AddDomainWithName(domainName)
		// }

		fmt.Printf("Group created at %s\n", group.Path)
	},
}

func init() {
	RootCmd.AddCommand(groupsCmd)

	groupsCmd.AddCommand(lsGroupsCmd)
	groupsCmd.AddCommand(addGroupsCmd)
}
