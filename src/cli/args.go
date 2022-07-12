package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NoArgsOrOneValidArg(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
	if err := cobra.ExactArgs(1)(cmd, args); err != nil {
		return err
	}
	if err := cobra.OnlyValidArgs(cmd, args); err != nil {
		fmt.Printf("Available segments: %s\n\n", cmd.ValidArgs)
		return err
	}

	return nil
}
