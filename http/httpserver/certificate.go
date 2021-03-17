package httpserver

import (
	"os"
	"path"

	"github.com/imkouga/gocore/store/filex"
)

const (
	currentShortDir           = "./"
	defaultCertificateFile    = `.server.crt`
	defaultCertificateKeyFile = `.server.key`
)

var (
	realDumpCertificateFile    = defaultCertificateFile
	realDumpCertificateKeyFile = defaultCertificateKeyFile
)
var (
	serverCrt = `-----BEGIN CERTIFICATE-----
MIIDczCCAlugAwIBAgIJAJwd4pRhh3z4MA0GCSqGSIb3DQEBCwUAMFAxCzAJBgNV
BAYTAkNOMQ8wDQYDVQQIDAZmdWppYW4xDzANBgNVBAcMBnhpYW1lbjELMAkGA1UE
CwwCZmoxEjAQBgNVBAMMCWxvY2FsaG9zdDAeFw0xOTAxMjIwNzQxMzBaFw0yOTAx
MTkwNzQxMzBaMFAxCzAJBgNVBAYTAkNOMQ8wDQYDVQQIDAZmdWppYW4xDzANBgNV
BAcMBnhpYW1lbjELMAkGA1UECwwCZmoxEjAQBgNVBAMMCWxvY2FsaG9zdDCCASIw
DQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAO5/Ge3pyNxFfpe0FXrA+VWjDgWz
Nn4tavpPIn96k1oT9I+PQXoBtnQgG5LU7IhHTaXrt0p/eu2K8cVM2vmcgQiAoqTj
vVBWvSVqALM8/IGj18XcotSlTZEh6hmdPotNcVYkLTI0KwAw8jsXmlEtDUgUUMtq
q0WbfCpuuz+K9Di5fc9DNk0doYQTGh2kslE3sSj0hohAyOLOCBXNGHmu41UeAKfe
vA+5byO9oDfRqiEulUcmcw5nH3CvGFSZ9jaXMGfb8YUhnFVIJXoZfkvR7CbTfaAS
ZE8ZGqbVcj6XYkvUN0+u+5JDMUZqufwcQR7L8rLts0S3Gbp1ZiB7zabFXDsCAwEA
AaNQME4wHQYDVR0OBBYEFH7vPv7Igv2kBCPtLr9tWa7v2f2bMB8GA1UdIwQYMBaA
FH7vPv7Igv2kBCPtLr9tWa7v2f2bMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEL
BQADggEBAMh0Wc7t64AC7PUtpvRrw6BOJKgwzEMfBvMsRuTHfVCovpv2J/DaWtSR
yM9cp6OuyJWlkTA/yDxwq9QeoaxYLJ5M4rOITGSnO2lORHhV+UkJyGPu+9oYmft0
zzYLDNmmwbkuwTWz/1hf3Novj4asltPv02pEeVfnpuIrEMR2d/IUrYzUvbCT28Pl
lLVK84LYbNal+SZDuV5EtmflImlcEcKYdaPbGRr8NB0wSA3CyNzNdbS7bBUj/vkB
5GIf4Qy1Wy+hvRFWKlhA9hzxTssxyxqrmBMMBlHXkWT181B7h24izw7P2rsuEeWO
pFZnF6wneZvUH6UE1seZhK0JpMqwSIo=
-----END CERTIFICATE-----`
	serverKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA7n8Z7enI3EV+l7QVesD5VaMOBbM2fi1q+k8if3qTWhP0j49B
egG2dCAbktTsiEdNpeu3Sn967YrxxUza+ZyBCICipOO9UFa9JWoAszz8gaPXxdyi
1KVNkSHqGZ0+i01xViQtMjQrADDyOxeaUS0NSBRQy2qrRZt8Km67P4r0OLl9z0M2
TR2hhBMaHaSyUTexKPSGiEDI4s4IFc0Yea7jVR4Ap968D7lvI72gN9GqIS6VRyZz
DmcfcK8YVJn2NpcwZ9vxhSGcVUglehl+S9HsJtN9oBJkTxkaptVyPpdiS9Q3T677
kkMxRmq5/BxBHsvysu2zRLcZunVmIHvNpsVcOwIDAQABAoIBAQC3CGRl6h10rwDQ
fCxf4Ol5h4Gjbj5L559KKqFXJEMhxl7SLicZ82aLCHkg3rgIfnBg/d3VFrDIzPFv
ceQ73JhKZi5sTTtlBKx0oj2XUR6Yf52BBCsS9ynoUBbRQZRWZECu02S8Or0lkGrW
Xu7XjbO7tZusAVkgOou4JPMfeQyk2nseMkSnwbSj7V75lCouVDJTE2WSy87j7KKn
Rq/iEJ683Pri7zss3PpkuJn5qQiOj2w0o1Faoj5b5AJlf9OdcMuLKJj6LLkQxoMC
2ZJvyp6PA9pnaMSZVQOKjnqB5qA5jIEUgzVpFs2L+oASqYPvRiJQspLCFx/xXr4Y
pri9HeRBAoGBAPpvIWKASkJpQ0Kjy/VNLZyTQZWAkfCiOivv9f/x7FiH0+zXs/IH
nqRvTWBf3Y0/aKw3KLTlMaeMe8LbykEk5cpI+gO4/RpY9IgHzc0xwxdlzClS8H34
aQMZdgxLr2dSbeRFzDcGgCYpuSDfNTyZYqdvvzU2xuBkovUJhPApknEFAoGBAPPM
DPb4kO/AHtzGV2/CXjPR5Kfy3lrrs4vpXaKvlJttlLdsHCPWOSGXJXSQJpfQgiPX
Bx5PDyWRNAjpodojI5HflY9QvDXFw/B3KJsbS2N3MWJ+3/SUXmz7lqs1uF59jKCk
zLA0j6m9Krdhqyr/HBUxkO7nU8V8wYvN7kQ0Exw/AoGAZ6fT06dyFSbolg6h/vhg
5qv0u2KqBUXAeisqUTPbNZGS4Dcv3f/VZA5FopxLYYlbU9zI1ob/FHCLUU4T2v7g
teeaxCuvZ1ZmcF96iXINZAPYi0ovDJTjMks5l0FEaqmtnoxdSHFCXYlrfPWmXVzH
frI8HFR94KcG5BF6msU6PdkCgYBFW/cQSFVLsDfXjaIQjJaqXXuVAHacHVR+aI0Y
HKXFtl5J9LrowyiL0ul4CQ7BwDNWKPXAfLONd4r7QiSm37pd5OMy28A/+Byvi+cE
gbZn/OAS6o+ikJdwn/8UhHsIfuWESn5lXv7ERqohc+rzl5KQwQI/xZZCqCSUNqkj
xEDvlwKBgB1/N7i5I2D/zlNHBMqPASGecLASt8xSIqT7zLXlZipY9l3jd0akwk3U
wX9U5jiVh22zUy9O0X5IyncmvjlhCqOt0TybVkMUKv47wkl+U9Iz2vCm2/WJz5co
odJDuzIbFIh3iZ6lQVw6ss0JpDjr8eSDfSgIKdILgTmD22NRfBX3
-----END RSA PRIVATE KEY-----`
)

func dumnCommonCertificateFile() error {

	dir := getDumpDir()
	realDumpCertificateFile = path.Clean(path.Join(dir, defaultCertificateFile))
	realDumpCertificateKeyFile = path.Clean(path.Join(dir, defaultCertificateKeyFile))
	if err := filex.WriteFullFile(realDumpCertificateFile, []byte(serverCrt)); err != nil {
		return err
	}
	return filex.WriteFullFile(realDumpCertificateKeyFile, []byte(serverKey))
}

func getDumpDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return currentShortDir
	}

	return dir
}
