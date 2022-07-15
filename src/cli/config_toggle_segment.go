package cli

import (
	"fmt"
	"oh-my-posh/color"
	"oh-my-posh/engine"
	"oh-my-posh/environment"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var (
	parsedConfig *engine.Config
	segments     []*engine.Segment
	segmentTypes []string
	listFlag     bool
	orderFlag    bool
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
		if !listFlag && len(args) == 0 {
			_ = cmd.Help()
			return
		}
		if listFlag {
			listSegmentsStatus()
			return
		}

		toggleSegment(args[0])
		engine.SyncAndWrite(parsedConfig)
	},
}

func init() { // nolint:gochecknoinits
	toggleSegmentCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "list segments status")
	toggleSegmentCmd.Flags().BoolVarP(&orderFlag, "order", "o", false, "order segments by status and name")
	configCmd.AddCommand(toggleSegmentCmd)
}

func listSegmentTypes() []string {
	parseConfigFile()
	return segmentTypes
}

func listSegmentsStatus() {
	if orderFlag {
		sortSegments()
	}

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
	env := &environment.ShellEnvironment{
		Version: cliVersion,
	}
	env.Init()
	defer env.Close()

	parsedConfig = engine.LoadConfig(env)

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
