# mklic

[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)[![Build Status](https://github.com/zfs123/mklic/workflows/Go/badge.svg)](https://github.com/zfs123/mklic/actions)[![Go Report Card](https://goreportcard.com/badge/github.com/zfs123/mklic)](https://goreportcard.com/report/github.com/zfs123/mklic)

A simple licensing tool and library for golang to generate a license with device information and verify it. 

The project is an upgrade of [lk](https://github.com/hyperboloide/lk). It provides a convenient way for software publishers to sign license keys with their private key and Information submitted by users of their products. Note that this implementation is quite basic and you can make it more secure with a few modification.

Translations: [English](README.md) | [简体中文](README_ZH.md)

## How does it works?
1. Generate a device id with the MAC address
2. Compute the fingerprint of the input data and some inherent information
3. RSA encryption with private key
4. AES encryption using device id
5. Verifying a license is the reverse process

## Usage

A command line helper `mklic` is also provided to generate deviceid and create licenses. Install it with the following command :
```
go install github.com/zfs123/mklic/mklic
```
See the usage bellow on how to use it :
```
$ mklic --help     
NAME:
   mklic - A license generate and validate tool

USAGE:
   mklic [global options] command [command options] [arguments...]

COMMANDS:
   sign, s    Create a license
   verify, v  Verify a license
   gen, g     Generate a device id
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

## Examples

- Generate a device id
```
$ mklic gen --help
NAME:
   mklic gen - Generate a device id

USAGE:
   mklic gen [command options] [arguments...]

OPTIONS:
   --macs value  macs for generate device id(if not defined then local)
   --help, -h    show help (default: false)
```
A device id  is used to identify a unique device. MAC address can usually be used to generate device id.

If you have two MAC addresses are `08:00:27:36:dd:98` and `08:00:27:36:dd:99`.
```
$ mklic gen --macs="08:00:27:36:dd:9808:00:27:36:dd:99"
R5D4NSCCG392WRE3XA74BHC9CMQ4UT62
```
The device id is `R5D4NSCCG392WRE3XA74BHC9CMQ4UT62`. If the flag `--mac` not set, mklic will use the first two interface of the machine. 

- Creat a license
```
$ mklic sign --help
NAME:
   mklic sign - Create a license

USAGE:
   mklic sign [command options] [arguments...]

OPTIONS:
   --prikey value  private key file
   --data value    data file(json format)
   --devid value   device id(if not defined then local)
   --output value  output file(if not defined then stdout)
   --help, -h      show help (default: false)
```
To sign a license, you should have a secret private key. And in order to bind the machine, mklic will use a device id. 

Data file is the information submitted by users when they apply for permission, which is used to save information or make some behavioral restrictions. mklic accepts a data file in JSON format. 
```
// data file
{
    "name": "test",
    "email": "test@test.com",
    "number": 999,
    "yes": true
}
```
Sign a license:
```
$ mklic sign --prikey="test/prikey" --data="test/datafile.json" --devid="R5D4NSCCG392WRE3XA74BHC9CMQ4UT62"
mkRZliS6SS4SSQyW8N+9S24IAlLBP7jDLaFUnADcmhq+KhhljlcILnStOKdJ+K+PHG6lOe6NiqCWgiJ9odDooVxCZdpDhZw0S6yjc5hVBb8QPwAVZS24UeGV9d15CXbbTqv5a5pXjHqurd4el0Z45E6YRtGtKyuMVRJzIbjgtdP14SC6p7Ahk1X5dfXXCaFPNCG/CeEsAxaNeQwMpp9Plbw5rwnb8mJzkMqZsOGKw+fkD2J96NcS1Y0bg2wjwpA5yUD8fxjiBKbHEQj3pmlCv7iJ7zR3+OLaSAJs5ObuK5EuyJ4RG69aEDXXGIdm2ltqnFTglnNeFO8ftvqILMrnlNg4t2inPx5MRqkaSfB34B1m8ncNFWtV1UmbvTZD+rVSGpHnVFjI+n1WUhvXIdVT/gTJImYLopZEERhdS/ZCW2OujQZXulSaVYPb7SLVHgmjKRrdLooC9sFeqRrxp0VnGyiPRWz3LcsBXVpr3/pVil0fscBDcl26uzphs9VHD1XFCx34zIiNvzBxHmCVw2m4sw==# 
```
- Verify a license
```
$ mklic verify --help
NAME:
   mklic verify - Verify a license

USAGE:
   mklic verify [command options] [arguments...]

OPTIONS:
   --lic value     license file
   --pubkey value  public key
   --devid value   device id(if not defined then local)
   --help, -h      show help (default: false)
```
The public key and device id are used to verify a license, and the signature of the license content is checked.

If the verification passes, the license content will be printed
```
$ mklic verify --pubkey="test/pubkey" --lic="test/licensefile" --devid="R5D4NSCCG392WRE3XA74BHC9CMQ4UT62"
{
    "devid": "R5D4NSCCG392WRE3XA74BHC9CMQ4UT62",
    "email": "test@test.com",
    "name": "test",
    "number": 999,
    "time": "2021-06-14 22:24:28 +0800 CST",
    "uuid": "34b4ea1c-f424-4dfb-8451-7022540ad404",
    "version": "v1.0",
    "yes": true
}
```

## Library

Get expansion pack:
```
$ go get github.com/zfs123/mklic
```
Usage examples:
```
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/zfs123/mklic"
)

const priKey = ""
const pubKey = ""
const macs = ""

var data = struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Number int    `json:"number"`
	Yes    bool   `json:"yes"`
}{
	"test",
	"test@test.com",
	999,
	true,
}

func main() {
	devid := string(mklic.GenDevId(macs))
	jd, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	lic, err := mklic.Sign([]byte(priKey), jd, devid)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(lic)

	cleartext, err := mklic.Verify([]byte(lic), []byte(pubKey), devid)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cleartext)
}
```