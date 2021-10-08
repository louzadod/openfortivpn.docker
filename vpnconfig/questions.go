package main

import (
	"errors"
	"net"
	"os"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"gopkg.in/ini.v1"
)

func ipValidate(ip interface{}) error {
	if net.ParseIP(ip.(string)) == nil {
		return errors.New("IP inválido")
	}
	return nil
}

func portValidate(val interface{}) error {
	port, _ := strconv.Atoi(val.(string))
	if port < 1 || port > 65536 {
		return errors.New("porta inválida")
	}
	return nil
}

func ask(defaultValue string, input *survey.Input, validator survey.Validator) string {
	var result string
	if defaultValue != "" {
		input.Default = defaultValue
	}
	err := survey.AskOne(input, &result, survey.WithValidator(validator))
	if err == terminal.InterruptErr {
		os.Exit(1)
	}
	return result
}

func confirm(input *survey.Confirm) bool {
	result := false
	err := survey.AskOne(input, &result)
	if err == terminal.InterruptErr {
		os.Exit(1)
	}
	return result
}

func password(input *survey.Password) string {
	var result string
	err := survey.AskOne(input, &result)
	if err == terminal.InterruptErr {
		os.Exit(1)
	}
	return result
}

func sel(key *ini.Key, input *survey.Select, tokenCerts []TokenCert) {
	var index int
	for _, cert := range tokenCerts {
		input.Options = append(input.Options, cert.name)
	}
	err := survey.AskOne(certQuestion, &index)
	if err == terminal.InterruptErr {
		os.Exit(1)
	}
	key.SetValue(tokenCerts[index].url)
}

var ipQuestion = &survey.Input{Message: "IP do servidor VPN:"}
var portQuestion = &survey.Input{Message: "Porta:", Default: "443"}
var certQuestion = &survey.Select{Message: "Selecione o certificado:"}

var savePinQuestion = &survey.Confirm{
	Message: "Deseja salvar o PIN para não precisar digitá-lo a cada conexão?", Default: true,
}
var enterPinQuestion = &survey.Password{Message: "PIN:"}
