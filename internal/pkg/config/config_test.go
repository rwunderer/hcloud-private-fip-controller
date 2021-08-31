package config

import (
	"fmt"
	"testing"
)

func testConfig() *Config {
	return &Config{
		ApiToken:  "aToken",
		HostIp: "2.3.4.5",
		IpAddress: "1.2.3.4",
		NetworkName: "test",
	}
}

type GenConfiguration func() *Config

func TestValidate(t *testing.T) {
	tests := []struct {
		name   string
		config GenConfiguration
		err    error
	}{
		{
			name: "test valid config",
			config: func() *Config {
				return testConfig()
			},
			err: nil,
		},
		{
			name: "test no token",
			config: func() *Config {
				conf := testConfig()
				conf.ApiToken = ""
				return conf
			},
			err: fmt.Errorf("API token not provided. Please set HCLOUD_TOKEN env var!"),
		},
		{
			name: "test no ip",
			config: func() *Config {
				conf := testConfig()
				conf.IpAddress = ""
				return conf
			},
			err: fmt.Errorf("Floating IP address not provided. Please set IP_ADDRESS env var!"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conf := test.config()
			err := conf.Validate()

			if err == nil {
				if test.err != nil {
					t.Fatalf("error should be [%v] but was [nil]", test.err)
				}
			} else {
				if test.err == nil {
					t.Fatalf("error should be [nil] but was [%v]", err)
				}
				if err.Error() != test.err.Error() {
					t.Fatalf("error should be [%v] but was [%v]", test.err, err)
				}
			}
		})
	}
}
