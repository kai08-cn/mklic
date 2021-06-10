package mklic

import (
	"crypto/sha1"
	"encoding/base64"
	"sort"
	"strconv"
)

// ComputeFingerprint
// 1. sort 2. base64(sha1(key1=value1&key2=value2&k3=value3...))
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
