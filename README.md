# Standalone ACME Client

<p align="center">
    <a title="Support this Project (Donate, Support-Licenses)" href="https://shop.oxl.app/collections/open-source">
        <img src="https://files.oxl.at/img/badge-oss-support.svg" alt="Support Badge (Donate, Support-Licenses)"/>
    </a>
</p>

----

[![Lint](https://github.com/O-X-L/acme-client/actions/workflows/lint.yml/badge.svg?branch=latest)](https://github.com/O-X-L/acme-client/actions/workflows/lint.yml)
[![Unit Test](https://github.com/O-X-L/acme-client/actions/workflows/unit_test.yml/badge.svg?branch=latest)](https://github.com/O-X-L/acme-client/actions/workflows/unit_test.yml)

This ACME-client is based on the awesome [go-acme/lego](https://github.com/go-acme/lego) library. ❤️

This client enables you to supply a simple configuration-file that will request certificates and save them to your filesystem - similar to how [dehydrated](https://github.com/dehydrated-io/dehydrated) does.

----

## Usage

### Install

1. Get the binary

  * Download pre-compiled binary from the [Releases](https://github.com/O-X-L/acme-client/releases)

  * Or build it yourself:

    * [Download & Install Go](https://go.dev/doc/install)
    * Build: 

      ```bash
      mkdir $REPO/build
      cd $REPO/src
      go mod tidy  # download dependencies
      GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o "../build/acme" ./cmd/main.go
      ```

2. Transfer the binary to your server
3. Prepare for challenges:

  * To use the `http-01` challenge - configure a web-root directory and a web-server that serves its content via plain HTTP
  * To use the `dns-01` challenge - create an account at a supported provider

----

### Config

It expects a YAML config-file in this format:

```yaml
---

email: 'test@waf.alpenmesh.com'
retries: 1  # retries per configured certificate if a validation error occurred
cooldown_sec: 1  # seconds to wait between requests/retries
path_web: '/var/www/acme'  # web-root-directory - has to contain '<path>/.well-known/acme-challenge/' and be writable for the service-user
path_certs: '/var/local/acme'
create_bundle: true  # optionally create certificate bundles (public: cert+ca, private: cert+ca+pk)
file_mode_cert: 0644  # default: 0640
file_mode_key: 0640  # default: 0600
file_group: 'ssl-cert'  # default: primary group of service-user
hook_cmd: 'echo "DONE"'  # hook command to be ran after all certificates were processed AND something changed

apps:
  - name: "App #1"
    id: 1
    certs:
      - id: 1
        challenge_type: "http-01"
        provider: "https://acme-staging-v02.api.letsencrypt.org/directory"
        domains:
          - 'aaa.waf.alpenmesh.com'
          - 'bbb.waf.alpenmesh.com'
          - 'ccc.waf.alpenmesh.com'

  - name: "App #2"
    id: 2
    certs:
      - id: 1
        challenge_type: "dns-01"
        provider: "cloudflare"
        provider_config:  # env-vars to pass the go-acme/lego execution
          CLOUDFLARE_API_KEY: '...'
        domains:
          - '*.alpenmesh.net'
```

For DNS-Provider config see: [go-acme/lego documentation](https://go-acme.github.io/lego/dns/index.html)

See also: [Examples](https://github.com/O-X-L/acme-client/blob/latest/examples/)

----

### Run

```bash
# move to permanent directory after upload
mv ./acme /usr/local/bin/acme

# to check the supported DNS-providers
/usr/local/bin/acme -show-providers

# run
/usr/local/bin/acme -path-cnf /etc/acme/acme.yml

# to connect over proxy
export HTTPS_PROXY=http://test-proxy.waf.alpenmesh.com:3128
/usr/local/bin/acme -path-cnf /etc/acme/acme.yml
```

----

### Result

```bash
root@srv:/var/local/acme# tree
├── account  # account cache
│   ├── account_09ff80dda58a752729e0506d726ba47590ff1413129666b581a2eee1fa01449b.key
│   └── account_3bfae30343be0ae9c6709cc568ac155d2c3cb562fdf487f6e858b8cd0006cd27.key
├── bundle_certs  # if 'create_bundle: true' | public-key bundles
│   ├── app_1_1.crt
│   └── app_2_1.crt
├── bundle_private  # if 'create_bundle: true' | bundles including private-key
│   ├── app_1_1.pem
│   └── app_2_1.pem
├── certs  # public-keys
│   ├── app_1_1.crt
│   └── app_2_1.crt
└── private  # private-keys
    ├── app_1_1.key
    └── app_2_1.key

root@srv:/var/local/acme# ls -l
drwx------ 2 acme acme     4096 Jan  5 23:32 account
drwxr-xr-x 2 acme ssl-cert 4096 Jan  5 23:34 bundle_certs
drwxr-x--- 2 acme ssl-cert 4096 Jan  5 23:34 bundle_private
drwxr-xr-x 2 acme ssl-cert 4096 Jan  5 23:34 certs
drwxr-x--- 2 acme ssl-cert 4096 Jan  5 23:34 private

root@srv:/var/local/acme# ls -l */*
-rw------- 1 acme acme      227 Jan  5 23:32 account/account_09ff80dda58a752729e0506d726ba47590ff1413129666b581a2eee1fa01449b.key
-rw------- 1 acme acme      227 Jan  5 23:32 account/account_3bfae30343be0ae9c6709cc568ac155d2c3cb562fdf487f6e858b8cd0006cd27.key
-rw-r--r-- 1 acme ssl-cert 3831 Jan  5 23:33 bundle_certs/app_1_1.crt
-rw-r--r-- 1 acme ssl-cert 3831 Jan  5 23:33 bundle_certs/app_2_1.crt
-rw-r----- 1 acme ssl-cert 5510 Jan  5 23:33 bundle_private/app_1_1.pem
-rw-r----- 1 acme ssl-cert 5506 Jan  5 23:33 bundle_private/app_2_1.pem
-rw-r--r-- 1 acme ssl-cert 1935 Jan  5 23:33 certs/app_1_1.crt
-rw-r--r-- 1 acme ssl-cert 1935 Jan  5 23:33 certs/app_2_1.crt
-rw-r----- 1 acme ssl-cert 1679 Jan  5 23:33 private/app_1_1.key
-rw-r----- 1 acme ssl-cert 1675 Jan  5 23:33 private/app_2_1.key
```

#### Output / Logs

Note: *This output is not from the exact config-example above.*

First run:

<details>

```
OXL ACME-Client | Version: 1.0.0 | License: MIT | Repo: https://git.OXL.at/acme-client | © 2026 OXL IT Services
2026/01/05 23:04:30 [INFO] [App: 1 'App #1' | Cert: 1] processing...
2026/01/05 23:04:30 [INFO] [App: 1 'App #1' | Cert: 1] updating: certificate file missing
2026/01/05 23:04:31 [INFO] [App: 1 'App #1' | Cert: 1] obtaining certificate...
2026/01/05 23:04:31 [INFO] Generating new account key for test@waf.alpenmesh.com
2026/01/05 23:04:32 [INFO] acme: Trying to resolve account by key
2026/01/05 23:04:32 [INFO] Registering new ACME account for test@waf.alpenmesh.com
2026/01/05 23:04:32 [INFO] acme: Registering account for test@waf.alpenmesh.com
2026/01/05 23:04:32 [INFO] [aaa.waf.alpenmesh.com, bbb.waf.alpenmesh.com, ccc.waf.alpenmesh.com] acme: Obtaining bundled SAN certificate
2026/01/05 23:04:33 [INFO] [aaa.waf.alpenmesh.com] AuthURL: https://acme-staging-v02.api.letsencrypt.org/acme/authz/255623553/21042904583
2026/01/05 23:04:33 [INFO] [bbb.waf.alpenmesh.com] AuthURL: https://acme-staging-v02.api.letsencrypt.org/acme/authz/255623553/21042904593
2026/01/05 23:04:33 [INFO] [ccc.waf.alpenmesh.com] AuthURL: https://acme-staging-v02.api.letsencrypt.org/acme/authz/255623553/21042904603
2026/01/05 23:04:33 [INFO] [aaa.waf.alpenmesh.com] acme: Could not find solver for: tls-alpn-01
2026/01/05 23:04:33 [INFO] [aaa.waf.alpenmesh.com] acme: use http-01 solver
2026/01/05 23:04:33 [INFO] [bbb.waf.alpenmesh.com] acme: Could not find solver for: tls-alpn-01
2026/01/05 23:04:33 [INFO] [bbb.waf.alpenmesh.com] acme: use http-01 solver
2026/01/05 23:04:33 [INFO] [ccc.waf.alpenmesh.com] acme: Could not find solver for: tls-alpn-01
2026/01/05 23:04:33 [INFO] [ccc.waf.alpenmesh.com] acme: use http-01 solver
2026/01/05 23:04:33 [INFO] [aaa.waf.alpenmesh.com] acme: Trying to solve HTTP-01
2026/01/05 23:04:48 [INFO] [aaa.waf.alpenmesh.com] The server validated our request
2026/01/05 23:04:48 [INFO] [bbb.waf.alpenmesh.com] acme: Trying to solve HTTP-01
2026/01/05 23:04:51 [INFO] [bbb.waf.alpenmesh.com] The server validated our request
2026/01/05 23:04:51 [INFO] [ccc.waf.alpenmesh.com] acme: Trying to solve HTTP-01
2026/01/05 23:04:56 [INFO] [ccc.waf.alpenmesh.com] The server validated our request
2026/01/05 23:04:56 [INFO] [aaa.waf.alpenmesh.com, bbb.waf.alpenmesh.com, ccc.waf.alpenmesh.com] acme: Validations succeeded; requesting certificates
2026/01/05 23:04:57 [INFO] Wait for certificate [timeout: 30s, interval: 500ms]
2026/01/05 23:05:00 [INFO] [aaa.waf.alpenmesh.com] Server responded with a certificate.
2026/01/05 23:05:00 [INFO] [App: 1 'App #1' | Cert: 2] processing...
2026/01/05 23:05:00 [WARN] [App: 1 'App #1' | Cert: 2] has 1 duplicate domains configured
2026/01/05 23:05:00 [INFO] [App: 1 'App #1' | Cert: 2] updating: certificate file missing
2026/01/05 23:05:01 [INFO] [App: 1 'App #1' | Cert: 2] obtaining certificate...
2026/01/05 23:05:01 [INFO] acme: Trying to resolve account by key
2026/01/05 23:05:01 [INFO] [ddd.waf.alpenmesh.com, eee.waf.alpenmesh.com] acme: Obtaining bundled SAN certificate
2026/01/05 23:05:02 [INFO] [ddd.waf.alpenmesh.com] AuthURL: https://acme-staging-v02.api.letsencrypt.org/acme/authz/255623553/21042910063
2026/01/05 23:05:02 [INFO] [eee.waf.alpenmesh.com] AuthURL: https://acme-staging-v02.api.letsencrypt.org/acme/authz/255623553/21042910073
2026/01/05 23:05:02 [INFO] [ddd.waf.alpenmesh.com] acme: Could not find solver for: tls-alpn-01
2026/01/05 23:05:02 [INFO] [ddd.waf.alpenmesh.com] acme: use http-01 solver
2026/01/05 23:05:02 [INFO] [eee.waf.alpenmesh.com] acme: Could not find solver for: tls-alpn-01
2026/01/05 23:05:02 [INFO] [eee.waf.alpenmesh.com] acme: use http-01 solver
2026/01/05 23:05:02 [INFO] [ddd.waf.alpenmesh.com] acme: Trying to solve HTTP-01
2026/01/05 23:05:09 [INFO] [ddd.waf.alpenmesh.com] The server validated our request
2026/01/05 23:05:09 [INFO] [eee.waf.alpenmesh.com] acme: Trying to solve HTTP-01
2026/01/05 23:05:14 [INFO] Skipping deactivating of valid auth: https://acme-staging-v02.api.letsencrypt.org/acme/authz/255623553/21042910063
2026/01/05 23:05:14 [INFO] Deactivating auth: https://acme-staging-v02.api.letsencrypt.org/acme/authz/255623553/21042910073
2026/01/05 23:05:14 [WARN] [App: 1 'App #1' | Cert: 2] attempt 1 failed: "error: one or more domains had a problem:
[eee.waf.alpenmesh.com] invalid authorization: acme: error: 400 :: urn:ietf:params:acme:error:dns :: While processing CAA for eee.waf.alpenmesh.com: DNS problem: SERVFAIL looking up CAA for com - the domain's nameservers may be malfunctioning
"
2026/01/05 23:05:15 [INFO] [App: 1 'App #1' | Cert: 2] obtaining certificate...
2026/01/05 23:05:15 [INFO] acme: Trying to resolve account by key
2026/01/05 23:05:16 [INFO] [ddd.waf.alpenmesh.com, eee.waf.alpenmesh.com] acme: Obtaining bundled SAN certificate
2026/01/05 23:05:16 [INFO] [ddd.waf.alpenmesh.com] AuthURL: https://acme-staging-v02.api.letsencrypt.org/acme/authz/255623553/21042910063
2026/01/05 23:05:16 [INFO] [eee.waf.alpenmesh.com] AuthURL: https://acme-staging-v02.api.letsencrypt.org/acme/authz/255623553/21042912623
2026/01/05 23:05:16 [INFO] [ddd.waf.alpenmesh.com] acme: authorization already valid; skipping challenge
2026/01/05 23:05:16 [INFO] [eee.waf.alpenmesh.com] acme: Could not find solver for: tls-alpn-01
2026/01/05 23:05:16 [INFO] [eee.waf.alpenmesh.com] acme: use http-01 solver
2026/01/05 23:05:16 [INFO] [eee.waf.alpenmesh.com] acme: Trying to solve HTTP-01
2026/01/05 23:05:23 [INFO] [eee.waf.alpenmesh.com] The server validated our request
2026/01/05 23:05:23 [INFO] [ddd.waf.alpenmesh.com, eee.waf.alpenmesh.com] acme: Validations succeeded; requesting certificates
2026/01/05 23:05:24 [INFO] Wait for certificate [timeout: 30s, interval: 500ms]
2026/01/05 23:05:24 [INFO] [ddd.waf.alpenmesh.com] Server responded with a certificate.
2026/01/05 23:05:24 [INFO] [App: 2 'App #2' | Cert: 1] processing...
2026/01/05 23:05:24 [INFO] [App: 2 'App #2' | Cert: 1] updating: certificate file missing
2026/01/05 23:05:25 [INFO] [App: 2 'App #2' | Cert: 1] obtaining certificate...
2026/01/05 23:05:26 [INFO] acme: Trying to resolve account by key
2026/01/05 23:05:26 [INFO] [fff.waf.alpenmesh.com, ggg.waf.alpenmesh.com, hhh.waf.alpenmesh.com] acme: Obtaining bundled SAN certificate
2026/01/05 23:05:27 [INFO] [fff.waf.alpenmesh.com] AuthURL: https://acme-staging-v02.api.letsencrypt.org/acme/authz/255623553/21042914213
2026/01/05 23:05:27 [INFO] [ggg.waf.alpenmesh.com] AuthURL: https://acme-staging-v02.api.letsencrypt.org/acme/authz/255623553/21042914223
2026/01/05 23:05:27 [INFO] [hhh.waf.alpenmesh.com] AuthURL: https://acme-staging-v02.api.letsencrypt.org/acme/authz/255623553/21042914233
2026/01/05 23:05:27 [INFO] [fff.waf.alpenmesh.com] acme: Could not find solver for: tls-alpn-01
2026/01/05 23:05:27 [INFO] [fff.waf.alpenmesh.com] acme: use http-01 solver
2026/01/05 23:05:27 [INFO] [ggg.waf.alpenmesh.com] acme: Could not find solver for: tls-alpn-01
2026/01/05 23:05:27 [INFO] [ggg.waf.alpenmesh.com] acme: use http-01 solver
2026/01/05 23:05:27 [INFO] [hhh.waf.alpenmesh.com] acme: Could not find solver for: tls-alpn-01
2026/01/05 23:05:27 [INFO] [hhh.waf.alpenmesh.com] acme: use http-01 solver
2026/01/05 23:05:27 [INFO] [fff.waf.alpenmesh.com] acme: Trying to solve HTTP-01
2026/01/05 23:05:32 [INFO] [fff.waf.alpenmesh.com] The server validated our request
2026/01/05 23:05:32 [INFO] [ggg.waf.alpenmesh.com] acme: Trying to solve HTTP-01
2026/01/05 23:05:40 [INFO] [ggg.waf.alpenmesh.com] The server validated our request
2026/01/05 23:05:40 [INFO] [hhh.waf.alpenmesh.com] acme: Trying to solve HTTP-01
2026/01/05 23:05:53 [INFO] [hhh.waf.alpenmesh.com] The server validated our request
2026/01/05 23:05:53 [INFO] [fff.waf.alpenmesh.com, ggg.waf.alpenmesh.com, hhh.waf.alpenmesh.com] acme: Validations succeeded; requesting certificates
2026/01/05 23:05:53 [INFO] Wait for certificate [timeout: 30s, interval: 500ms]
2026/01/05 23:05:54 [INFO] [fff.waf.alpenmesh.com] Server responded with a certificate.
2026/01/05 23:05:54 Executing hook: "echo "DONE""
DONE
```

</details>

Second run:

<details>

```
OXL ACME-Client | Version: 1.0.0 | License: MIT | Repo: https://git.OXL.at/acme-client | © 2026 OXL IT Services
2026/01/05 23:12:41 [INFO] [App: 1 'App #1' | Cert: 1] processing...
2026/01/05 23:12:41 [INFO] [App: 1 'App #1' | Cert: 1] skipping: cert is valid
2026/01/05 23:12:41 [INFO] [App: 1 'App #1' | Cert: 2] processing...
2026/01/05 23:12:41 [WARN] [App: 1 'App #1' | Cert: 2] has 1 duplicate domains configured
2026/01/05 23:12:41 [INFO] [App: 1 'App #1' | Cert: 2] skipping: cert is valid
2026/01/05 23:12:41 [INFO] [App: 2 'App #2' | Cert: 1] processing...
2026/01/05 23:12:41 [INFO] [App: 2 'App #2' | Cert: 1] skipping: cert is valid
```

</details>
