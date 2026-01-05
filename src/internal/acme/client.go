package acme

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns"
	"github.com/go-acme/lego/v4/providers/http/webroot"
)

const (
	// as lego-module does not seem to return a simple list
	//   copied from: https://github.com/go-acme/lego/tree/master/providers/dns
	//   todo: move to shared config
	PROVIDERS_DNS = `acmedns
active24
alidns
aliesa
allinkl
alwaysdata
anexia
arvancloud
auroradns
autodns
axelname
azion
azure
azuredns
baiducloud
beget
binarylane
bindman
bluecat
bookmyname
brandit
bunny
checkdomain
civo
clouddns
cloudflare
cloudns
cloudru
cloudxns
conoha
conohav3
constellix
corenetworks
cpanel
derak
desec
designate
digitalocean
directadmin
dnshomede
dnsimple
dnsmadeeasy
dnspod
dode
domeneshop
dreamhost
duckdns
dyn
dyndnsfree
dynu
easydns
edgecenter
edgedns
edgeone
efficientip
epik
exec
exoscale
f5xc
freemyip
gandi
gandiv5
gcloud
gcore
gigahostno
glesys
godaddy
googledomains
gravity
hetzner
hostingde
hostinger
hostingnl
hosttech
httpnet
httpreq
huaweicloud
hurricane
hyperone
ibmcloud
iij
iijdpf
infoblox
infomaniak
internal
internetbs
inwx
ionos
ionoscloud
ipv64
ispconfig
ispconfigddns
iwantmyname
joker
keyhelp
liara
lightsail
limacity
linode
liquidweb
loopia
luadns
mailinabox
manageengine
manual
metaname
metaregistrar
mijnhost
mittwald
myaddr
mydnsjp
mythicbeasts
namecheap
namedotcom
namesilo
nearlyfreespeech
neodigit
netcup
netlify
nicmanager
nicru
nifcloud
njalla
nodion
ns1
octenium
oraclecloud
otc
ovh
pdns
plesk
porkbun
rackspace
rainyun
rcodezero
regfish
regru
rfc2136
rimuhosting
route53
safedns
sakuracloud
scaleway
selectel
selectelv2
selfhostde
servercow
shellrent
simply
sonic
spaceship
stackpath
syse
technitium
tencentcloud
timewebcloud
transip
ultradns
uniteddomains
variomedia
vegadns
vercel
versio
vinyldns
virtualname
vkcloud
volcengine
vscale
vultr
webnames
webnamesca
websupport
wedos
westcn
yandex
yandex360
yandexcloud
zoneedit
zoneee
zonomi`
)

func NewACMEClientHttp(user *User, caURL, webrootPath string) (*lego.Client, error) {
	config := lego.NewConfig(user)
	config.CADirURL = caURL

	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}

	ps, err := webroot.NewHTTPProvider(webrootPath)
	if err != nil {
		return nil, err
	}
	client.Challenge.SetHTTP01Provider(ps)

	return client, nil
}

func NewACMEClientDns(user *User, providerName string, envVars map[string]string) (*lego.Client, error) {
	config := lego.NewConfig(user)

	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}

	for k, v := range envVars {
		os.Setenv(k, v)
	}
	cp, err := dns.NewDNSChallengeProviderByName(providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to load dns provider %s: %v", providerName, err)
	}
	client.Challenge.SetDNS01Provider(cp)

	return client, nil
}

func IsSupportedProvider(name string) bool {
	providerList := strings.Split(PROVIDERS_DNS, "\n")
	return slices.Contains(providerList, name)
}
