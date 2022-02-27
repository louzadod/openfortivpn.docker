package main

import (
	"gopkg.in/ini.v1"
)

type VPNConfig struct {
	File        *ini.File
	Host        *ini.Key
	Port        *ini.Key
	UserCert    *ini.Key
	TrustedCert *ini.Key
}

func LoadConfig(cfgFile string) VPNConfig {
	iniFile, _ := ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, cfgFile)

	section := iniFile.Section("")
	return VPNConfig{
		File:        iniFile,
		Host:        section.Key("host"),
		Port:        section.Key("port"),
		UserCert:    section.Key("user-cert"),
		TrustedCert: section.Key("trusted-cert"),
	}
}

func (c VPNConfig) IsComplete() bool {
	req := c.Host.Value() != "" && c.Port.Value() != "" && c.UserCert.Value() != ""
	optional := IsDNS(c.Host.Value()) || c.TrustedCert.Value() != ""
	return req && optional
}
