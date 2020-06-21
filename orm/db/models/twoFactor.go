package models

import (
	"encoding/base32"
	"strings"
	"time"

	"github.com/hgfischer/go-otp"

	"cherry/utils"
)

type TwoFactor struct {
	Id                int `gorm:"primary_key"`
	MemberId          int
	Otp               string `sql:"-"`
	OtpSecret         string
	LastVerifyAt      time.Time
	LastVerifyAtStamp int64 `sql:"-"`
	Activated         bool  `gorm:"default:null"`
	Type              string
	RefreshedAt       time.Time
	RefreshedAtStamp  int64 `sql:"-"`
}

func (tf *TwoFactor) InitializeTimestamp() {
	var nilTime time.Time
	if !tf.LastVerifyAt.Equal(nilTime) {
		tf.LastVerifyAtStamp = tf.LastVerifyAt.Unix()
	}
	if !tf.RefreshedAt.Equal(nilTime) {
		tf.RefreshedAtStamp = tf.RefreshedAt.Unix()
	}
}

func (tf *TwoFactor) Verify() bool {
	result := true
	switch tf.Type {
	case "TwoFactor::Sms":
		result = tf.smsVerify()
	}
	return result
}

func (tf *TwoFactor) GetCurrentOtpCode() string {
	totp := otp.TOTP{Secret: strings.ToUpper(tf.OtpSecret), IsBase32Secret: true}
	return totp.Get()
}

func (tf *TwoFactor) smsVerify() bool {
	return tf.OtpSecret == tf.Otp
}

func (tf *TwoFactor) GenerateOtpSecret() {
	secret := base32.StdEncoding.EncodeToString([]byte(utils.RandStringRunes(16)))
	var nb [16]byte
	for i, c := range []byte(secret) {
		if i > 15 {
			break
		}
		nb[i] = c
	}
	tf.OtpSecret = strings.ToLower(string(nb[:]))
}

func (tf *TwoFactor) NeedRefreshLastVerifyAt() bool {
	if tf.LastVerifyAt.Before(time.Now().Add(-time.Minute * 10)) {
		return true
	}
	return false
}
