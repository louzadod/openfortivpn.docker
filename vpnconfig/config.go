package main

import (
	"errors"
	"fmt"
	"github.com/mgutz/ansi"
	"gopkg.in/ini.v1"
	"os"
)

type VPNConfig struct {
	FileName    string
	File        *ini.File
	Host        *ini.Key
	Port        *ini.Key
	UserCert    *ini.Key
	TrustedCert *ini.Key
}

var ErrCertsNotFound = errors.New("certs not found")

func LoadConfig(cfgFile string) (*VPNConfig, error) {
	iniFile, err := ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, cfgFile)

	if err != nil {
		return nil, err
	}

	section := iniFile.Section("")
	return &VPNConfig{
		FileName:    cfgFile,
		File:        iniFile,
		Host:        section.Key("host"),
		Port:        section.Key("port"),
		UserCert:    section.Key("user-cert"),
		TrustedCert: section.Key("trusted-cert"),
	}, nil
}

func (c *VPNConfig) IsComplete() bool {
	req := c.Host.Value() != "" && c.Port.Value() != "" && c.UserCert.Value() != ""
	optional := IsDNSName(c.Host.Value()) || c.TrustedCert.Value() != ""
	return req && optional
}

func (c *VPNConfig) DeleteKey(key string) {
	section := c.File.Section("")
	section.DeleteKey(key)
}

func (c *VPNConfig) VerifyServerHostname() error {
	return VerifyHostname(c.Host.Value(), c.Port.Value())
}

func (c *VPNConfig) Save() error {
	if c.TrustedCert.Value() == "" {
		c.DeleteKey("trusted-cert")
	}
	return c.File.SaveTo(c.FileName)
}

func (c *VPNConfig) IsNameBased() bool {
	return IsDNSName(c.Host.Value())
}

func (c *VPNConfig) AskPort() {
	answer := ask(c.Port.Value(), portQuestion, portValidate)
	c.Port.SetValue(answer)
}

func (c *VPNConfig) AskHost() {
	answer := ask(c.Host.Value(), ipQuestion, ipValidate)
	c.Host.SetValue(answer)
}

func (c *VPNConfig) ConfirmCertificate() {
	if hash, err := c.VerifyServer(); err != nil {
		if hash == c.TrustedCert.Value() {
			return
		}
		if hash == "" {
			fmt.Printf("%s Não foi possível obter um certificado do servidor\n  %s", redDot, err.Error())
			os.Exit(1)
		}
		fmt.Printf("%s Não foi possível validar o certificado do servidor\n  Fingerprint: %s\n", redDot, hash)
		if confirm(invalidCertificateQuestion) {
			c.TrustedCert.SetValue(hash)
		} else {
			os.Exit(1)
		}
	}
}

func (c *VPNConfig) SelectCertificate(tokenChan chan []TokenCert) error {
	fmt.Printf("%s Carregando certificados...", blueDot)
	tokenCerts := <-tokenChan
	switch totalCerts := len(tokenCerts); totalCerts {
	case 0:
		fmt.Printf("\r%s Erro! Nenhum certificado foi encontrado.\n", redDot)
		fmt.Println("  - O token foi inserido?")
		fmt.Println("  - Existe algum certificado elegível no token?")
		fmt.Println("  - Há mais de um token conectado?")
		return ErrCertsNotFound
	case 1:
		fmt.Printf("\r%s Só há um certificado elegível no token. Selecionado automaticamente:\n", blueDot)
		fmt.Printf("  %s\n", ansi.Color(tokenCerts[0].name, "white+d"))
		c.UserCert.SetValue(tokenCerts[0].url)
		return nil
	default:
		c.UserCert.SetValue(sel(certQuestion, tokenCerts))
		return nil
	}
}

func (c *VPNConfig) ConfirmSavePIN() {
	if confirm(savePinQuestion) {
		pinValue := password(enterPinQuestion)
		c.UserCert.SetValue(fmt.Sprintf("%spin-value=%s", c.UserCert.Value(), pinValue))
	}
}

func (c *VPNConfig) VerifyServer() (string, error) {
	var hash string
	var err error
	if hash, err = GetServerCertificateHash(c.Host, c.Port); err != nil {
		return "", err
	}
	return hash, c.VerifyServerHostname()
}
