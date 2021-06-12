# mklic

[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)[![Build Status](https://github.com/zfs123/mklic/workflows/Go/badge.svg)](https://github.com/zfs123/mklic/actions)[![Go Report Card](https://goreportcard.com/badge/github.com/zfs123/mklic)](https://goreportcard.com/report/github.com/zfs123/mklic)

一个简单的支持设备绑定的许可制作和验证工具，也可以作为 golang 库来使用。 

本项目参考 [lk](https://github.com/hyperboloide/lk)，它为软件开发者提供一个生成和管理许可的思路。注意本项目只有较为基础的实现，开发者可以通过少量的修改使你的许可管理更加安全。


## 工作原理
1. 使用机器的MAC地址生成设备id
2. 计算输入数据加上一些许可固有数据的指纹，用于验证时比对
3. 使用私钥进行RSA加密
4. 使用设备id进行AES加密
5. 验证许可是相反的过程

## Usage

本项目提供一个命令行工具 `mklic`，方便创建和验证许可，可通过以下指令安装：
```
go install github.com/zfs123/mklic/mklic
```
查看使用：
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

## 使用样例

- 生成一个设备 id
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
设备 id 可用于标识唯一的设备。MAC 地址通常用于生成设备 id.

假如你的机器有两个网卡地址分别为 `08:00:27:36:dd:98` 和 `08:00:27:36:dd:99`.
```
$ mklic gen --macs="08:00:27:36:dd:9808:00:27:36:dd:99"
R5D4NSCCG392WRE3XA74BHC9CMQ4UT62
```
得到的设备 id 是 `R5D4NSCCG392WRE3XA74BHC9CMQ4UT62`. 假如 `--mac` 参数为空，则默认使用当前机器的前两个网卡地址。

- 创建一个许可
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
首先你应该有一个保密的私钥和一个设备 id 用于表明该许可将使用的设备。

数据文件是用户申请许可时提交的信息，用于保存信息或做一些行为上的限制。mklic 接收 json 格式的输入数据。

```
// data file
{
    "name": "test",
    "email": "test@test.com",
    "number": 999,
    "yes": true
}
```
签发许可：
```
$ mklic sign --prikey="test/prikey" --data="test/datafile.json" --devid="R5D4NSCCG392WRE3XA74BHC9CMQ4UT62"
mkRZliS6SS4SSQyW8N+9S24IAlLBP7jDLaFUnADcmhq+KhhljlcILnStOKdJ+K+PHG6lOe6NiqCWgiJ9odDooVxCZdpDhZw0S6yjc5hVBb8QPwAVZS24UeGV9d15CXbbTqv5a5pXjHqurd4el0Z45E6YRtGtKyuMVRJzIbjgtdP14SC6p7Ahk1X5dfXXCaFPNCG/CeEsAxaNeQwMpp9Plbw5rwnb8mJzkMqZsOGKw+fkD2J96NcS1Y0bg2wjwpA5yUD8fxjiBKbHEQj3pmlCv7iJ7zR3+OLaSAJs5ObuK5EuyJ4RG69aEDXXGIdm2ltqnFTglnNeFO8ftvqILMrnlNg4t2inPx5MRqkaSfB34B1m8ncNFWtV1UmbvTZD+rVSGpHnVFjI+n1WUhvXIdVT/gTJImYLopZEERhdS/ZCW2OujQZXulSaVYPb7SLVHgmjKRrdLooC9sFeqRrxp0VnGyiPRWz3LcsBXVpr3/pVil0fscBDcl26uzphs9VHD1XFCx34zIiNvzBxHmCVw2m4sw==# 
```
- 验证一个许可
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
公钥和设备id用于验证许可，在验证时会比对内部数据的指纹。

假如验证通过，则打印许可内容：
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

获取库：
```
$ go get github.com/zfs123/mklic
```
使用例子：
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