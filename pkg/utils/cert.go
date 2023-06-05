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
	// 100 years
	DefaultCAduration = time.Hour * 24 * 365 * 100
)

// host could be ip or dns
func NewServerCertKey(host string, alternateIPs []net.IP, alternateDNS []string) (serverCert, serverKey, caCert []byte, err error) {

	// ----------- generate self-signed ca
	validFrom := time.Now().Add(-time.Hour) // valid an hour earlier to avoid flakes due to clock skew

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

	caCertificate, err := x509.ParseCertificate(caDERBytes)
	if err != nil {
		return nil, nil, nil, err
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
func NewServerCertKeyForLocalNode(alternateDNS []string, destCertFilePath, destKeyFilePath, destCaFilePath string) error {
	host, _ := os.Hostname()
	alternateIPs := []net.IP{}

	ipv4List, ipv6List, err := GetAllInterfaceUnicastAddrWithoutMask()
	if err != nil {
		return errors.Errorf("fail to get local ip, error=%v", err)
	}
	alternateIPs = append(alternateIPs, ipv4List...)
	alternateIPs = append(alternateIPs, ipv6List...)

	serverCert, serverKey, caCert, err := NewServerCertKey(host, alternateIPs, alternateDNS)
	if err != nil {
		return err
	}

	if err := cert.WriteCert(destCertFilePath, serverCert); err != nil {
		return err
	}
	if err := cert.WriteCert(destCaFilePath, caCert); err != nil {
		return err
	}
	if err := keyutil.WriteKey(destKeyFilePath, serverKey); err != nil {
		return err
	}
	return nil
}

// CanReadCertAndKey returns true if the certificate and key files already exists,
// otherwise returns false. If lost one of cert and key, returns error.
func CanReadCertAndKey(certPath, keyPath string) (bool, error) {
	return cert.CanReadCertAndKey(certPath, keyPath)
}
