package mklic

import (
	"bytes"
	"crypto/aes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/zfs123/gorsa"
)

func Sign(prikey, data []byte, devid string) (string, error) {
	if devid == "" {
		devid = string(GenDevId(""))
	}

	m := make(map[string]interface{})
	if data != nil {
		if err := json.Unmarshal(data, &m); err != nil {
			return "", err
		}
	}
	// add default information to license
	m["time"] = time.Now().Truncate(time.Second).Local().String()
	m["uuid"] = uuid.New().String()
	m["devid"] = devid
	m["version"] = "v1.0"

	fingerprint, err := computeFingerprint(m)
	if err != nil {
		return "", err
	}
	cryptfp, err := gorsa.PriKeyEncrypt(fingerprint, string(prikey))
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	m["fingerprint"] = cryptfp

	lic, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	license := licAesEncrypt(lic, devid)
	return license, nil
}

func Verify(lic, pubkey []byte, devid string) ([]byte, error) {
	if devid == "" {
		devid = string(GenDevId(""))
	}
	cleartext, err := licAesDecrypt(lic, devid)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal(cleartext, &m); err != nil {
		return nil, err
	}

	if _, exist := m["fingerprint"]; !exist {
		return nil, errors.New("content error")
	}

	cryptfp := m["fingerprint"].(string)
	clearfp, err := gorsa.PublicDecrypt(cryptfp, string(pubkey))
	if err != nil {
		return nil, err
	}

	delete(m, "fingerprint")
	verifyfp, err := computeFingerprint(m)
	if err != nil {
		return nil, err
	}
	if verifyfp != clearfp {
		return nil, errors.New("fingerprint error")
	}

	out, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return nil, err
	}
	return out, nil
}

func computeFingerprint(m map[string]interface{}) (string, error) {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var text string
	for i, v := range keys {
		var tmp string
		switch m[v].(type) {
		case float64:
			tmp = strconv.Itoa(int(m[v].(float64)))
		case string:
			tmp = m[v].(string)
		case bool:
			tmp = strconv.FormatBool(m[v].(bool))
		default:
			continue
		}
		if i != 0 {
			text += "&"
		}
		text = text + v + "=" + tmp
	}
	h := sha1.New()
	h.Write([]byte(text))
	sha1finger := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(sha1finger), nil
}

func licAesEncrypt(in []byte, key string) string {
	block, _ := aes.NewCipher([]byte(key))
	in = PKCS7Padding(in, block.BlockSize())
	encrypted := make([]byte, len(in))
	size := block.BlockSize()

	for bs, be := 0, size; bs < len(in); bs, be = bs+size, be+size {
		block.Encrypt(encrypted[bs:be], in[bs:be])
	}
	return base64.StdEncoding.EncodeToString(encrypted)
}

func licAesDecrypt(in []byte, key string) ([]byte, error) {
	content, err := base64.StdEncoding.DecodeString(string(in))
	if err != nil {
		return nil, err
	}

	block, _ := aes.NewCipher([]byte(key))
	decrypted := make([]byte, len(content))
	size := block.BlockSize()

	for bs, be := 0, size; bs < len(content); bs, be = bs+size, be+size {
		block.Decrypt(decrypted[bs:be], content[bs:be])
	}

	return PKCS7UnPadding(decrypted), nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(s []byte) []byte {
	length := len(s)
	padding := int(s[length-1])
	return s[:(length - padding)]
}

func GenDevId(macstr string) []byte {
	macs, nmac := getMacs(macstr)
	if nmac == 0 {
		return nil
	}
	macsha1 := macSha1(macs, nmac)
	b32 := macBase32(macsha1)
	return b32
}
