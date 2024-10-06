package main

import (
	"fmt"
	"os"

	"singularity/straphangerctl/stations"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use: "straphangerctl",
	Short: "Help.",
	Long: "View help page and global flags",
}

func init() {
	cmd.AddCommand(stations.Command)
}

func main(){
	cmd.Execute()
}

func Execute() {
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("Unable to execute command: %v\n", err)
		os.Exit(1)
	}
}

