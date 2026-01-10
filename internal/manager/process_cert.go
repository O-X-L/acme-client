package manager

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"git.oxl.at/acme-client/internal/acme"
	"git.oxl.at/acme-client/internal/u"
	"git.oxl.at/acme-client/pkg/config"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	acme_logger "github.com/go-acme/lego/v4/log"
	"github.com/go-acme/lego/v4/registration"
)

func obtainCert(name string, svc config.Service, domains []string) error {
	if config.ModeCheck {
		return nil
	}

	user, err := getOrCreateACMEUser(svc.Provider)
	if err != nil {
		return err
	}

	var client *lego.Client
	if svc.ChallengeType == config.CHALLENGE_TYPE_DNS {
		client, err = acme.NewACMEClientDns(user, svc.Provider, svc.ProviderConfig)
	} else {
		client, err = acme.NewACMEClientHttp(user, svc.Provider, config.Config.PathWeb)
	}
	if err != nil {
		return err
	}

	if err := ensureACMEUserRegistration(client, user); err != nil {
		return err
	}

	res, err := client.Certificate.Obtain(certificate.ObtainRequest{
		Domains: domains,
		Bundle:  false,
	})

	if err != nil {
		return err
	}

	crtPath := filepath.Join(config.PathCertsPublic, name+".crt")
	keyPath := filepath.Join(config.PathCertsPrivate, name+".key")

	if err := os.WriteFile(crtPath, res.Certificate, config.Config.FileModeCert); err != nil {
		return err
	}
	if err := os.WriteFile(keyPath, res.PrivateKey, config.Config.FileModeKey); err != nil {
		return err
	}
	SetOwnership(crtPath)
	SetOwnership(keyPath)

	if *config.Config.CreateBundle {
		bundlePublic := append(res.Certificate, res.IssuerCertificate...)
		bundlePrivate := append(bundlePublic, res.PrivateKey...)
		bundlePublicPath := filepath.Join(config.PathCertsBundlePublic, name+".crt")
		bundlePrivatePath := filepath.Join(config.PathCertsBundlePrivate, name+".pem")
		os.WriteFile(bundlePublicPath, bundlePublic, config.Config.FileModeCert)
		os.WriteFile(bundlePrivatePath, bundlePrivate, config.Config.FileModeKey)
		SetOwnership(bundlePublicPath)
		SetOwnership(bundlePrivatePath)
	}

	return nil
}

func getOrCreateACMEUser(provider string) (*acme.User, error) {
	hash := sha256.Sum256([]byte(provider + config.Config.Email))
	hashStr := hex.EncodeToString(hash[:])
	keyPath := filepath.Join(config.PathAccountCache, fmt.Sprintf("account_%s.key", hashStr))

	var privKey crypto.PrivateKey
	data, err := os.ReadFile(keyPath)
	if err == nil {
		block, _ := pem.Decode(data)
		if block != nil {
			privKey, _ = x509.ParseECPrivateKey(block.Bytes)
		}
	}

	if privKey == nil {
		u.Logf("Generating new account key for %s", config.Config.Email)
		newKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}
		privKey = newKey

		keyBytes, _ := x509.MarshalECPrivateKey(newKey)
		pemKey := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})
		_ = os.WriteFile(keyPath, pemKey, 0600)
	}

	return &acme.User{Email: config.Config.Email, Key: privKey}, nil
}

func ensureACMEUserRegistration(client *lego.Client, user *acme.User) error {
	reg, err := client.Registration.ResolveAccountByKey()
	if err != nil {
		u.Logf("Registering new ACME account for %s", user.Email)
		reg, err = client.Registration.Register(registration.RegisterOptions{
			TermsOfServiceAgreed: true,
		})
		if err != nil {
			return fmt.Errorf("failed to register ACME user: %v", err)
		}
	}
	user.Registration = reg
	return nil
}

func Run() {
	acme_logger.Logger = log.New(os.Stdout, "[ACME] ", log.LstdFlags)

	anyChanged := false

	for _, grp := range config.Config.Groups {
		for _, svc := range grp.Services {
			changed, err := processService(grp, svc)
			if err != nil {
				u.LogErrorf(
					"[%s] failed to obtain service-certs %d (%d domains) via challenge '%s' of provider '%s'",
					grp.Name, svc.ID, len(svc.Domains), svc.ChallengeType, svc.Provider,
				)

			} else if changed {
				anyChanged = true
			}
		}
	}

	if anyChanged && strings.TrimSpace(config.Config.HookCmd) != "" {
		u.Logf("Executing hook: \"%s\"", config.Config.HookCmd)
		cmd := exec.Command("sh", "-c", config.Config.HookCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			u.LogErrorf("Hook command failed: %v", err)
		}
	}
}
