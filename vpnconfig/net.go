package main

import (
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

var tlsConfig = tls.Config{InsecureSkipVerify: true}
var netDialer = net.Dialer{Timeout: 10 * time.Second}

func GetServerCertificateHash(host fmt.Stringer, port fmt.Stringer) (string, error) {
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

func VerifyHostname(host string, port string) error {
	target := fmt.Sprintf("%s:%s", host, port)
	conn, err := tls.DialWithDialer(&netDialer, "tcp", target, &tls.Config{})
	if err != nil {
		return err
	} else {
		defer conn.Close()
	}
	return conn.VerifyHostname(host)
}
