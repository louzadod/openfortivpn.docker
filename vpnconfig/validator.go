package main

import (
	"errors"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var DNSName = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`
var rxDNSName = regexp.MustCompile(DNSName)

func IsIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func IsDNSName(str string) bool {
	if str == "" || len(strings.Replace(str, ".", "", -1)) > 255 {
		return false
	}
	return !IsIP(str) && rxDNSName.MatchString(str)
}

func ipValidate(val interface{}) error {
	host := val.(string)
	if IsDNSName(host) || IsIP(host) {
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
