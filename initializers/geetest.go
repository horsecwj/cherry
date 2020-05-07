package initializers

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"math/rand"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Geetest struct {
	CaptchaID  string `yaml:"id"`
	PrivateKey string `yaml:"key"`
}

var GeetestConfig Geetest

func InitGeetest() {
	path_str, _ := filepath.Abs("config/geetest.yml")
	content, err := ioutil.ReadFile(path_str)
	if err != nil {
		log.Fatal(err)
	}
	yaml.Unmarshal(content, &GeetestConfig)
}

func (geetest *Geetest) GenerateToken(ip string) string {
	var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")
	b := make([]rune, 32)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return base64.StdEncoding.EncodeToString([]byte(string(b) + "-" + ip))
}
