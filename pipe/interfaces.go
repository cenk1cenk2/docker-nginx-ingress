package pipe

type (
	ConfigurationServerJson struct {
		Listen  string            `json:"listen"            validate:"required"`
		Options map[string]string `json:"options,omitempty"`
	}

	ConfigurationUpstreamJson struct {
		Servers []string          `json:"servers"           validate:"required"`
		Options map[string]string `json:"options,omitempty"`
	}

	ConfigurationJson struct {
		Server   ConfigurationServerJson   `json:"server"   validate:"required"`
		Upstream ConfigurationUpstreamJson `json:"upstream" validate:"required"`
	}

	Configuration []ConfigurationJson

	ServerTemplate struct {
		Listen   string
		Upstream string
		Options  map[string]string
	}

	UpstreamTemplate struct {
		Upstream string
		Servers  []string
		Options  map[string]string
	}
)
