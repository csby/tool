package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	mutex sync.RWMutex

	Ca       ConfigCa     `json:"ca"`
	Ta       string       `json:"ta"`
	Template string       `json:"template"`
	Client   ConfigClient `json:"client"`
}

func (s *Config) LoadFromFile(filePath string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, s)
}

func (s *Config) SaveToFile(filePath string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	bytes, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return err
	}

	fileFolder := filepath.Dir(filePath)
	_, err = os.Stat(fileFolder)
	if os.IsNotExist(err) {
		os.MkdirAll(fileFolder, 0777)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprint(file, string(bytes[:]))

	return err
}

func (s *Config) String() string {
	bytes, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(bytes[:])
}

func (s *Config) FormatString() string {
	bytes, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return ""
	}

	return string(bytes[:])
}

func (s *Config) Example() {
	s.Ca.Cert = `-----BEGIN CERTIFICATE-----
MIIDejCCAmKgAwIBAgIRAMFh7yWv3oAiHLq3LYKmyLgwDQYJKoZIhvcNAQELBQAw
ZjELMAkGA1UEBhMCQ04xCTAHBgNVBAgTADEJMAcGA1UEBxMAMQkwBwYDVQQJEwAx
CzAJBgNVBAoTAmNhMQwwCgYDVQQLEwNkZXYxGzAZBgNVBAMMEuW8gOWPkeeUqOag
ueivgeS5pjAeFw0xOTAxMTcwNTU0MzNaFw0yOTAxMTQwNTU0MzNaMGYxCzAJBgNV
BAYTAkNOMQkwBwYDVQQIEwAxCTAHBgNVBAcTADEJMAcGA1UECRMAMQswCQYDVQQK
EwJjYTEMMAoGA1UECxMDZGV2MRswGQYDVQQDDBLlvIDlj5HnlKjmoLnor4HkuaYw
ggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDQPF4EMpvXEOvpry691sOj
laC2dixmMiFMsyxHDbNDcGZ/lMXwv4qhG/vZSNqURsHsNcjPhGAqURZKDeN0lpX6
FQL7ge/B3wvtaDK3o4RwrDmmdnFyuxFWraOfzUNTqp7T+fKJ3k2jkeHWzJ8xJK/E
sfBfrKkOm60kBbA8uAjwwwEhZWN1ChpTY6TQp3tvusIvbA+zTeROMOWcv71o5Zjv
yPpFq8zo1AquznISVZ6FfcSIlKGq2CZA6FGxBfzhXmiXRZ/hKO6pUlbGT63/60BY
nOd7Pr6nDFwcxiO9eVZiwJuj2ZQYmRronkBLHHtqLilWZw8akLtsq1MWp7DQkHoL
AgMBAAGjIzAhMA4GA1UdDwEB/wQEAwICpDAPBgNVHRMBAf8EBTADAQH/MA0GCSqG
SIb3DQEBCwUAA4IBAQA4899EzH9SkGWZiyLsPzN3t3hUu/j7k3CZgx5SRJ0sGgzh
LhZQnTcqfdDJJBYrTlLjipHF77CU3BL2W9LhR4GEW/ZS5SG7jd0k9za+o0oifGxo
/r03EH0T1EipnGRUUNbIQaXYPNZvI4ViwEEGFxMh+ii6h+s2OKGatTPzj1CJL7Vv
xuCWeDR2ZN5Z5HhxLK2g29Ab6e0XlVX9iTw7vNjUEkdJ9UX6qfNoIAEtZMDgnf1N
d3CD2xK0SZjV8cv+iEA2YjLV3Xcdx9HlFaRRAY5pWBS3dqQ+e4/Gjr9Ut1V9QEPX
+kdXNC25hhcWvt6td0WE1GkZ+kaXoTIpWgf4kTF9
-----END CERTIFICATE-----`

	s.Ca.Key = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA0DxeBDKb1xDr6a8uvdbDo5WgtnYsZjIhTLMsRw2zQ3Bmf5TF
8L+KoRv72UjalEbB7DXIz4RgKlEWSg3jdJaV+hUC+4Hvwd8L7Wgyt6OEcKw5pnZx
crsRVq2jn81DU6qe0/nyid5No5Hh1syfMSSvxLHwX6ypDputJAWwPLgI8MMBIWVj
dQoaU2Ok0Kd7b7rCL2wPs03kTjDlnL+9aOWY78j6RavM6NQKrs5yElWehX3EiJSh
qtgmQOhRsQX84V5ol0Wf4SjuqVJWxk+t/+tAWJznez6+pwxcHMYjvXlWYsCbo9mU
GJka6J5ASxx7ai4pVmcPGpC7bKtTFqew0JB6CwIDAQABAoIBACBlg+r7RKuNAmb0
zjzSsNU+biZ09CtiGTQpm/Xh98XCMvAeYT03T4YJKIGBiCARchIhvAAtBBkRTHpw
9rYox2SE6FXgvUBYRy7ESz+uvOgao012l+fVmrj1gsNV1+eoX9VyyX0RjNqp04zr
gMjQgFgFpvP7lMGlDqbQC78mkp2CMUrZgxYyoYvtyDchBNGi74/gF4XB84cUF41q
sEJvRlIAE2+X4xAvGYsXSr14R0B+FdNWmRTRhy5T7C6VRQenQ++xHA18se17zJDE
xovp+P1iGHSqW8vwhN1ewlQS+V8ReYhOcxTtryYDSZRHytoUzy1lmRbk+2G2m99v
vT0Fi2ECgYEA4dqwZS/kGkqdoTySJuzSi9NQXthyHMI/XlJkZdnXQCv0JEPtm2He
V6FgHweqsdzdBYiYZ87gaJJND7w5V2R6RTlxJcKWy2xtPsklfvqynCB77FwcN/kE
AUmBe7LZnqBxak5ZZCp7VJ+OHSqzQF/twuqf7DLUzKhuwr5LX6NUAKUCgYEA7Aeq
f+88gNLbzdMfk1CY12Gj/dbnmR8V0ZJpDx6Vm+vxn66Aj+5Zld9P2P+LXu4xgprX
DklfYuVDhJRy+IWFV1DZGnGCLwi6QifqOdmFP4Cp0q3HBXeXP6lMGI3a76Xqi5PE
zB3nCZUHTzoWS8PbiQO2GdqgQcF5lh5r3fZuYO8CgYEAiiGIFLrPlUzhT0WOVYQt
2RqqYveaAwDCZNubT7eFsFexf6ST455dr9agxRmZSiK8gq/iFksucIZZ3y/NGif3
p/LTwrZaJ5vuzKGU7y5AosAzSoGjOJBx5J+iM5dVBXE3LD8y6NIaj8zty+TbsJl6
/uUkXf9QqsKwyyY7TwjDTYkCgYEAr/brIuPdrIEHA8TBRpeQywntM8Jy+VIWdx57
Gp2Hli3p/k0fZa6htT0+Dx487nIQETmU2P7UkSsxEfeGW0dX6IfKw48aKiyMh+Ow
GJ7VW0l10i5iMO4+oWR66ddgAMgmOxbYYgBtDVTAlU1N8AA9fEF7En1OepdeKQ/X
BSzCj5sCgYEA3Q7At47u06neS6xvk63kqnf7BOoK7XV46JGJAxx+Ykt1P+Ipvx07
JiPMao1nUnEjYVkKd66ePiIGp7FJ6tRd9uDwq3yH9htDLXPzG0V6jN1Re2n2DHPx
L2L5G8C9+yYyIvaM+4wDXuFiA7X9Eo46NTzEtx0tzj+3b6+OK/S33Do=
-----END RSA PRIVATE KEY-----`

	s.Ca.KeyPassword = ""

	s.Template = `client
dev tun
proto udp
remote 172.16.115.96 
port 1194
resolv-retry infinite
nobind
persist-key
persist-tun
cipher AES-256-CBC
remote-cert-tls server
key-direction 1
verb 3

<ca>
%s
</ca>

<tls-auth>
%s
</tls-auth>

<cert>
%s
</cert>

<key>
%s
</key>`

	s.Client.ExpiredDays = 365
}
