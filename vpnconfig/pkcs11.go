package main

import (
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/miekg/pkcs11"
	"net/url"
	"strings"
)

type TokenCert struct {
	name string
	url  string
}

type TokenInfo struct {
	locked   bool
	finalTry bool
	countLow bool
}

func GetTokenCertificates() chan []TokenCert {
	tokenChan := make(chan []TokenCert, 1)
	go listCerts(tokenChan)
	return tokenChan
}

var searchTemplate = []*pkcs11.Attribute{
	pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_CERTIFICATE),
}
var attrTemplate = []*pkcs11.Attribute{
	pkcs11.NewAttribute(pkcs11.CKA_ID, nil),
	pkcs11.NewAttribute(pkcs11.CKA_LABEL, nil),
	pkcs11.NewAttribute(pkcs11.CKA_VALUE, nil),
}

func listCerts(ch chan<- []TokenCert) {
	tokenCerts, _ := getAcceptableCertificates()
	ch <- tokenCerts
}

func initToken() (*pkcs11.Ctx, error) {
	ctx := pkcs11.New("/usr/lib/libeToken.so")
	err := ctx.Initialize()
	return ctx, err
}

func GetTokenInfo() (TokenInfo, error) {
	info := TokenInfo{}

	ctx, err := initToken()
	if err != nil {
		return info, err
	}
	defer ctx.Destroy()
	defer ctx.Finalize()

	slots, err := ctx.GetSlotList(true)
	if err != nil {
		return info, err
	}

	if len(slots) == 0 {
		// returns empty token error
		return info, errors.New("not found")
	}

	tokenInfo, err := ctx.GetTokenInfo(slots[0])
	if err != nil {
		return info, err
	}

	info.locked = tokenInfo.Flags&pkcs11.CKF_USER_PIN_LOCKED != 0
	info.finalTry = tokenInfo.Flags&pkcs11.CKF_USER_PIN_FINAL_TRY != 0
	info.countLow = tokenInfo.Flags&pkcs11.CKF_USER_PIN_COUNT_LOW != 0

	return info, nil
}

func getAcceptableCertificates() ([]TokenCert, error) {
	var certs []TokenCert

	p, err := initToken()
	if err != nil {
		return nil, err
	}

	defer p.Destroy()
	defer p.Finalize()

	slots, err := p.GetSlotList(true)
	if err != nil {
		return nil, err
	}

	if len(slots) == 0 {
		return []TokenCert{}, nil
	}

	tokenInfo, _ := p.GetTokenInfo(slots[0])

	session, err := p.OpenSession(slots[0], pkcs11.CKF_SERIAL_SESSION)
	if err != nil {
		return nil, err
	}
	defer p.CloseSession(session)

	err = p.FindObjectsInit(session, searchTemplate)
	if err != nil {
		return nil, err
	}

	hObjects, _, err := p.FindObjects(session, 1024)
	if err != nil {
		return nil, err
	}

	for _, hObject := range hObjects {
		attrs, err := p.GetAttributeValue(session, hObject, attrTemplate)
		if err != nil {
			continue
		}

		cert, err := x509.ParseCertificate(attrs[2].Value)
		if err != nil {
			continue
		}

		if cert.IsCA {
			continue
		}

		id, label := attrs[0], attrs[1]
		certs = append(certs, TokenCert{
			name: fmt.Sprintf("%s [%s]", cert.Subject.CommonName, cert.Issuer.CommonName),
			url:  getCertificateURL(tokenInfo, id, label),
		})
	}

	_ = p.FindObjectsFinal(session)

	return certs, nil
}

func percentEncode(value []byte) string {
	var buf strings.Builder
	for _, v := range value {
		buf.WriteString(fmt.Sprintf("%%%02X", v))
	}
	return buf.String()
}

func formatPkcs11UriParam(name string, value string, buf *strings.Builder) {
	buf.WriteString(fmt.Sprintf("%s=%s;", name, value))
}

func getCertificateURL(info pkcs11.TokenInfo, id *pkcs11.Attribute, object *pkcs11.Attribute) string {
	var buf strings.Builder
	buf.WriteString("pkcs11:")
	formatPkcs11UriParam("model", url.PathEscape(info.Model), &buf)
	formatPkcs11UriParam("manufacturer", url.PathEscape(info.ManufacturerID), &buf)
	formatPkcs11UriParam("serial", url.PathEscape(info.SerialNumber), &buf)
	formatPkcs11UriParam("token", url.PathEscape(info.Label), &buf)
	formatPkcs11UriParam("id", percentEncode(id.Value), &buf)
	formatPkcs11UriParam("object", url.PathEscape(string(object.Value)), &buf)
	formatPkcs11UriParam("type", "cert", &buf)
	return buf.String()
}
