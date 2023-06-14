// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	"github.com/pkg/errors"
	"k8s.io/client-go/util/cert"
	"k8s.io/client-go/util/keyutil"
	netutils "k8s.io/utils/net"
	"math/big"
	"net"
	"os"
	"time"
)

const (
	// 200 years
	DefaultCAduration = time.Hour * 24 * 365 * 200
)

// host could be ip or dns
func NewServerCertKey(host string, alternateIPs []net.IP, alternateDNS []string, caCertPath, caKeyPath string) (serverCert, serverKey, caCert []byte, err error) {
	validFrom := time.Now().Add(-time.Hour) // valid an hour earlier to avoid flakes due to clock skew
	var caCertificate *x509.Certificate
	var caKey *rsa.PrivateKey
	var caDERBytes []byte
	// ----------- generate self-signed ca
	if caCertPath == "" || caKeyPath == "" {

		caKey, err := rsa.GenerateKey(cryptorand.Reader, 2048)

		if err != nil {
			return nil, nil, nil, err
		}

		CommonName := types.TlsCaCommonName
		Organization := []string{types.TlsCaCommonName}

		tmpl := x509.Certificate{
			SerialNumber: new(big.Int).SetInt64(0),
			Subject: pkix.Name{
				CommonName:   CommonName,
				Organization: Organization,
			},
			DNSNames:              []string{CommonName},
			NotBefore:             validFrom.UTC(),
			NotAfter:              validFrom.Add(DefaultCAduration).UTC(),
			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			BasicConstraintsValid: true,
			IsCA:                  true,
		}

		caDERBytes, err := x509.CreateCertificate(cryptorand.Reader, &tmpl, &tmpl, caKey.Public(), caKey)
		if err != nil {
			return nil, nil, nil, err
		}

		caCertificate, err = x509.ParseCertificate(caDERBytes)
		if err != nil {
			return nil, nil, nil, err
		}
	} else {
		//read ca cert
		caDERBytes, err = os.ReadFile(caCertPath)
		if err != nil {
			return nil, nil, nil, err
		}
		blockCert, _ := pem.Decode(caDERBytes)
		caCertificate, err = x509.ParseCertificate(blockCert.Bytes)
		if err != nil {
			return nil, nil, nil, err
		}

		//read ca key
		cakeyByte, err := os.ReadFile(caKeyPath)
		if err != nil {
			return nil, nil, nil, err
		}
		blockKey, _ := pem.Decode(cakeyByte)
		caKey, err = x509.ParsePKCS1PrivateKey(blockKey.Bytes)
		if err != nil {
			return nil, nil, nil, err
		}

	}

	// --------------- server cert

	privatekey, err := rsa.GenerateKey(cryptorand.Reader, 2048)
	if err != nil {
		return nil, nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			CommonName: fmt.Sprintf("%s@%d", host, time.Now().Unix()),
		},
		NotBefore:             validFrom.UTC(),
		NotAfter:              validFrom.Add(DefaultCAduration).UTC(),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	if ip := netutils.ParseIPSloppy(host); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, host)
	}

	template.IPAddresses = append(template.IPAddresses, alternateIPs...)
	template.DNSNames = append(template.DNSNames, alternateDNS...)

	derBytes, err := x509.CreateCertificate(cryptorand.Reader, &template, caCertificate, &privatekey.PublicKey, caKey)
	if err != nil {
		return nil, nil, nil, err
	}

	// Generate cert
	certBuffer := bytes.Buffer{}
	if err := pem.Encode(&certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return nil, nil, nil, err
	}

	caBuffer := bytes.Buffer{}
	if err := pem.Encode(&caBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: caDERBytes}); err != nil {
		return nil, nil, nil, err
	}

	// Generate key
	keyBuffer := bytes.Buffer{}
	if err := pem.Encode(&keyBuffer, &pem.Block{Type: keyutil.RSAPrivateKeyBlockType, Bytes: x509.MarshalPKCS1PrivateKey(privatekey)}); err != nil {
		return nil, nil, nil, err
	}

	return certBuffer.Bytes(), keyBuffer.Bytes(), caBuffer.Bytes(), nil
}

// generate cert for local host name and local ip, and write to files
// alternateDNS could be pod dns name
func NewServerCertKeyForLocalNode(alternateDNS []string, alternateIPs []net.IP, caCertPath, caKeyPath, destCertFilePath, destKeyFilePath, destCaCertPath string) error {
	host, _ := os.Hostname()

	ipv4List, ipv6List, err := GetAllInterfaceUnicastAddrWithoutMask()
	if err != nil {
		return errors.Errorf("fail to get local ip, error=%v", err)
	}
	alternateIPs = append(alternateIPs, ipv4List...)
	alternateIPs = append(alternateIPs, ipv6List...)

	serverCert, serverKey, caCert, err := NewServerCertKey(host, alternateIPs, alternateDNS, caCertPath, caKeyPath)
	if err != nil {
		return err
	}

	if err := cert.WriteCert(destCertFilePath, serverCert); err != nil {
		return err
	}

	if err := cert.WriteCert(destCaCertPath, caCert); err != nil {
		return err
	}

	if err := keyutil.WriteKey(destKeyFilePath, serverKey); err != nil {
		return err
	}
	return nil
}
