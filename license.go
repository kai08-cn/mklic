package mklic

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/zfs123/gorsa"
)

//Sign sign a license with private key
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
	license := aesEncrypt(lic, devid)
	return license, nil
}

//Verify verify a license and return the content
func Verify(lic, pubkey []byte, devid string) ([]byte, error) {
	if devid == "" {
		devid = string(GenDevId(""))
	}
	cleartext, err := aesDecrypt(lic, devid)
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

//GenDevId generate a device id with macs (if not define, use macs of local ethernet interfaces)
func GenDevId(macstr string) []byte {
	macs, nmac := getMacs(macstr)
	if nmac == 0 {
		return nil
	}
	macsha1 := macSha1(macs, nmac)
	b32 := macBase32(macsha1)
	return b32
}
