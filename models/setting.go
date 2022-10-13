package models

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"strings"
)

type Setting struct {
	Id          int64  `gorm:"primary_key;AUTO_INCREMENT;unique;column:id" json:"id"`                           //协议ID
	ConfigName  string `json:"ConfigName" form:"ConfigName" gorm:"unique;not null;size:128" binding:"required"` //配置名称类型
	ConfigValue string `json:"ConfigValue" form:"ConfigValue" gorm:"not null;size:4096" binding:"required"`     //配置值
	Version     int    `json:"Version" form:"Version" gorm:"not null;size:32" binding:"required"`               //版本
}

func (setting *Setting) CreateSetting() error {
	ret, _ := setting.QueryDefaultSetting()
	if ret == nil {
		result := db.Create(setting)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (setting *Setting) CreateDefaultSetting() {
	if _, err := os.Stat("/var/decept-agent/ssh/private.pem"); os.IsNotExist(err) {
		GenerateRsaKey()
	}
}

func (setting *Setting) QueryDefaultSetting() (*Setting, error) {
	var ret Setting
	if err := db.Where("version = ?", 1).Take(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func GenerateRsaKey() error {
	zap.L().Info("generate rsa key files...")
	privateKeyDir := "/var/decept-agent/ssh/"
	privateKeyPath := "/var/decept-agent/ssh/private.pem"
	publicKeyPath := "/var/decept-agent/ssh/public.pem"

	var privateKeyFile *os.File
	var publicKeyFile *os.File
	if _, err := os.Stat(privateKeyDir); os.IsNotExist(err) {
		err = os.MkdirAll(privateKeyDir, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
			zap.L().Error(err.Error())
			return err
		}
	}

	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			fmt.Println(err.Error())
			zap.L().Error(err.Error())
			return err
		}
		derStream := x509.MarshalPKCS1PrivateKey(privateKey)
		block := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: derStream,
		}

		privateKeyFile, err = os.Create(privateKeyPath)
		defer privateKeyFile.Close()
		if err != nil {
			fmt.Println(err.Error())
			zap.L().Error(err.Error())
			return err
		}
		err = pem.Encode(privateKeyFile, block)
		if err != nil {
			fmt.Println(err.Error())
			zap.L().Error(err.Error())
			return err
		}
		// 生成公钥文件
		publicKey := &privateKey.PublicKey
		derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
		if err != nil {
			fmt.Println(err.Error())
			zap.L().Error(err.Error())
			return err
		}
		block = &pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: derPkix,
		}
		publicKeyFile, err = os.Create(publicKeyPath)
		defer publicKeyFile.Close()
		if err != nil {
			zap.L().Error(err.Error())
			fmt.Println(err.Error())
			return err
		}
		err = pem.Encode(publicKeyFile, block)
		if err != nil {
			fmt.Println(err.Error())
			zap.L().Error(err.Error())
			return err
		}
	}
	publicFileBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		fmt.Println(err.Error())
		zap.L().Error(err.Error())
		return err
	}
	keyString := GetBetweenStr(string(publicFileBytes), "-----BEGIN RSA PUBLIC KEY-----", "-----END RSA PUBLIC KEY-----")
	encodedData := base64.StdEncoding.EncodeToString([]byte(keyString))
	setting := Setting{
		ConfigName:  "SSHKey",
		ConfigValue: encodedData,
		Version:     1,
	}
	if err = setting.CreateSetting(); err != nil {
		return err
	}
	zap.L().Info("generate rsa key files success")

	return nil
}

func GetBetweenStr(str, start, end string) string {
	startLen := len(start)
	str = string([]byte(str)[startLen:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}
