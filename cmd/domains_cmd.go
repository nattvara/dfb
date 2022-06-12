package cmd

import (
	"os"

	"github.com/nattvara/dfb/internal/groups"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// domainsCmd represents the domains command
var domainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "Domain commands",
}

var ListIncludeRepositories bool
var ListIncludeSymlink bool

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List configured domains",
	Run: func(cmd *cobra.Command, args []string) {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)

		header := table.Row{"#", "Group", "Domain"}
		if ListIncludeRepositories {
			header = append(header, "Repositories")
		}
		if ListIncludeSymlink {
			header = append(header, "Symlink")
		}
		t.AppendHeader(header)

		groups := groups.FetchGroups()
		for _, group := range groups {
			domains := group.Domains()

			var i int
			for _, domain := range domains {
				i += 1

				var symlink string
				if domain.Symlink != nil {
					symlink = domain.Symlink.Source
				}

				row := table.Row{
					i,
					group.Name,
					domain.Name,
				}
				if ListIncludeRepositories {
					row = append(row, domain.Repositories)
				}
				if ListIncludeSymlink {
					row = append(row, symlink)
				}
				t.AppendRow(row)
			}
		}

		t.Render()
	},
}

func init() {
	RootCmd.AddCommand(domainsCmd)

	lsCmd.Flags().BoolVarP(&ListIncludeRepositories, "include-repositories", "r", false, "include repositories in output")
	lsCmd.Flags().BoolVarP(&ListIncludeSymlink, "include-symlink", "s", false, "include symlink path in output")
	domainsCmd.AddCommand(lsCmd)
}
