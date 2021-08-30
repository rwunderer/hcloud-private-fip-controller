package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	ApiToken    string `yaml:"hcloudToken"`
	HostIp      string `yaml:"hostIP"`
	IpAddress   string `yaml:"ipAddress"`
	NetworkName string `yaml:"networkName"`
}

func (conf *Config) ReadFile(inFile string) error {
	content, err := ioutil.ReadFile(inFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	if err := yaml.Unmarshal(content, &conf); err != nil {
		return fmt.Errorf("failed to unmarshal config file: %v", err)
	}

	return nil
}

func (conf *Config) Validate() error {
	var errs []string

	if conf.ApiToken == "" {
		errs = append(errs, "API token not provided. Please set HCLOUD_TOKEN env var!")
	}

	if conf.HostIp == "" {
		errs = append(errs, "Host IP not provided. Please set HOST_IP env var!")
	}

	if conf.IpAddress == "" {
		errs = append(errs, "Floating IP address not provided. Please set IP_ADDRESS env var!")
	}

	if conf.NetworkName == "" {
		errs = append(errs, "Network name not provided. Please set HCLOUD_NETWORK env var!")
	}

	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, ", "))
	}

	return nil
}
