package main

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"net"
	"strconv"
)

func IsIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func IsDNS(host string) bool {
	return govalidator.IsDNSName(host)
}

func ipValidate(val interface{}) error {
	host := val.(string)
	if IsDNS(host) || IsIP(host) {
		return nil
	}
	return errors.New("IP ou host inválido")
}

func portValidate(val interface{}) error {
	port, _ := strconv.Atoi(val.(string))
	if port < 1 || port > 65536 {
		return errors.New("porta inválida")
	}
	return nil
}
