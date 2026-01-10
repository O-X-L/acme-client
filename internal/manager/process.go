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
	"time"

	"git.oxl.at/acme-client/internal/acme"
	"git.oxl.at/acme-client/internal/u"
	"git.oxl.at/acme-client/pkg/config"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	acme_logger "github.com/go-acme/lego/v4/log"
	"github.com/go-acme/lego/v4/registration"
)

func processDomainBatch(grp config.Group, svc config.Service, batchID int, domains []string) (bool, error) {
	key := fmt.Sprintf("[Group: %d '%s' | Service: %d | Batch: %d]", grp.ID, grp.Name, svc.ID, batchID)
	certBaseName := fmt.Sprintf("%d_%d_%d", grp.ID, svc.ID, batchID)

	updateNeeded, reason := NeedsUpdate(certBaseName, domains)
	if !updateNeeded {
		u.Logf("%s skipping: cert is valid", key)
		return false, nil
	}

	u.Logf("%s updating: %s", key, reason)

	for attempt := uint(0); attempt <= config.Config.Retries; attempt++ {
		time.Sleep(time.Second * time.Duration(config.Config.CooldownSec))
		u.Logf("%s obtaining certificate...", key)
		err := obtainCert(certBaseName, svc, domains)
		if err == nil {
			return true, nil
		}
		u.LogWarningf("%s attempt %d failed: \"%v\"", key, attempt+1, err)
	}

	u.LogWarningf("%s could not process certificate after %d retries", key, config.Config.Retries)
	return false, fmt.Errorf("failed")
}

func processDomainsInBatches(
	grp config.Group, svc config.Service, domains []string,
	callback func(config.Group, config.Service, int, []string) (bool, error),
) (bool, int, error) {
	anyChanged := false
	var err error
	processedBatches := 0
	for batchID, batchDomains := range u.BuildBatches(domains, int(config.Config.MaxDomains)) {
		processedBatches += 1
		changed, err := callback(grp, svc, batchID+1, batchDomains)
		if changed {
			anyChanged = true
		}
		if err != nil {
			return anyChanged, processedBatches, err
		}
	}
	return anyChanged, processedBatches, err
}

func processService(grp config.Group, svc config.Service) (bool, error) {
	key := fmt.Sprintf("[Group: %d '%s' | Service: %d]", grp.ID, grp.Name, svc.ID)
	u.Logf("%s processing...", key)

	if len(svc.Domains) == 0 {
		u.Logf("%s skipping: service has no domains configured", key)
		return false, nil
	}

	providedDomainCount := len(svc.Domains)
	svc.Domains = u.RemoveDuplicates(svc.Domains)
	uniqueDomainCount := len(svc.Domains)
	if uniqueDomainCount != providedDomainCount {
		u.LogWarningf("%s has %d duplicate domains configured", key, providedDomainCount-uniqueDomainCount)
	}

	changed, _, err := processDomainsInBatches(grp, svc, svc.Domains, processDomainBatch)
	return changed, err
}

func obtainCert(name string, svc config.Service, domains []string) error {
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
		log.Printf("Executing hook: \"%s\"", config.Config.HookCmd)
		cmd := exec.Command("sh", "-c", config.Config.HookCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Hook command failed: %v", err)
		}
	}
}
