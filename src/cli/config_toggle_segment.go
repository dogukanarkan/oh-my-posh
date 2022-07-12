package cli

import (
	"fmt"
	"io/ioutil"
	"oh-my-posh/color"
	"oh-my-posh/engine"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// toggleSegment represents the toggle segment command
var toggleSegment = &cobra.Command{
	Use:   "toggle",
	Short: "Enable/disable segment",
	Long: `Enable/disable segment.

You can enable/disable specific segment which given type name on runtime.

Example usage:

> oh-my-posh config toggle aws

> oh-my-posh config toggle spotify`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configFilePath := os.Getenv("POSH_THEME")

		configFile, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			fmt.Print(err)
		}

		var config engine.Config
		if err := yaml.Unmarshal([]byte(configFile), &config); err != nil {
			fmt.Print(err)
		}

		for _, block := range config.Blocks {
			for _, segment := range block.Segments {
				if string(segment.Type) == args[0] {
					segment.Enabled = !segment.Enabled
					printStatus(segment)

					break
				}
			}
		}

		modifiedConfig, err := yaml.Marshal(config)
		if err != nil {
			fmt.Print(err)
		}

		ioutil.WriteFile(configFilePath, modifiedConfig, 0644)
	},
}

func init() { // nolint:gochecknoinits
	configCmd.AddCommand(toggleSegment)
}

func printStatus(segment *engine.Segment) {
	ansiPayload := "\x1b[%sm%s\x1b[0m"
	ansiColors := color.DefaultColors{}
	red := ansiColors.AnsiColorFromString("red", false)
	green := ansiColors.AnsiColorFromString("green", false)

	var status string
	if segment.Enabled {
		status = fmt.Sprintf(ansiPayload, green, "ON")
	} else {
		status = fmt.Sprintf(ansiPayload, red, "OFF")
	}

	fmt.Printf("%s segment turned %s.\n", strings.Title(string(segment.Type)), status)
}
