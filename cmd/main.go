package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"git.oxl.at/acme-client/internal/acme"
	"git.oxl.at/acme-client/internal/manager"
	"git.oxl.at/acme-client/internal/u"
	"git.oxl.at/acme-client/pkg/config"
)

func initCertDir() error {
	config.PathCertsPublic = filepath.Join(config.Config.PathCerts, config.DIR_CERTS_PUBLIC)
	config.PathCertsPrivate = filepath.Join(config.Config.PathCerts, config.DIR_CERTS_PRIVATE)
	config.PathCertsBundlePublic = filepath.Join(config.Config.PathCerts, config.DIR_BUNDLE_PUBLIC)
	config.PathCertsBundlePrivate = filepath.Join(config.Config.PathCerts, config.DIR_BUNDLE_PRIVATE)
	config.PathAccountCache = filepath.Join(config.Config.PathCerts, config.DIR_ACCOUNT)

	err := os.MkdirAll(config.PathCertsPublic, 0755)
	if err != nil {
		return fmt.Errorf("failed to create certificate-public directory: %s - %v", config.PathCertsPublic, err)
	}

	err = os.MkdirAll(config.PathCertsPrivate, 0750)
	if err != nil {
		return fmt.Errorf("failed to create certificate-private directory: %s - %v", config.PathCertsPrivate, err)
	}

	if *config.Config.CreateBundle {
		err = os.MkdirAll(config.PathCertsBundlePublic, 0755)
		if err != nil {
			return fmt.Errorf("failed to create certificate-public-bundle directory: %s - %v", config.PathCertsBundlePublic, err)
		}
		err = os.MkdirAll(config.PathCertsBundlePrivate, 0750)
		if err != nil {
			return fmt.Errorf("failed to create certificate-private-bundle directory: %s - %v", config.PathCertsBundlePrivate, err)
		}
	}

	err = os.MkdirAll(config.PathAccountCache, 0700)
	if err != nil {
		return fmt.Errorf("failed to create account-cache directory: %s - %v", config.PathAccountCache, err)
	}

	if config.Config.FileGroup != "" {
		manager.SetOwnership(config.PathCertsPublic)
		manager.SetOwnership(config.PathCertsPrivate)
		manager.SetOwnership(config.PathCertsBundlePublic)
		manager.SetOwnership(config.PathCertsBundlePrivate)
	}
	return nil
}

func main() {
	fmt.Printf("OXL ACME-Client | Version: %s | License: MIT | Repo: https://git.OXL.at/acme-client | © 2026 OXL IT Services\n", config.VERSION)

	var pathConfig string
	var showProviders bool
	flag.StringVar(&pathConfig, "c", "acme.yml", "Path to config file")
	flag.BoolVar(&showProviders, "show-providers", false, "Only show supported DNS-providers and exit")
	flag.BoolVar(&config.CheckMode, "check", false, "Only validate the config-file")
	flag.Parse()

	if showProviders {
		fmt.Printf("Supported DNS-Providers:\n\n%v\n", acme.PROVIDERS_DNS)
		os.Exit(0)
	}

	u.InitLogModes()

	cnf, err := config.LoadConfig(pathConfig)
	if err != nil {
		u.LogError(fmt.Sprintf("Failed to load config: %v", err))
		os.Exit(1)
	}
	config.Config = cnf

	err = config.ValidateConfig(cnf)
	if err != nil {
		u.LogError(fmt.Sprintf("Got invalid config: %v", err))
		os.Exit(1)
	}

	if !config.ValidateSchema(cnf) {
		os.Exit(1)
	}

	if config.CheckMode {
		u.Log("Config is valid")
		os.Exit(0)
	}

	err = initCertDir()
	if err != nil {
		u.LogError(fmt.Sprintf("%v", err))
		os.Exit(1)
	}

	if config.Config.PathWeb != "" {
		testFile := filepath.Join(config.Config.PathWeb, ".test")
		err = os.WriteFile(testFile, []byte(""), 0644)
		if err != nil {
			u.LogError(fmt.Sprintf("configured path_web is not writable: %v", err))
			os.Exit(1)
		}
		os.Remove(testFile)
	}

	config.RenewalDays = time.Duration(config.Config.RenewalDays) * 24 * time.Hour

	manager.Run()
}
