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

var (
	configFilePath = os.Getenv("POSH_THEME")
	parsedConfig   *engine.Config
	segments       []*engine.Segment
	segmentTypes   []string
)

// toggleSegment represents the toggle segment command
var toggleSegmentCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Enable/disable segment",
	Long: `Enable/disable segment.

You can enable/disable specific segment which given type name on runtime.

Example usage:

> oh-my-posh config toggle aws

> oh-my-posh config toggle spotify`,
	ValidArgs: listSegmentTypes(),
	Args:      NoArgsOrOneValidArg,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}

		toggleSegment(args[0])
		rewriteConfigFile()
	},
}

func init() { // nolint:gochecknoinits
	configCmd.AddCommand(toggleSegmentCmd)
}

func listSegmentTypes() []string {
	parseConfigFile()
	return segmentTypes
}

func parseConfigFile() {
	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Print(err)
	}

	if err := yaml.Unmarshal([]byte(configFile), &parsedConfig); err != nil {
		fmt.Print(err)
	}

	for _, block := range parsedConfig.Blocks {
		for _, segment := range block.Segments {
			segments = append(segments, segment)
			segmentTypes = append(segmentTypes, string(segment.Type))
		}
	}
}

func toggleSegment(typeName string) {
	segment, err := findSegmentByTypeName(typeName)
	if err != nil {
		return
	}

	segment.Enabled = !segment.Enabled
	printStatus(segment)
}

func findSegmentByTypeName(typeName string) (*engine.Segment, error) {
	for _, segment := range segments {
		if string(segment.Type) == typeName {
			return segment, nil
		}
	}

	return nil, fmt.Errorf("segment %s not found.", typeName)
}

func rewriteConfigFile() {
	modifiedConfig, err := yaml.Marshal(parsedConfig)
	if err != nil {
		fmt.Print(err)
	}

	ioutil.WriteFile(configFilePath, modifiedConfig, 0644)
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
