package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/nattvara/dfb/internal/groups"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// domainsCmd represents the domains command
var domainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "Domain commands",
	Long: `
A domain is a directory in the home directory to backup,
this could be a symlink to some other directory on another
volume
`,
}

var ListIncludeRepositories bool
var ListIncludeSymlink bool
var CreatSymlinkPath string

var lsDomainsCmd = &cobra.Command{
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

var notAddedDomainsCmd = &cobra.Command{
	Use:   "not-added [group]",
	Short: "List directories not added as domains in home directory",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		group, err := groups.GetGroupFromString(args[0])
		if err != nil {
			log.Fatalf(err.Error())
		}

		homedir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		files, err := ioutil.ReadDir(homedir)
		if err != nil {
			log.Fatal(err)
		}

		var notFound []string
		domains := group.DomainsMap()
		for _, f := range files {
			if _, ok := domains[f.Name()]; !ok {
				notFound = append(notFound, f.Name())
			}
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		header := table.Row{"#", "Name"}
		t.AppendHeader(header)
		for i, path := range notFound {
			t.AppendRow(table.Row{
				i,
				path,
			})
		}
		t.Render()
	},
}

var addDomainsCmd = &cobra.Command{
	Use:   "add [group] [domain]",
	Short: "Add new domain",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		group, err := groups.GetGroupFromString(args[0])
		if err != nil {
			log.Fatalf(err.Error())
		}

		domainName := args[1]

		if group.DomainExists(domainName) {
			fmt.Println("Domain already exists.")
			return
		}

		if CreatSymlinkPath != "" {
			group.AddDomainWithNameAndSymlink(domainName, CreatSymlinkPath)
		} else {
			group.AddDomainWithName(domainName)
		}

		fmt.Println("Domain created.")
	},
}

var rmDomainsCmd = &cobra.Command{
	Use:   "rm [group] [domain]",
	Short: "Remove a domain",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		group, err := groups.GetGroupFromString(args[0])
		if err != nil {
			log.Fatalf(err.Error())
		}
		domainName := args[1]

		if !group.DomainExists(domainName) {
			log.Fatalf("Domain %s doesn't exists.", domainName)
		}

		domain := group.DomainsMap()[domainName]

		fmt.Printf("deleting record of domain %s\n", domainName)
		domain.DeleteConfigFile()

		if domain.Symlink != nil {
			fmt.Println("deleting symlink to real directory")
			domain.Symlink.DeleteProxy()
		}

		fmt.Printf(`
NOTE: this does not delete the actual directory
      it will simply not be included in any more backups
      neither will it be removed from previous backups
      you will have to delete the data of "%s" yourself
		`, domainName)
	},
}

func init() {
	RootCmd.AddCommand(domainsCmd)

	lsDomainsCmd.Flags().BoolVarP(&ListIncludeRepositories, "include-repositories", "r", false, "include repositories in output")
	lsDomainsCmd.Flags().BoolVarP(&ListIncludeSymlink, "include-symlink", "s", false, "include symlink path in output")
	domainsCmd.AddCommand(lsDomainsCmd)

	domainsCmd.AddCommand(notAddedDomainsCmd)

	addDomainsCmd.Flags().StringVarP(&CreatSymlinkPath, "symlink", "s", "", "domain content is symlinked to another location")
	domainsCmd.AddCommand(addDomainsCmd)

	domainsCmd.AddCommand(rmDomainsCmd)
}
