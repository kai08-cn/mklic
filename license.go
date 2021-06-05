package lictool

import (
	"crypto/aes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"sort"

	"github.com/wenzhenxi/gorsa"
)

/*
type License struct {
	plain     string
	cleartext string
	devid     string
}

func NewLicense() *License {
	return &License{}
}
*/

func Sign() (string, error) {
	return "", nil
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

	verifyfp, err := computeFingerprint(m)
	if err != nil {
		return nil, err
	}
	if verifyfp != clearfp {
		return nil, errors.New("fingerprint error")
	}

	m["fingerprint"] = clearfp
	out, err := json.MarshalIndent(m, "", "")
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
		if i != 0 {
			text += "&"
		}
		tmp, ok := m[v].(string)
		if !ok {
			return "", errors.New("fingerprint error")
		}
		text = text + v + "=" + tmp
	}
	h := sha1.New()
	h.Write([]byte(text))
	sha1finger := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(sha1finger), nil
}

func licAesDecrypt(in []byte, key string) ([]byte, error) {
	content := make([]byte, 0)
	_, err := base64.StdEncoding.Decode(content, in)
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

func PKCS7UnPadding(s []byte) []byte {
	length := len(s)
	padding := int(s[length-1])
	return s[:(length - padding)]
}

func AesDecrypt(data, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	decrypted := make([]byte, len(data))
	size := block.BlockSize()

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Decrypt(decrypted[bs:be], data[bs:be])
	}
	return decrypted
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
