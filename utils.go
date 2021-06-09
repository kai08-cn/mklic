package mklic

import (
	"crypto/sha1"
	"fmt"
	"net"
	"sort"
	"strings"
)

func getMacs(macs string) ([]byte, int) {
	if macs == "" {
		// default use 2 interface macs, eth0 eth1
		macs = getNumEtherMacs(2)
	}
	nlen := len(macs)
	if nlen < 17 || nlen%17 != 0 {
		return nil, 0
	}

	num := len(macs) / 17
	res := make([]byte, 0)
	for i := 0; i < num; i++ {
		tmp := macs[i*17 : (i+1)*17]
		mac := make([]byte, 6)
		fmt.Sscanf(tmp, "%02X:%02X:%02X:%02X:%02X:%02X", &mac[0], &mac[1], &mac[2], &mac[3], &mac[4], &mac[5])
		res = append(res, mac...)
	}

	return res, num
}

func macSha1(macs []byte, nmac int) []byte {
	sha1output := make([]byte, 20)
	sha1one := make([]byte, 20)
	for i := 0; i < nmac; i++ {
		tmp := macs[i*6 : (i+1)*6]
		h := sha1.New()
		h.Write([]byte(tmp))
		sha1one = h.Sum(nil)
		for j := 0; j < 20; j++ {
			sha1output[j] = sha1output[j] ^ sha1one[j]
		}
	}
	return sha1output
}

func macBase32(s []byte) []byte {
	b32set := "23456789WXYZABCDEFGHJKLMNPQRSTUV"
	n := len(s)
	if n%5 != 0 {
		return nil
	}
	output := make([]byte, n/5*8)
	i, j := 0, 0
	for count := n; count >= 5; count -= 5 {
		inpos := s[j*5 : j*5+5]
		output[i] = b32set[inpos[0]>>3]
		i++
		output[i] = b32set[(inpos[1]&0x3E)>>1]
		i++
		output[i] = b32set[(inpos[2]&0x0F)<<1|(inpos[3]&0x80)>>7]
		i++
		output[i] = b32set[(inpos[3]&0x03)<<3|(inpos[4]&0xE0)>>5]
		i++
		output[i] = b32set[(inpos[4] & 0x1F)]
		i++
		output[i] = b32set[(inpos[0]&0x07)<<2|(inpos[1]&0xC0)>>6]
		i++
		output[i] = b32set[(inpos[1]&0x01)<<4|(inpos[2]&0xF0)>>4]
		i++
		output[i] = b32set[(inpos[3]&0x7C)>>2]
		i++
		j++
	}
	return output
}

func getNumEtherMacs(num int) string {
	macs := ""
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	eths := make([]string, 0)
	ens := make([]string, 0)
	ethinters := make(map[string]net.Interface)
	eninters := make(map[string]net.Interface)
	for _, i := range interfaces {
		if strings.HasPrefix(i.Name, "eth") {
			eths = append(eths, i.Name)
			ethinters[i.Name] = i
		} else if strings.HasPrefix(i.Name, "en") {
			ens = append(ens, i.Name)
			eninters[i.Name] = i
		} else {
			continue
		}
	}
	if len(eths) == 0 && len(ens) == 0 {
		return ""
	}
	var internames []string
	var inters map[string]net.Interface
	if len(eths) != 0 {
		internames = eths
		inters = ethinters
	} else {
		internames = ens
		inters = eninters
	}
	sort.Strings(internames)
	cnt := 0
	for _, it := range internames {
		macAddr := inters[it].HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}
		macs += macAddr
		cnt++
		if cnt >= num {
			break
		}
	}
	return macs
}
