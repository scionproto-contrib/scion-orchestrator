package scionca

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ThalesGroup/crypto11"
	"github.com/netsys-lab/scion-orchestrator/conf"
	"github.com/netsys-lab/scion-orchestrator/pkg/fileops"
	"github.com/scionproto/scion/pkg/scrypto/cppki"
)

type SCIONCertificateAuthority struct {
	ConfigDir         string
	CaConfig          conf.CA
	CaCertificate     *x509.Certificate
	CaPrivateKey      crypto.Signer
	ISD               string
	CertValidityHours int
}

func loadPKCS11Signer(KeyId []byte, pkcs11Conf crypto11.Config) (crypto.Signer, error) {

	ctx, err := crypto11.Configure(&pkcs11Conf)
	if err != nil {
		return nil, err
	}
	signer, err := ctx.FindKeyPair(nil, []byte(KeyId))
	if err != nil {
		return nil, err
	}

	return signer, nil
}

func NewSCIONCertificateAuthority(configDir string, isd string, conf conf.CA) *SCIONCertificateAuthority {
	return &SCIONCertificateAuthority{
		ConfigDir: configDir,
		ISD:       isd,
		CaConfig:  conf,
	}
}

func (ca *SCIONCertificateAuthority) LoadCA() error {

	caDir := filepath.Join(ca.ConfigDir, "crypto", "ca")

	if ca.CaConfig.PKCS11 != nil {
		pkcs11Conf := crypto11.Config{
			Path:       ca.CaConfig.PKCS11.ModulePath,
			TokenLabel: ca.CaConfig.PKCS11.TokenLabel,
			Pin:        ca.CaConfig.PKCS11.Pin,
		}
		signer, err := loadPKCS11Signer([]byte(ca.CaConfig.PKCS11.KeyLabel), pkcs11Conf)
		if err != nil {
			return fmt.Errorf("Failed to initialize signer: %v\n", err)
		}
		ca.CaPrivateKey = signer

		log.Printf("[CA] CA is using PKCS11 module %s for signing\n", ca.CaConfig.PKCS11.ModulePath)
	} else {
		// Load CA private key
		keyFile := filepath.Join(caDir, "cp-ca.key")

		caKeyPEM, err := os.ReadFile(keyFile)
		if err != nil {
			return fmt.Errorf("Failed to read CA private key: %v\n", err)
		}
		caKey, err := loadKey(caKeyPEM)
		if err != nil {
			return fmt.Errorf("Failed to load CA key: %v\n", err)
		}
		ca.CaPrivateKey = caKey

		log.Println("[CA] CA is using cp-ca.key for signing")
	}

	// Load CA certificate

	caCertFiles, err := fileops.ListFilesByPrefixAndSuffix(caDir, "ISD", ".crt")
	if err != nil {
		return fmt.Errorf("Failed to list CA certificate files: %v\n", err)
	}
	// TODO: We assume only one .crt file here
	certFile := caCertFiles[0]

	caCertPEM, err := os.ReadFile(certFile)
	if err != nil {
		return fmt.Errorf("Failed to read CA certificate: %v\n", err)

	}

	caCert, err := loadCert(caCertPEM)
	if err != nil {
		return fmt.Errorf("Failed to load CA cert and key: %v\n", err)
	}

	ca.CaCertificate = caCert
	return nil
}

func (ca *SCIONCertificateAuthority) IssueCertificateFromCSR(csrFile string, dstFile, isd string, as string) error {
	log.Println("[CA] Issuing certificate from CSR, csrFile ", csrFile)
	csrPEM, err := ioutil.ReadFile(csrFile)
	if err != nil {
		return fmt.Errorf("Failed to read CSR: %v\n", err)
	}
	csrPEM = formatPEMString(string(csrPEM))

	// Parse the CSR
	block, _ := pem.Decode(csrPEM)
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		return fmt.Errorf("failed to decode CSR PEM")
	}
	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse CSR: %v", err)
	}

	// Validate the CSR
	err = csr.CheckSignature()
	if err != nil {
		return fmt.Errorf("CSR signature verification failed: %v", err)
	}

	// Calculate the Subject Key Identifier
	/*pubKeyBytes, err := x509.MarshalPKIXPublicKey(csr.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %v", err)
	}*/
	subjectKeyID, err := cppki.SubjectKeyID(csr.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to calculate subject key ID: %v", err)
	}

	// Create a serial number for the certificate
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %v", err)
	}

	customOID := asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 55324, 1, 2, 1}

	// Create the leaf certificate template
	leafCertTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: csr.Subject.CommonName,
			ExtraNames: []pkix.AttributeTypeAndValue{
				{
					Type:  customOID,
					Value: fmt.Sprintf("%s-%s", isd, as),
				},
			},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, 0, ca.CertValidityHours), // Valid for X hours (default 72)
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageTimeStamping},
		SubjectKeyId: subjectKeyID[:],
	}

	// Set the public key from the CSR
	leafCertTemplate.PublicKey = csr.PublicKey

	// Sign the leaf certificate with the CA key
	leafCertDER, err := x509.CreateCertificate(rand.Reader, &leafCertTemplate, ca.CaCertificate, csr.PublicKey, ca.CaPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to create leaf certificate: %v", err)
	}

	// Encode the leaf certificate
	leafCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: leafCertDER})

	// Append the CA certificate to the leaf certificate
	fullChainPEM := append(leafCertPEM, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: ca.CaCertificate.Raw})...)

	// Save the full chain to a file
	err = os.WriteFile(dstFile, fullChainPEM, 0644)
	if err != nil {
		return fmt.Errorf("failed to write full chain certificate: %v", err)
	}

	return nil
}

// loadCert loads the certificate from PEM-encoded data
func loadCert(certPEM []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to decode certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %v", err)
	}

	return cert, nil
}

// loadKey loads the private key from PEM-encoded data
func loadKey(keyPEM []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(keyPEM)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode private key PEM")
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	// Assert that the parsed key is of type *ecdsa.PrivateKey
	key, ok := parsedKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an ECDSA private key")
	}

	return key, nil
}

func formatPEMString(pemStr string) []byte {
	// Remove existing newlines (if any) and extract the PEM content
	pemStr = strings.ReplaceAll(pemStr, "\n", "")

	// Ensure the BEGIN and END markers are preserved correctly
	beginMarker := "-----BEGIN CERTIFICATE REQUEST-----"
	endMarker := "-----END CERTIFICATE REQUEST-----"

	// Extract the actual PEM data
	pemData := strings.TrimPrefix(pemStr, beginMarker)
	pemData = strings.TrimSuffix(pemData, endMarker)
	pemData = strings.TrimSpace(pemData) // Remove extra spaces

	// Reformat by inserting a newline every 64 characters
	var formattedPem bytes.Buffer
	formattedPem.WriteString(beginMarker + "\n")

	for i := 0; i < len(pemData); i += 64 {
		end := i + 64
		if end > len(pemData) {
			end = len(pemData)
		}
		formattedPem.WriteString(pemData[i:end] + "\n")
	}

	formattedPem.WriteString(endMarker + "\n")

	return formattedPem.Bytes()
}
