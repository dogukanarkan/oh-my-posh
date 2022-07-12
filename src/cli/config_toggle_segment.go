package cli

import (
	"fmt"
	"io/ioutil"
	"oh-my-posh/color"
	"oh-my-posh/engine"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	configFilePath = os.Getenv("POSH_THEME")
	parsedConfig   *engine.Config
	segments       []*engine.Segment
	segmentTypes   []string
	listFlag       bool
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
		if !cmd.HasFlags() && len(args) == 0 {
			_ = cmd.Help()
			return
		}
		if listFlag {
			listSegmentsStatus()
			return
		}

		toggleSegment(args[0])
		rewriteConfigFile()
	},
}

func init() { // nolint:gochecknoinits
	toggleSegmentCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "list segments status")
	configCmd.AddCommand(toggleSegmentCmd)
}

func listSegmentTypes() []string {
	parseConfigFile()
	return segmentTypes
}

func listSegmentsStatus() {
	sortSegments()

	for _, segment := range segments {
		printColoredSegmentStatus(segment, "->")
	}
}

func sortSegments() {
	sort.Slice(segments, func(i, j int) bool {
		if segments[i].Enabled != segments[j].Enabled {
			return segments[i].Enabled && !segments[j].Enabled
		}

		return segments[i].Type < segments[j].Type
	})
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
	printColoredSegmentStatus(segment, "segment turned")
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

func printColoredSegmentStatus(segment *engine.Segment, message string) {
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

	fmt.Printf("%s %s %s.\n", strings.Title(string(segment.Type)), message, status)
}
