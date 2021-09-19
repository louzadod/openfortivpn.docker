package main

import (
	"fmt"
	"github.com/mgutz/ansi"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var cfg = kingpin.Arg("file", "Arquivo de configuração.").Required().ExistingFile()
var reconfigure = kingpin.Flag("reconfigure", "Reconfigura.").Bool()

func main() {
	var err error
	kingpin.Parse()

	tokenChan := make(chan []TokenCert, 1)
	go listCerts(tokenChan)

	config := LoadConfig(*cfg)

	if *reconfigure || config.Host.Value() == "" {
		answer := ask(config.Host.Value(), ipQuestion, ipValidate)
		config.Host.SetValue(answer)
	}

	if *reconfigure || config.Port.Value() == "" {
		answer := ask(config.Port.Value(), portQuestion, portValidate)
		config.Port.SetValue(answer)
	}

	if *reconfigure || config.TrustedCert.Value() == "" {
		var result string
		if result, err = getServerCertificateHash(config.Host, config.Port); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		config.TrustedCert.SetValue(result)
	}

	if *reconfigure || config.UserCert.Value() == "" {
		fmt.Printf("%s Carregando certificados...", blueDot)

		tokenCerts := <-tokenChan
		switch totalCerts := len(tokenCerts); totalCerts {
		case 0:
			fmt.Println("\rErro: Nenhum certificado foi encontrado")
			fmt.Println("Ou o token não foi inserido ou há mais de um token conectado")
			os.Exit(1)
		case 1:
			fmt.Printf("\r%s Só há um certificado elegível no token. Selecionado automaticamente:\n", blueDot)
			fmt.Printf("  %s\n", ansi.Color(tokenCerts[0].name, "white+d"))
			config.UserCert.SetValue(tokenCerts[0].url)
		default:
			sel(config.UserCert, certQuestion, tokenCerts)
		}

		if confirm(savePinQuestion) {
			pinValue := password(enterPinQuestion)
			config.UserCert.SetValue(fmt.Sprintf("%s;pin-value=%s", config.UserCert.Value(), pinValue))
		}
	}

	if err := config.File.SaveTo(*cfg); err != nil {
		fmt.Printf("Não foi possível salvar as configurações: %s\n", err)
	}
}
