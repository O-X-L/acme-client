package manager

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"git.oxl.at/acme-client/internal/u"
	"git.oxl.at/acme-client/pkg/config"
)

// dummy certs created through easyrsa: https://gist.github.com/superstes/5eda95dcf4016e9376c6e8b60f8b1f6b
const (
	rsaCaCert = `-----BEGIN CERTIFICATE-----
MIIDeDCCAmCgAwIBAgIUG9OaWB9UHsYlZT++vrtgYSScfU8wDQYJKoZIhvcNAQEL
BQAwFjEUMBIGA1UEAwwLT1hMIFRFU1QgQ0EwHhcNMjYwMTAxMTczODMzWhcNNDUw
OTE4MTczODMzWjAWMRQwEgYDVQQDDAtPWEwgVEVTVCBDQTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBAKwemXgGvI7CjLUkZldGCV0K3hb+Kd8qlVvTtjzE
ZsMHYYX1ke5WOvVIq3lfDJvwBEPKZ3ij0iT+qAxeURvw5j7GhxTZMfn7iLil1nIg
4beUDLaNqeMPlrw6Dm6aTp2eXyimdPNgQ9mIieNkCCgfRYaQK+THEojYH2GbyHCe
8g57MrcsIu5eHJ2lS25yeSFxPrpk86kQmnj8s+zwDb3VuCJXLh19aAAB6wJMR6YL
YtAUn5Iz6o4N5eCNcCnNiRPqz7Pp1Bjl1wE+ltNnmFCIryAKIJGvbNKTtkcqWmh1
fr21YQu+AvU1VMjAZ0Z9exRn1yD3JrwPCVhrd3DKI9GRVgECAwEAAaOBvTCBujAM
BgNVHRMEBTADAQH/MB0GA1UdDgQWBBRJFsfS3qM8jaZYPq6jULm6sRan8TBRBgNV
HSMESjBIgBRJFsfS3qM8jaZYPq6jULm6sRan8aEapBgwFjEUMBIGA1UEAwwLT1hM
IFRFU1QgQ0GCFBvTmlgfVB7GJWU/vr67YGEknH1PMAsGA1UdDwQEAwIBBjArBgNV
HR8EJDAiMCCgHqAchhpodHRwczovL3Rlc3Qub3hsLmF0L2NhLmNybDANBgkqhkiG
9w0BAQsFAAOCAQEAYuzrH7aZ2BzBsWPKvWuf66ciXyjtPL0luXOouJ/UqY7pd0IW
Pn4giOUxSM6p3J3cOGmJ+q/gUQ5LBx8w+yayz512CRbJUMo0xPU3sDOInq9z01lR
l8z5PybIPHYwk/gLiucEQPP7cXwpKTeYbBNlMAq4VyTzy4/bwrOFF8alVKwZ6YuG
Gm080qzFv44B5waFk/WefId/sKvxdgOe4teY8nZwTZo0JRDbAGjYvaI7krsKkBq+
S9KCRz9HRaxQn/0lFh3xz4Xc0uQCpFonl2NmhLCw4wnixpt6kN+9tELcJfIpI1w3
gntP1aGEuFotJmWq/xLwtWKN3AJwXyl83na1Pw==
-----END CERTIFICATE-----`

	rsaSubCaCert = `-----BEGIN CERTIFICATE-----
MIIDeDCCAmCgAwIBAgIQHZzFIfoXJge4FqnJSO8ltzANBgkqhkiG9w0BAQsFADAW
MRQwEgYDVQQDDAtPWEwgVEVTVCBDQTAeFw0yNjAxMDExNzM4NDFaFw00NTA5MTgx
NzM4NDFaMBoxGDAWBgNVBAMMD09YTCBURVNUIFN1Yi1DQTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBAOdvi43vupVLLltiyGVZJuQJLLATIay5BPjqbKpE
XOdM4NeCAkTSZ0lK6l7Ys5ILFBIrZfsJkjvSnA5tZ7Wask2Lask2LkwDYl4XEgw4X
EAOQjhvJgDw2d8mjlDFau38j4/DnsZQ0bqCnYi5bkJz0HsXCKhHVbMpj13HHeFdK
OR5vr1KosxFGpA3Em8tbNbakf1cjVDUHv4qtB1tffVj1pKhL5rAOEbTH9sh0gAyY
uNkrfG05ufrVtYmfsUeyV8+KvrWWeEISJudZIh9FBu+UiV5uwe/I4O2afFSX3+Lu
je3JV8nAZmA17iiPzgNuwAvesMVYOjytUYiXP6MstWWNM7wNQukCAwEAAaOBvTCB
ujAMBgNVHRMEBTADAQH/MB0GA1UdDgQWBBRwDEm6aRU5xOIiy9PLmmG8/08BbTBR
BgNVHSMESjBIgBRJFsfS3qM8jaZYPq6jULm6sRan8aEapBgwFjEUMBIGA1UEAwwL
T1hMIFRFU1QgQ0GCFBvTmlgfVB7GJWU/vr67YGEknH1PMAsGA1UdDwQEAwIBBjAr
BgNVHR8EJDAiMCCgHqAchhpodHRwczovL3Rlc3Qub3hsLmF0L2NhLmNybDANBgkq
hkiG9w0BAQsFAAOCAQEAYb/bSqV6TXNP9hLbxlesgWVRMPMOIl/u5DFQikFe0RnS
bZUl vwo2MeTPSEV6MDQCqiKu8MLetDfv9u/BT8pVZq3FJJZ9sKVfawNjGBkMHnqvi
Fss O8kdPwjMGkHyQgzXzAHVjQ/L5gIdyKKRuE5bMc9bv2rz/4RYpRzG1T+Dwi+8
B2wq v9M+lQ/kcU4e3Fd/QNFv8XCJM65ye/6jVDgzrXSsufgd9PsOw6q/4FGd4kw
1FHuP v/i0qFpv9RSp+QjFis27d/dVrJj+JsTKfXb+zWjA2WaP/rzOigGe8wyqug
juCVMv 2D3lWHeqJxIFwf7I9f7GwMT7/ehAKLsPq0P4gw==
-----END CERTIFICATE-----`

	rsaApp01Cert = `-----BEGIN CERTIFICATE-----
MIIDszCCApugAwIBAgIQRjWept5cIha83iCH/4re0TANBgkqhkiG9w0BAQsFADAa
MRgwFgYDVQQDDA9PWEwgVEVTVCBTdWItQ0EwHhcNMjYwMTAxMTgwMzMzWhcNNDUw
OTE4MTgwMzMzWjAZMRcwFQYDVQQDDA5PWEwgVEVTVCBBcHAwMTCCASIwDQYJKoZI
hvcNAQEBBQADggEPADCCAQoCggEBAMz+rV2aSHJnTqX/4zfieLBRaXR9GiTq0IvY
JRiBBmRsGETef1/d67Xatt7LELz1dd7ysBC6hZcA90OyDe3LlPKhZ5V/FuYPMxJe
qcTjwGEE6TKe12zeVlPxEc2Fig92hu80cDS4OakHJQiR2ym/pIhk1Ij2JU0ZKsTv
hCa7N8+FIBrzgPvnKBTAegERQanMai2W5Dq2zSCnlg2kwEXUwixoEWUjE2IbeYcN
ywZ1VQ4l/On1nlETVHzAEQQKvi40CwxTxbXjQcRjEn8L3dqszwGhz8/EN1una8OX
AQ4ITekSuaIg6NcQhmWIDew4U/DOe1sPdyViWHoZdP8Sq8NwGU0CAwEAAaOB9TCB
8jAJBgNVHRMEAjAAMB0GA1UdDgQWBBQhJd7Hv9Z23rSBvVMwFdMr5b2I8DBOBgNV
HSMERzBFgBTgjbA3cXDd6uUwqNGWUXbIIyQrzaEapBgwFjEUMBIGA1UEAwwLT1hM
IFRFU1QgQ0GCEQDg3piHNauuYF7fthiARqcTMBMGA1UdJQQMMAoGCCsGAQUFBwMB
MAsGA1UdDwQEAwIFoDArBgNVHR8EJDAiMCCgHqAchhpodHRwczovL3Rlc3Qub3hs
LmF0L2NhLmNybDAnBgNVHREEIDAeggt0ZXN0Lm94bC5hdIIPd3d3LnRlc3Qub3hs
LmF0MA0GCSqGSIb3DQEBCwUAA4IBAQBp0cpGBygtz8ET21smC6X3rKh5WFShYVJK
LjQYYO/P7hQa6F27/0gdbg3ckcbIJeMq9NumgNyOkMnSK8zG6yRojSm7H5xr/RBj
NhOu2EHs1fhUdsOPoP860NlxuV23xd+fFY0dDLuIVZAcwzt/QQGEm4MOmQBwZEwo
CaEjUeMcFOFjytVZEoYbmunNhgWe2T/wXhRY+Bh4BKGtT5XqyorwJF6s+krkGI7e
H9cg4ZcjZXnRlvmPUmRfOPHbdr30eT9qyDBaLrH9pZoO3n7W44ciCSQu5RNM1ekx
Vbi33hpyUMoJ4wMw2Z2LMPYgLiPq9wev+2J52PiEjt5awUNskwps
-----END CERTIFICATE-----`

	rsaApp01Key = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDM/q1dmkhyZ06l
/+M34niwUWl0fRok6tCL2CUYgQZkbBhE3n9f3eu12rbeyxC89XXe8rAQuoWXAPdD
sg3ty5TyoWeVfxbmDzMSXqnE48BhBOkyntds3lZT8RHNhYoPdobvNHA0uDmpByUI
kdspv6SIZNSI9iVNGSrE74QmuzfPhSAa84D75ygUwHoBEUGpzGotluQ6ts0gp5YN
pMBF1MIsaBFlIxNiG3mHDcsGdVUOJfzp9Z5RE1R8wBEECr4uNAsMU8W140HEYxJ/
C93arM8Boc/PxDdbp2vDlwEOCE3pErmiIOjXEIZliA3sOFPwzntbD3clYlh6GXT/
EqvDcBlNAgMBAAECggEAAn2BWbyb14M3IIKpaaXiKEP+X17xI8XxNLVVWJOKf2JC
x/DZpHRbvnOrgJXRkNOIiHF/HIFG1Vt0yGwnLIc2AGl3RM6bHHhxb1cC9oAHDwap
eUnn4CNJT3VJM05z8Wb2tMK+i3p/zwArt/RXdnZqwX32Qy8hNyVqtkYDJTF8w8Kg
heoos0f66sA379D+gik1V0jtkGuK+6GvuwKwn4TQoRmhbpnKG6tBBLTRvRhY9Py3
d3uHdw7kys2mBsxGg6DvKsiDbVwMsLgklYqbQdlXpFczhBfL3eg60jDGgEgityPj
+9Nb6qndPlLZTgjmFVpZDA3bGGLMBA+cmfDUXGDCkQKBgQD6eNkcvwGzoATvSJqe
fcQNK1auqgeIptG0OdhWtODonghTrvk3392nzDnPENQ0vVoEZFp/B4GmD0C6vAxR
ror03hlnF77y50S7ITu0PojL3WBELj+khjXrTLgPDJKdZxseM8Ix+4OKNZxKll/3
cPzNxBT0J9E/HoXBKJqmfLlgCQKBgQDRhOKyWcNwKF0xhbHjgXs+qwYDSn5DCAJ9
pORW90jllU5d6q5I5Ns/M72FvkPupR6NOE2SXFKwwusCebVNPI9p2JBU+5b0czwt
ZyohbLV1nTXa183BIWLxpbpm7ibdttBht0gRt+KMUGW4ym4qwMe9+gTTkGet9fvR
n32q8Nl4JQKBgQC277DTOAaMJSG1irezbnPUkoS7CWB7RCwBkAYcPfvOqi22wSMw
1gbUWWsXe1kiM/IdJxaZlOfyW02RlWsB9ZN0CQtQqp1CV/txGXK70Lik/UkkQwsk
pQmYk+4Sv8INyJKb2n3Jd8O5HDLMn0v0M8fQmZgcQ0Cm8yoJzBg35PyX0QKBgQCT
QvJPZvYN6+DswMpyXHyyZGR6ha5PEN8nTnFLis1KyHFnY16ST4CmYIhx63Q11Qqv
OYaUO53HLYkemUrL+afXPmbbxGWqdSAzzVH4Yj78Zr4Gji3I891meRSV6geZSDgY
pkjaY0OxWYTVqDpchFkf9w3TYajtxXn0MUSTlGtVnQKBgGx2+NUlGCO0bAGUiWwo
WAHZMH6D8LKSYtKWNNB4ugEAHKkG6W0vQj9fuXi0y08feEyixwsvkhGKmqc/8esI
91gTa8EI9G0f+My4/bLYViLNCM7oGa6h5SKkpJTs0BYrZYf1THjIrM02hEK/kEQY
H7+F6J0VyZPKZzHVUCjnHI88
-----END PRIVATE KEY-----`

	ecCaCert = `-----BEGIN CERTIFICATE-----
MIICKjCCAa+gAwIBAgIUEJ+v+8rbSEywHrUsm5iZ6GBobA8wCgYIKoZIzj0EAwIw
FjEUMBIGA1UEAwwLT1hMIFRFU1QgQ0EwHhcNMjYwMTAxMTgyNjI0WhcNNDUwOTE4
MTgyNjI0WjAWMRQwEgYDVQQDDAtPWEwgVEVTVCBDQTB2MBAGByqGSM49AgEGBSuB
BAAiA2IABPr2tMX5CQ01dLhSdFBa8EdYdhHjperr5pwZbGOrLYDqnbcRkspFzzeK
QyHHXv1DvwKu2/UWIl/BkD33MKjEsjgjSPUn8aahuG+pKM2ysKWCZqYWcDEqoU+Z
J5/AtSGWgKOBvTCBujAMBgNVHRMEBTADAQH/MB0GA1UdDgQWBBQOqDMkf+R6lJo0
9T6fi/15YunYtDBRBgNVHSMESjBIgBQOqDMkf+R6lJo09T6fi/15YunYtKEapBgw
FjEUMBIGA1UEAwwLT1hMIFRFU1QgQ0GCFBCfr/vK20hMsB61LJuYmehgaGwPMAsG
A1UdDwQEAwIBBjArBgNVHR8EJDAiMCCgHqAchhpodHRwczovL3Rlc3Qub3hsLmF0
L2NhLmNybDAKBggqhkjOPQQDAgNpADBmAjEA5uGISy3avy/WSvJoePza6/FtoCI5
1DWGSlmHXTL/0R59thlv0kNnPE0uKJDxJo8pAjEAqRaqAQR180Zn126/HZrSaOf9
NjC6qAYukrkqmg2sF6vBnwwcKhapQmtF4Ulz48cL
-----END CERTIFICATE-----`

	ecSubCaCert = `-----BEGIN CERTIFICATE-----
MIICKTCCAa+gAwIBAgIQeN5jQOOK0Zr8WELEVLw3lDAKBggqhkjOPQQDAjAWMRQw
EgYDVQQDDAtPWEwgVEVTVCBDQTAeFw0yNjAxMDExODI2MjlaFw00NTA5MTgxODI2
MjlaMBoxGDAWBgNVBAMMD09YTCBURVNUIFN1Yi1DQTB2MBAGByqGSM49AgEGBSuB
BAAiA2IABL2XFdQQ1rp2VKZwU+BBsUZNoyJWYldAJGSRmFFC2LxarCwdfgFHecsE
Z0N8CqiB8QHBkcs+XArT7CcQG7KGZJB95sNwYnM5gnlJoiVDQPDVDSv7MALvysgU
A9QGmTRazaOBvTCBujAMBgNVHRMEBTADAQH/MB0GA1UdDgQWBBTEJPmJ0h4k94wx
ljoivM6bSvCnXDBRBgNVHSMESjBIgBQOqDMkf+R6lJo09T6fi/15YunYtKEapBgw
FjEUMBIGA1UEAwwLT1hMIFRFU1QgQ0GCFBCfr/vK20hMsB61LJuYmehgaGwPMAsG
A1UdDwQEAwIBBjArBgNVHR8EJDAiMCCgHqAchhpodHRwczovL3Rlc3Qub3hsLmF0
L2NhLmNybDAKBggqhkjOPQQDAgNoADBlAjBgMhV0//IczKpWW6CLWP/FVzkwTjtd
2hNT/CUSSJkcKI3sK5yd3i2xQCveuuJM3ocCMQChzmo5yrlkOmcn6n5S3Y1qwp78
+y7AK6Ja0Yoz4ZBVfUSeHmBdfHwGZV/L0/6lXs8=
-----END CERTIFICATE-----`

	ecApp01Cert = `-----BEGIN CERTIFICATE-----
MIICZDCCAemgAwIBAgIQeiFENTqXTY5ReYWH18kQxTAKBggqhkjOPQQDAjAaMRgw
FgYDVQQDDA9PWEwgVEVTVCBTdWItQ0EwHhcNMjYwMTAxMTgyNjM2WhcNNDUwOTE4
MTgyNjM2WjAZMRcwFQYDVQQDDA5PWEwgVEVTVCBBcHAwMTB2MBAGByqGSM49AgEG
BSuBBAAiA2IABI+NnpByFqpxXPmHjzgGelhztw7rvFgOnqmQhY0Zef3Ksu9I0eop
L5HO4BWahxiQrAc7BjJeubqzFYU0wunzxhysFfuaHfalQ8+zbyeVPoHRysavl1wh
lRP1/mzb9/NkkqOB9DCB8TAJBgNVHRMEAjAAMB0GA1UdDgQWBBT8LEwfHt+VUbhA
r18pAkYR9fS9yTBNBgNVHSMERjBEgBTEJPmJ0h4k94wxljoivM6bSvCnXKEapBgw
FjEUMBIGA1UEAwwLT1hMIFRFU1QgQ0GCEHjeY0DjitGa/FhCxFS8N5QwEwYDVR0l
BAwwCgYIKwYBBQUHAwEwCwYDVR0PBAQDAgWgMCsGA1UdHwQkMCIwIKAeoByGGmh0
dHBzOi8vdGVzdC5veGwuYXQvY2EuY3JsMCcGA1UdEQQgMB6CC3Rlc3Qub3hsLmF0
gg93d3cudGVzdC5veGwuYXQwCgYIKoZIzj0EAwIDaQAwZgIxAOBsOb+qybSbgJa3
qJK3gOF3/n5kjfM/i8MI3BlYhTMbvlSZX233ifUIkRAtYo/KXQIxAIlb6atGJCl/
w2ymJ4kJE0CBPd/jbxJ6mplGhxSbDQunUh/Ia6ddo/QFXyWEBVuQwQ==
-----END CERTIFICATE-----`

	ecApp01Key = `-----BEGIN PRIVATE KEY-----
MIG2AgEAMBAGByqGSM49AgEGBSuBBAAiBIGeMIGbAgEBBDCkB5jHGUU9tRf8LNvP
3LzCX2dEU4H6r5eBaRL6heFZBZNT7ucscehU1w03sLyPy5GhZANiAASPjZ6Qchaq
cVz5h484BnpYc7cO67xYDp6pkIWNGXn9yrLvSNHqKS+RzuAVmocYkKwHOwYyXrm6
sxWFNMLp88YcrBX7mh32pUPPs28nlT6B0crGr5dcIZUT9f5s2/fzZJI=
-----END PRIVATE KEY-----`

	ecApp02Cert = `-----BEGIN CERTIFICATE-----
MIILwjCCC0igAwIBAgIRAKNS8w8w3C3Dr1jStaeW0FwwCgYIKoZIzj0EAwIwGjEY
MBYGA1UEAwwPT1hMIFRFU1QgU3ViLUNBMB4XDTI2MDEwMTE4MjgwM1oXDTQ1MDkx
ODE4MjgwM1owGTEXMBUGA1UEAwwOT1hMIFRFU1QgQXBwMDIwdjAQBgcqhkjOPQIB
BgUrgQQAIgNiAATqK5ceyfRukB/M0DhWc6ZFe9bB0i6joK3BlX2rwokdDobi2OxV
6i497TmMYOO/YSpzPTpA1xVYynbv4oMo+wH2JxSqFQ8/xgLKz/NKLDbFKSMwDKHj
lInOxMTyAkGMVDOjggpRMIIKTTAJBgNVHRMEAjAAMB0GA1UdDgQWBBS71aB7hTXf
dNUL4Djpyy8G1s5ZxzBNBgNVHSMERjBEgBTEJPmJ0h4k94wxljoivM6bSvCnXKEa
pBgwFjEUMBIGA1UEAwwLT1hMIFRFU1QgQ0GCEHjeY0DjitGa/FhCxFS8N5QwEwYD
VR0lBAwwCgYIKwYBBQUHAwEwCwYDVR0PBAQDAgWgMCsGA1UdHwQkMCIwIKAeoByG
Gmh0dHBzOi8vdGVzdC5veGwuYXQvY2EuY3JsMIIJgQYDVR0RBIIJeDCCCXSCFHhq
eHVvdGloLnRlc3Qub3hsLmF0ghRoY2tqaWdlai50ZXN0Lm94bC5hdIIUdWhqeW9w
cnoudGVzdC5veGwuYXSCFHBsamJkaHNkLnRlc3Qub3hsLmF0ghRxYmZsZmhsbC50
ZXN0Lm94bC5hdIIUdWhsZ2lhZmsudGVzdC5veGwuYXSCFG9ucmFsdXVsLnRlc3Qu
b3hsLmF0ghR2cmhxd2Z3eS50ZXN0Lm94bC5hdIIUZWZjeWVqb3kudGVzdC5veGwu
YXSCFHBia2hqd3Z1LnRlc3Qub3hsLmF0ghRibWt3bmpweC50ZXN0Lm94bC5hdIIU
emJteHh3cnkudGVzdC5veGwuYXSCFHh3enR2a2J1LnRlc3Qub3hsLmF0ghRybXR5
a2l6ai50ZXN0Lm94bC5hdIIUZG1hbWRibWkudGVzdC5veGwuYXSCFGViaXB0Z3Zj
LnRlc3Qub3hsLmF0ghR6Z2JnYmd3eC50ZXN0Lm94bC5hdIIUcHhiemVuY2sudGVz
dC5veGwuYXSCFHlvY2V1YWtxLnRlc3Qub3hsLmF0ghRwbHBkamh3ci50ZXN0Lm94
bC5hdIIUZWdudGhicXcudGVzdC5veGwuYXSCFHlnd3pjY3NvLnRlc3Qub3hsLmF0
ghRod3l2aGJrYi50ZXN0Lm94bC5hdIIUa2x5cGpybnQudGVzdC5veGwuYXSCFHZn
andweWdlLnRlc3Qub3hsLmF0ghRjcHNyYXN1ai50ZXN0Lm94bC5hdIIUeGhua3B2
anoudGVzdC5veGwuYXSCFHNrZ2ZxcmFiLnRlc3Qub3hsLmF0ghR5ZWVxeHFsZy50
ZXN0Lm94bC5hdIIUaWx5cm1hcXoudGVzdC5veGwuYXSCFHdieWlhaGtoLnRlc3Qu
b3hsLmF0ghRidmFsbG9uai50ZXN0Lm94bC5hdIIUaGFscGpmcHYudGVzdC5veGwu
YXSCFGN5Y2Jrcmd3LnRlc3Qub3hsLmF0ghRrZXVwbnF5dC50ZXN0Lm94bC5hdIIU
ZW5va2pnZ2UudGVzdC5veGwuYXSCFGdlbHJneWRmLnRlc3Qub3hsLmF0ghRudWVq
Zndqci50ZXN0Lm94bC5hdIIUZGh0am9kYmgudGVzdC5veGwuYXSCFGxweXVodXNw
LnRlc3Qub3hsLmF0ghR3eGplZ2pleS50ZXN0Lm94bC5hdIIUZm1la3dncGIudGVz
dC5veGwuYXSCFGV6Z29ncmpzLnRlc3Qub3hsLmF0ghRpbmFuZ2toeS50ZXN0Lm94
bC5hdIIUbHdtd2NpaW8udGVzdC5veGwuYXSCFGFmb3Nzd2ZlLnRlc3Qub3hsLmF0
ghRydXdoeWFwdC50ZXN0Lm94bC5hdIIUYnNjeXpvZmUudGVzdC5veGwuYXSCFHhu
a3JhZG5lLnRlc3Qub3hsLmF0ghR1cnpid3drdy50ZXN0Lm94bC5hdIIUcHpscHBw
Z20udGVzdC5veGwuYXSCFGllZnVydXpvLnRlc3Qub3hsLmF0ghR4eW91anNwdy50
ZXN0Lm94bC5hdIIUYm9kZm9qbHcudGVzdC5veGwuYXSCFHh1dGZiZ29rLnRlc3Qu
b3hsLmF0ghRsZmJjdnpveS50ZXN0Lm94bC5hdIIUeGdyZGlqbHYudGVzdC5veGwu
YXSCFGtsamtwc2xkLnRlc3Qub3hsLmF0ghRidmllbmh5cy50ZXN0Lm94bC5hdIIU
a2d4ZWZ4dnkudGVzdC5veGwuYXSCFGNleHdhZmNnLnRlc3Qub3hsLmF0ghRxeGFl
c3dqYy50ZXN0Lm94bC5hdIIUeWZ6dGtrdXcudGVzdC5veGwuYXSCFHppaHFxaXZi
LnRlc3Qub3hsLmF0ghRmY3dwaWVrai50ZXN0Lm94bC5hdIIUZ2F0c3JndXMudGVz
dC5veGwuYXSCFGtvZHd2ZWxuLnRlc3Qub3hsLmF0ghRod3pqdHd4di50ZXN0Lm94
bC5hdIIUcGllY2d4aHoudGVzdC5veGwuYXSCFGZwZ2tuYnJ2LnRlc3Qub3hsLmF0
ghR6a3BkeGtleS50ZXN0Lm94bC5hdIIUcmhib3d2Y3IudGVzdC5veGwuYXSCFHdt
ZGJta3piLnRlc3Qub3hsLmF0ghRvZG9rZ2RrbC50ZXN0Lm94bC5hdIIUbmVwZGx3
eWoudGVzdC5veGwuYXSCFGRjYmNtaWptLnRlc3Qub3hsLmF0ghRkZHBid3pvdy50
ZXN0Lm94bC5hdIIUc3J1b3B3dGUudGVzdC5veGwuYXSCFGV3dWRqbXR3LnRlc3Qu
b3hsLmF0ghRuYXF3emR1dC50ZXN0Lm94bC5hdIIUbmJ3b3dwd2wudGVzdC5veGwu
YXSCFGxtb3ZqZ3dqLnRlc3Qub3hsLmF0ghRjcWtwanhlby50ZXN0Lm94bC5hdIIU
Z2ZydmJ0c2cudGVzdC5veGwuYXSCFGF6bHhsbHluLnRlc3Qub3hsLmF0ghRxeHFj
Y2Nkai50ZXN0Lm94bC5hdIIUdGZrYnRxamsudGVzdC5veGwuYXSCFGlzemdzemRm
LnRlc3Qub3hsLmF0ghRncm9rdmplci50ZXN0Lm94bC5hdIIUdGxramRsa2UudGVz
dC5veGwuYXSCFGR2ZXV0ZWdqLnRlc3Qub3hsLmF0ghRybGJibGRsby50ZXN0Lm94
bC5hdIIUYXFnbWdtZGoudGVzdC5veGwuYXSCFGZ2aGltaHNsLnRlc3Qub3hsLmF0
ghRheXNzYm1jYy50ZXN0Lm94bC5hdIIUbW91dmxhbncudGVzdC5veGwuYXSCFHN0
b2dyb25vLnRlc3Qub3hsLmF0ghRha2VnaHl3eC50ZXN0Lm94bC5hdIIUcGx5YXNq
dmYudGVzdC5veGwuYXSCFHZvemJzZnVwLnRlc3Qub3hsLmF0ghRndmxqZWpodC50
ZXN0Lm94bC5hdIIUaWpjcHdtaHIudGVzdC5veGwuYXSCFHBjZGJzamhqLnRlc3Qu
b3hsLmF0ghRiY2JvZmlwcS50ZXN0Lm94bC5hdIIUYWZoeWV5cGgudGVzdC5veGwu
YXSCFHZtaXlyc3Z2LnRlc3Qub3hsLmF0ghRkYnJnd2FnaS50ZXN0Lm94bC5hdIIU
dmZnc2l1YXgudGVzdC5veGwuYXSCFGZwbWJhYnNrLnRlc3Qub3hsLmF0ghRlY3N5
a3Bxdy50ZXN0Lm94bC5hdDAKBggqhkjOPQQDAgNoADBlAjEAiufIWnPg8eYJjQEC
Oj40hgW6VMYUoPzDsfyt0soOnwPMRlAiLoUPeu1jc+SurZkPAjA1rMzZso02atOn
3zgANX1da4OPM0MAmSLRR/CyykJ9wib47wgtiEBDDEIsV5fPQNM=
-----END CERTIFICATE-----`

	ecApp02Key = `-----BEGIN PRIVATE KEY-----
MIG2AgEAMBAGByqGSM49AgEGBSuBBAAiBIGeMIGbAgEBBDB6f4Pf4a6mNQDp9ZcR
PKOPaQNKsGlIltSTu9+hePCUKpbR5x9+aBpAsB7OxdALOR2hZANiAATqK5ceyfRu
kB/M0DhWc6ZFe9bB0i6joK3BlX2rwokdDobi2OxV6i497TmMYOO/YSpzPTpA1xVY
ynbv4oMo+wH2JxSqFQ8/xgLKz/NKLDbFKSMwDKHjlInOxMTyAkGMVDM=
-----END PRIVATE KEY-----`

	// expired
	ecApp03Key = `-----BEGIN PRIVATE KEY-----
MIG2AgEAMBAGByqGSM49AgEGBSuBBAAiBIGeMIGbAgEBBDC3Dlaxs9cC6VsD/7nD
5ARtzQbgBonZ+5ex5Kqan9QACJWhUv0XwKHqCESHBTRI4RmhZANiAAQ8Zzu8WBp5
yXhIn8qNVBXg0kPSGGrtPIyXnpfA2G9YfDCZa1+5VsMlOVl+W1cUcGZvkh1nkp6c
t0g3azGvDm7o6FZJ4rLkMlvZqIAo1geLhu5Ki7ea9OSVi2yYCtWIEc0=
-----END PRIVATE KEY-----`

	ecApp03Cert = `-----BEGIN CERTIFICATE-----
MIICYzCCAemgAwIBAgIQbCzBbT06N+88XdZafEWTojAKBggqhkjOPQQDAjAaMRgw
FgYDVQQDDA9PWEwgVEVTVCBTdWItQ0EwHhcNMjYwMTAxMTg0MzU3WhcNMjYwMTAy
MTg0MzU3WjAZMRcwFQYDVQQDDA5PWEwgVEVTVCBBcHAwMzB2MBAGByqGSM49AgEG
BSuBBAAiA2IABDxnO7xYGnnJeEifyo1UFeDSQ9IYau08jJeel8DYb1h8MJlrX7lW
wyU5WX5bVxRwZm+SHWeSnpy3SDdrMa8ObujoVknisuQyW9mogCjWB4uG7kqLt5r0
5JWLbJgK1YgRzaOB9DCB8TAJBgNVHRMEAjAAMB0GA1UdDgQWBBTud+y69BBYUF8S
7UnjfzHSLTDWqzBNBgNVHSMERjBEgBTEJPmJ0h4k94wxljoivM6bSvCnXKEapBgw
FjEUMBIGA1UEAwwLT1hMIFRFU1QgQ0GCEHjeY0DjitGa/FhCxFS8N5QwEwYDVR0l
BAwwCgYIKwYBBQUHAwEwCwYDVR0PBAQDAgWgMCsGA1UdHwQkMCIwIKAeoByGGmh0
dHBzOi8vdGVzdC5veGwuYXQvY2EuY3JsMCcGA1UdEQQgMB6CC3Rlc3Qub3hsLmF0
gg93d3cudGVzdC5veGwuYXQwCgYIKoZIzj0EAwIDaAAwZQIwUNsyn6NKuH7id10B
JHJqe7xigyWQo0Yqr0L9JePqN7IPowfF4byy/CJ2+hShgj/FAjEA6qYda7onFkLY
qVgtNLtbYdwr+v0xbd6fSdQ/bD+ahpisDoWVbWq1qseFmdTfcsq5
-----END CERTIFICATE-----`
)

func setupTestDir(t *testing.T) string {
	tmpDir, _ := os.MkdirTemp("", "manager_test")
	config.PathCertsPublic = filepath.Join(tmpDir, config.DIR_CERTS_PUBLIC)
	config.PathCertsPrivate = filepath.Join(tmpDir, config.DIR_CERTS_PRIVATE)
	config.PathCertsBundlePublic = filepath.Join(tmpDir, config.DIR_BUNDLE_PUBLIC)
	config.PathCertsBundlePrivate = filepath.Join(tmpDir, config.DIR_BUNDLE_PRIVATE)
	config.PathAccountCache = filepath.Join(tmpDir, config.DIR_ACCOUNT)
	os.MkdirAll(config.PathCertsPublic, 0755)
	os.MkdirAll(config.PathCertsPrivate, 0755)
	os.MkdirAll(config.PathCertsBundlePublic, 0755)
	os.MkdirAll(config.PathCertsBundlePrivate, 0755)
	os.MkdirAll(config.PathAccountCache, 0755)

	config.Config = &config.ConfigFile{
		PathCerts:    tmpDir,
		CreateBundle: u.PtrBool(false),
	}
	return tmpDir
}

func TestCorrectKeyPairs(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		baseName string
		cert     string
		key      string
		domains  []string
	}{
		{"RSA App01", "rsa_app01", rsaApp01Cert, rsaApp01Key, []string{"test.oxl.at", "www.test.oxl.at"}},
		{"EC App1", "ec_app1", ecApp01Cert, ecApp01Key, []string{"test.oxl.at", "www.test.oxl.at"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.WriteFile(filepath.Join(tmpDir, config.DIR_CERTS_PUBLIC, tt.baseName+".crt"), []byte(tt.cert), 0644)
			os.WriteFile(filepath.Join(tmpDir, config.DIR_CERTS_PRIVATE, tt.baseName+".key"), []byte(tt.key), 0600)

			needed, reason := NeedsUpdate(tt.baseName, tt.domains)
			if needed {
				t.Errorf("%s failed: reason: %s", tt.name, reason)
			}
		})
	}
}

func TestBadPath_CryptoMismatch_Extended(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name           string
		cert           string
		key            string
		domains        []string
		expectedReason string
	}{
		{
			name:           "EC cert <=> RSA key",
			cert:           ecApp01Cert,
			key:            rsaApp01Key,
			domains:        []string{"test.oxl.at", "www.test.oxl.at"},
			expectedReason: "private key and public key mismatch",
		},
		{
			name:           "EC cert <=> Wrong EC key (App01 vs App02)",
			cert:           ecApp01Cert,
			key:            ecApp02Key,
			domains:        []string{"test.oxl.at", "www.test.oxl.at"},
			expectedReason: "private key and public key mismatch",
		},
		{
			name:           "EC key <=> Sub-CA cert (Wrong cert type)",
			cert:           ecSubCaCert,
			key:            ecApp01Key,
			domains:        []string{"test.oxl.at", "www.test.oxl.at"},
			expectedReason: "private key and public key mismatch",
		},
		{
			name:           "RSA cert <=> Wrong RSA key (Missing pair)",
			cert:           rsaApp01Cert,
			key:            rsaSubCaCert, // Using a CA cert as a placeholder for a 'wrong' key block
			expectedReason: "invalid key pem",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := "mismatch_" + tt.name
			os.WriteFile(filepath.Join(tmpDir, config.DIR_CERTS_PUBLIC, base+".crt"), []byte(tt.cert), 0644)
			os.WriteFile(filepath.Join(tmpDir, config.DIR_CERTS_PRIVATE, base+".key"), []byte(tt.key), 0600)

			needed, reason := NeedsUpdate(base, tt.domains)
			if !needed || reason != tt.expectedReason {
				t.Errorf("%s: expected reason '%s', got '%s'", tt.name, tt.expectedReason, reason)
			}
		})
	}
}

func TestNeedsUpdate(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer os.RemoveAll(tmpDir)

	t.Run("Missing Certificate", func(t *testing.T) {
		needed, reason := NeedsUpdate("nonexistent", []string{"test.at"})
		if !needed || reason != "certificate file missing" {
			t.Errorf("Expected missing cert error, got: %s", reason)
		}
	})

	t.Run("Invalid PEM", func(t *testing.T) {
		base := "invalid_pem"
		os.WriteFile(filepath.Join(tmpDir, config.DIR_CERTS_PUBLIC, base+".crt"), []byte("not a certificate"), 0644)
		os.WriteFile(filepath.Join(tmpDir, config.DIR_CERTS_PRIVATE, base+".key"), []byte(rsaApp01Key), 0600)

		needed, reason := NeedsUpdate(base, []string{"test.at"})
		if !needed || reason != "invalid certificate pem" {
			t.Errorf("Expected invalid pem error, got: %s", reason)
		}
	})

	t.Run("Domain Sorting Persistence", func(t *testing.T) {
		base := "sort_test"
		os.WriteFile(filepath.Join(tmpDir, config.DIR_CERTS_PUBLIC, base+".crt"), []byte(rsaApp01Cert), 0644)
		os.WriteFile(filepath.Join(tmpDir, config.DIR_CERTS_PRIVATE, base+".key"), []byte(rsaApp01Key), 0600)

		// Provide domains in reverse order of how they appear in the cert [test.oxl.at, www.test.oxl.at]
		revDomains := []string{"www.test.oxl.at", "test.oxl.at"}

		needed, reason := NeedsUpdate(base, revDomains)
		if needed {
			t.Errorf("Should NOT need update; sorting should handle domain order. Reason: %s", reason)
		}
	})

	t.Run("Missing Bundle File", func(t *testing.T) {
		base := "bundle_test"
		config.Config.CreateBundle = u.PtrBool(true)
		os.WriteFile(filepath.Join(tmpDir, config.DIR_CERTS_PUBLIC, base+".crt"), []byte(ecApp01Cert), 0644)
		os.WriteFile(filepath.Join(tmpDir, config.DIR_CERTS_PRIVATE, base+".key"), []byte(ecApp01Key), 0600)

		needed, reason := NeedsUpdate(base, []string{"test.oxl.at", "www.test.oxl.at"})
		if !needed || reason != "public-bundle file missing" {
			t.Errorf("Expected update needed due to missing bundle file, got needed=%v, reason=%s", needed, reason)
		}
	})
}

func TestPublicKeyMatches(t *testing.T) {
	block, _ := pem.Decode([]byte(rsaApp01Cert))
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("failed to parse cert: %v", err)
	}

	keyBlock, _ := pem.Decode([]byte(rsaApp01Key))
	privKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		privKey, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	}
	if err != nil {
		t.Fatalf("failed to parse RSA key: %v", err)
	}

	if !publicKeyMatches(cert.PublicKey, privKey) {
		t.Error("Helper failed to match valid RSA pair")
	}

	ecKeyBlock, _ := pem.Decode([]byte(ecApp01Key))
	ecPrivKey, err := x509.ParsePKCS8PrivateKey(ecKeyBlock.Bytes)
	if err != nil {
		ecPrivKey, err = x509.ParseECPrivateKey(ecKeyBlock.Bytes)
	}
	if err != nil {
		t.Fatalf("failed to parse EC key: %v", err)
	}

	if publicKeyMatches(cert.PublicKey, ecPrivKey) {
		t.Error("Helper incorrectly matched RSA public key with EC private key")
	}
}

func TestNeedsUpdate_Expiry(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer os.RemoveAll(tmpDir)

	base := "expired_app"
	os.WriteFile(filepath.Join(tmpDir, config.DIR_CERTS_PUBLIC, base+".crt"), []byte(ecApp03Cert), 0644)
	os.WriteFile(filepath.Join(tmpDir, config.DIR_CERTS_PRIVATE, base+".key"), []byte(ecApp03Key), 0600)

	needed, reason := NeedsUpdate(base, []string{"test.oxl.at", "www.test.oxl.at"})
	if !needed || (reason != "expired" && reason != "expiring soon") {
		t.Errorf("Expected renewal due to expiry, got needed=%v, reason=%s", needed, reason)
	}
}
