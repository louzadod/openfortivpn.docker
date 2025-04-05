package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var cfg = kingpin.Arg("file", "Arquivo de configuração.").Required().ExistingFile()
var reconfigure = kingpin.Flag("reconfigure", "Reconfigura.").Bool()

func main() {
	var err error
	var config *VPNConfig
	kingpin.Parse()

	config, err = LoadConfig(*cfg)
	if err != nil {
		fmt.Printf("%s Não foi possível ler o arquivo de configuração:\n  %s\n", redDot, err)
		os.Exit(1)
	}

	if !*reconfigure {
		checkToken()
	}

	if !*reconfigure && config.IsComplete() {
		config.ConfirmCertificate()
		save(config)
		os.Exit(0)
	}

	tokenChan := GetTokenCertificates()

	config.AskHost()
	config.AskPort()
	config.ConfirmCertificate()

	err = config.SelectCertificate(tokenChan)
	if err != nil {
		os.Exit(1)
	}

	config.ConfirmSavePIN()

	save(config)
}

func save(config *VPNConfig) {
	if err := config.Save(); err != nil {
		fmt.Printf("Não foi possível salvar as configurações: %s\n", err)
		os.Exit(1)
	}
}

func checkToken() {
	info, err := GetTokenInfo()
	if err != nil {
		fmt.Printf("%s Não foi possível inicializar o token: %s\n", redDot, err)
		os.Exit(1)
	}

	if info.locked {
		fmt.Printf("%s O PIN do token está bloqueado.\n", redDot)
		os.Exit(1)
	}

	if info.finalTry {
		if !confirm(lastTryPinQuestion) {
			os.Exit(1)
		}
	}
}
