package segments

import (
	"fmt"
	"oh-my-posh/environment"
	"oh-my-posh/properties"
	"strings"
)

type Oci struct {
	props properties.Properties
	env   environment.Environment

	Profile string
	Region  string
}

func (a *Oci) Template() string {
	return " {{ .Profile }}@{{ .Region }}"
}

func (a *Oci) Init(props properties.Properties, env environment.Environment) {
	a.props = props
	a.env = env
}

func (a *Oci) Enabled() bool {
	a.Profile = a.env.Getenv("OCI_CLI_PROFILE")
	a.Region = a.getRegion()

	return a.Profile != ""
}

func (a *Oci) getRegion() string {
	configPath := a.env.Getenv("OCI_CLI_CONFIG_FILE")
	config := a.env.FileContent(configPath)
	configLines := strings.Split(config, "\n")
	configSection := fmt.Sprintf("[%s]", a.Profile)

	var sectionActive bool
	for _, line := range configLines {
		if strings.HasPrefix(line, configSection) {
			sectionActive = true
			continue
		}

		if sectionActive && strings.HasPrefix(line, "region") {
			splitted := strings.Split(line, "=")
			if len(splitted) >= 2 {
				return strings.TrimSpace(splitted[1])
			}
		}
	}

	return ""
}
