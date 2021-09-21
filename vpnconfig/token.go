package main

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

type TokenCert struct {
	name string
	url  string
}

var tlsConfig = tls.Config{InsecureSkipVerify: true}
var netDialer = net.Dialer{Timeout: 10 * time.Second}

func getServerCertificateHash(host fmt.Stringer, port fmt.Stringer) (string, error) {
	address := fmt.Sprintf("%s:%s", host, port)

	conn, err := tls.DialWithDialer(&netDialer, "tcp", address, &tlsConfig)
	if err != nil {
		return "", fmt.Errorf("não foi possível obter o certificado do servidor: %s", err.Error())
	} else {
		defer conn.Close()
	}
	firstCert := conn.ConnectionState().PeerCertificates[0]
	return fmt.Sprintf("%x", sha256.Sum256(firstCert.Raw)), nil
}

func getCertByURL(certUrl string) (TokenCert, error) {
	pemBytes, err := exec.Command("p11tool", "--export", certUrl).Output()
	if err != nil {
		return TokenCert{}, err
	}
	block, _ := pem.Decode(pemBytes)
	cert, _ := x509.ParseCertificate(block.Bytes)
	// não considera certificados CA
	if cert.IsCA {
		return TokenCert{}, nil
	}

	return TokenCert{
		name: fmt.Sprintf("%s [%s]", cert.Subject.CommonName, cert.Issuer.CommonName),
		url:  certUrl,
	}, nil
}

func listCerts(ch chan<- []TokenCert) {
	var err error
	var output []byte
	var cert TokenCert

	output, _ = exec.Command("p11tool", "--list-all", "--only-urls").Output()
	certsUrls := strings.Split(strings.TrimSpace(string(output)), "\n")
	subjects := make([]TokenCert, 0, len(certsUrls))

	for _, certUrl := range certsUrls {
		if cert, err = getCertByURL(certUrl); err != nil {
			continue
		}
		if cert.url != "" {
			subjects = append(subjects, cert)
		}
	}

	ch <- subjects
}
