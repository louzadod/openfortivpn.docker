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
	if err!= nil {
		fmt.Printf("%s Não foi possível ler o arquivo de configuração:\n  %s\n", redDot, err)
		os.Exit(1)
	}

	if !*reconfigure && config.IsComplete() {
		os.Exit(0)
	}

	tokenChan := GetTokenCertificates()

	config.AskHost()
	config.AskPort()

	err = config.VerifyServer()
	if err != nil {
		fmt.Printf("%s Não foi possível validar o servidor [%s]:\n  %s\n", redDot, config.Host.Value(), err)
		os.Exit(1)
	}

	err = config.SelectCertificate(tokenChan)
	if err != nil {
		os.Exit(1)
	}

	config.ConfirmSavePIN()

	err = config.Save()
	if err != nil {
		fmt.Printf("Não foi possível salvar as configurações: %s\n", err)
		os.Exit(1)
	}
}
