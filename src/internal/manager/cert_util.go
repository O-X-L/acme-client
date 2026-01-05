package manager

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"time"

	"git.oxl.at/acme-client/internal/config"
)

func NeedsUpdate(baseName string, desired []string) (bool, string) {
	crtPath := filepath.Join(config.PathCertsPublic, baseName+".crt")
	keyPath := filepath.Join(config.PathCertsPrivate, baseName+".key")

	crtData, err := os.ReadFile(crtPath)
	if err != nil {
		return true, "certificate file missing"
	}

	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return true, "private key file missing"
	}

	if *config.Config.CreateBundle {
		if _, err := os.Stat(filepath.Join(config.PathCertsBundlePublic, baseName+".crt")); err != nil {
			return true, "public-bundle file missing"
		}
		if _, err := os.Stat(filepath.Join(config.PathCertsBundlePrivate, baseName+".pem")); err != nil {
			return true, "private-bundle file missing"
		}
	}

	block, _ := pem.Decode(crtData)
	if block == nil {
		return true, "invalid certificate pem"
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return true, "certificate parse error"
	}

	if len(cert.RawTBSCertificate) == 0 {
		return true, "certificate chain/ca is empty"
	}

	keyBlock, _ := pem.Decode(keyData)
	if keyBlock == nil {
		return true, "invalid key pem"
	}
	privKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		privKey, err = x509.ParseECPrivateKey(keyBlock.Bytes)
	}

	if err == nil {
		if !publicKeyMatches(cert.PublicKey, privKey) {
			return true, "private key and public key mismatch"
		}
	}

	if time.Until(cert.NotAfter) < 0 {
		return true, "expired"
	}

	if time.Until(cert.NotAfter) < config.RENEWAL_DAYS*24*time.Hour {
		return true, "expiring soon"
	}

	sort.Strings(desired)
	actual := cert.DNSNames
	sort.Strings(actual)
	if !reflect.DeepEqual(desired, actual) {
		return true, "domain mismatch"
	}

	return false, ""
}

func publicKeyMatches(pub crypto.PublicKey, priv crypto.PrivateKey) bool {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return reflect.DeepEqual(pub, &k.PublicKey)
	case *ecdsa.PrivateKey:
		return reflect.DeepEqual(pub, &k.PublicKey)
	default:
		return false
	}
}

func SetOwnership(path string) {
	if config.Config.FileGroup == "" {
		return
	}
	g, err := user.LookupGroup(config.Config.FileGroup)
	if err != nil {
		return
	}
	gid, _ := strconv.Atoi(g.Gid)
	os.Chown(path, -1, gid)
}
