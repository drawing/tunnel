package ca

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"log"
	"math/big"
	"time"
)

var now = time.Now()

var caCertificate *x509.Certificate
var caPriKey *rsa.PrivateKey

func LoadCA() {
	if caCertificate != nil {
		return
	}
	certpem, err := ioutil.ReadFile("config/ca.pem")
	if err != nil {
		log.Fatalln("ca cert", err)
	}
	keypem, err := ioutil.ReadFile("config/ca.key")
	if err != nil {
		log.Fatalln("ca key", err)
	}

	certblock, _ := pem.Decode(certpem)

	caCertificate, err = x509.ParseCertificate(certblock.Bytes)
	if err != nil {
		log.Fatalln("ca key", err)
	}

	keyblock, _ := pem.Decode(keypem)
	caPriKey, err = x509.ParsePKCS1PrivateKey(keyblock.Bytes)
	if err != nil {
		log.Fatalln("ca key", err)
	}
}

func Generate() ([]byte, []byte, error) {
	LoadCA()

	template := &x509.Certificate{
		IsCA: false,
		BasicConstraintsValid: true,
		SubjectKeyId:          []byte{1, 2, 3},
		SerialNumber:          big.NewInt(1238),
		Subject: pkix.Name{
			Country:      []string{"Earth"},
			Organization: []string{"Mother Nature"},
		},
		DNSNames:  []string{"www.baidu.com"},
		NotBefore: now,
		NotAfter:  time.Now().AddDate(5, 5, 5),
		// see http://golang.org/pkg/crypto/x509/#KeyUsage
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	// generate private key
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	publickey := &privatekey.PublicKey

	// create a self-signed certificate. template = parent
	var parent = caCertificate
	cert, err := x509.CreateCertificate(rand.Reader, template, parent, publickey, caPriKey)
	if err != nil {
		return nil, nil, err
	}

	var pemblock = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privatekey)}

	var certblock = &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert}

	pemkey := pem.EncodeToMemory(pemblock)
	pemcert := pem.EncodeToMemory(certblock)

	ioutil.WriteFile("cer.pem", pemcert, 0777)
	ioutil.WriteFile("pri.pem", pemkey, 0777)

	return pemcert, pemkey, nil
}
