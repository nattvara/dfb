package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/nattvara/dfb/internal/groups"
	"github.com/spf13/cobra"
)

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
	RunE: func(cmd *cobra.Command, args []string) error {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)

		header := table.Row{"#", "Group"}
		t.AppendHeader(header)

		groups, err := groups.FetchGroups()
		if err != nil {
			log.Println(err.Error())
			return errors.New("failed to read groups")
		}

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

		return nil
	},
}

var addGroupsCmd = &cobra.Command{
	Use:   "add [group]",
	Short: "Add new group",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := groups.GetGroupFromString(args[0])
		if err != nil {
			log.Println(err.Error())
			return fmt.Errorf("group %s already exists", args[0])
		}

		group := groups.New(args[0])
		group.Create()

		fmt.Printf("Group created at %s\n", group.Path)

		return nil
	},
}

var repoGroupsCmd = &cobra.Command{
	Use:   "repos [group]",
	Short: "List restic repos for a group",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		group, err := groups.GetGroupFromString(args[0])
		if err != nil {
			log.Println(err.Error())
			return errors.New("failed to list repos")
		}

		repos, err := group.Repositories()
		if err != nil {
			log.Println(err.Error())
			return errors.New("failed to list repos")
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		header := table.Row{"#", "Name", "Repository"}
		t.AppendHeader(header)
		for i, repo := range repos {
			t.AppendRow(table.Row{
				i + 1,
				repo.Name,
				repo.ResticPath,
			})
		}
		t.Render()

		return nil
	},
}
var addRepoGroupsCmd = &cobra.Command{
	Use:   "add-repo [group] [repo-name] [restic-repo]",
	Short: "Add restic repo to a group",
	Args:  cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		group, err := groups.GetGroupFromString(args[0])
		if err != nil {
			log.Println(err.Error())
			return errors.New("failed to add restic repo")
		}

		repo := groups.Repository{
			Name:       args[1],
			ResticPath: args[2],
		}

		err = group.AddRepository(repo)
		if err != nil {
			log.Println(err.Error())
			return errors.New("failed to add restic repo")
		}

		fmt.Println("repository was added")

		return nil
	},
}

func init() {
	RootCmd.AddCommand(groupsCmd)

	groupsCmd.AddCommand(lsGroupsCmd)
	groupsCmd.AddCommand(addGroupsCmd)
	groupsCmd.AddCommand(repoGroupsCmd)
	groupsCmd.AddCommand(addRepoGroupsCmd)
}
