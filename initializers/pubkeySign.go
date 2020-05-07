package initializers

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"net/url"
	"sort"
	"strings"

	. "cherry/models"
)

func PublicKeySignVerified(method, path string, service *Service, params *map[string]string) (bool, error) {
	plaintext := method + "|" + path + "|"
	var keys []string
	for k, _ := range *params {
		if k == "signature" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, key := range keys {
		if i > 0 {
			plaintext += "&"
		}
		plaintext += key + "=" + strings.ToLower(url.QueryEscape((*params)[key]))
	}
	block, _ := pem.Decode([]byte(service.CustomKey))
	if block == nil {
		return false, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	h := sha256.New()
	h.Write([]byte(plaintext))
	hashed := h.Sum(nil)
	opts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA256}
	decoded, err := base64.StdEncoding.DecodeString((*params)["signature"])
	if err != nil {
		return false, errors.New("base64 decode error")
	}
	if e := rsa.VerifyPSS(pub, crypto.SHA256, hashed, []byte(decoded), opts); e != nil {
		return false, nil
	} else {
		return true, nil
	}
}

func PrivateKeySign(method, path string, service *Service, params *map[string]string) (string, error) {
	plaintext := method + "|" + path + "|"
	var keys []string
	for k, _ := range *params {
		if k == "signature" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, key := range keys {
		if i > 0 {
			plaintext += "&"
		}
		plaintext += key + "=" + strings.ToLower(url.QueryEscape((*params)[key]))
	}
	block, _ := pem.Decode([]byte(service.PrivateKey))
	if block == nil {
		return "", errors.New("private key error")
	}
	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	h.Write([]byte(plaintext))
	hashed := h.Sum(nil)
	opts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA256}
	sig, err := rsa.SignPSS(rand.Reader, pri, crypto.SHA256, hashed, opts)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString([]byte(sig)), nil
}
