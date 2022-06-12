package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dfb",
	Short: "Domain based Filesystem Backup.",
	Long: `
_________/\\\__________/\\\\\___/\\\________
 ________\/\\\________/\\\///___\/\\\________
  ________\/\\\_______/\\\_______\/\\\________
   ________\/\\\____/\\\\\\\\\____\/\\\________
    ___/\\\\\\\\\___\////\\\//_____\/\\\\\\\\\__
     __/\\\////\\\______\/\\\_______\/\\\////\\\_
      _\/\\\__\/\\\______\/\\\_______\/\\\__\/\\\_
       _\//\\\\\\\\\______\/\\\_______\/\\\\\\\\\\_
        __\/////////_______\///________\//////////__
`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
