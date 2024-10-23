package main

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"os"
)

func ask(defaultValue string, input *survey.Input, validator survey.Validator) string {
	var result string
	if defaultValue != "" {
		input.Default = defaultValue
	}
	err := survey.AskOne(input, &result, survey.WithValidator(validator))
	if errors.Is(err, terminal.InterruptErr) {
		os.Exit(1)
	}
	return result
}

func confirm(input *survey.Confirm) bool {
	result := false
	err := survey.AskOne(input, &result)
	if errors.Is(err, terminal.InterruptErr) {
		os.Exit(1)
	}
	return result
}

func password(input *survey.Password) string {
	var result string
	err := survey.AskOne(input, &result)
	if errors.Is(err, terminal.InterruptErr) {
		os.Exit(1)
	}
	return result
}

func sel(input *survey.Select, tokenCerts []TokenCert) string {
	var index int
	for _, cert := range tokenCerts {
		input.Options = append(input.Options, cert.name)
	}
	err := survey.AskOne(certQuestion, &index)
	if errors.Is(err, terminal.InterruptErr) {
		os.Exit(1)
	}
	return tokenCerts[index].url
}

var ipQuestion = &survey.Input{Message: "IP ou domínio do servidor da VPN:"}
var portQuestion = &survey.Input{Message: "Porta:", Default: "443"}
var certQuestion = &survey.Select{Message: "Selecione o certificado:"}

var savePinQuestion = &survey.Confirm{
	Message: "Deseja salvar o PIN para não precisar digitá-lo a cada conexão?", Default: true,
}
var invalidCertificateQuestion = &survey.Confirm{
	Message: "Gostaria de se conectar mesmo assim?",
	Default: false,
}
var enterPinQuestion = &survey.Password{Message: "PIN:"}
var lastTryPinQuestion = &survey.Confirm{
	Message: "É a última tentativa do PIN. Deseja continuar?", Default: true,
}
