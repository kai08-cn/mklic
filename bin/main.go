package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/zfs123/lictool"
)

func main() {
	app := cli.NewApp()
	app.Name = "licgv"
	app.Usage = "A license generate and validate tool"

	app.Commands = append(app.Commands,
		&cli.Command{
			Name:    "sign",
			Aliases: []string{"s"},
			Usage:   "Create a license",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "private_key", Usage: "private key file"},
				&cli.StringFlag{Name: "data", Usage: "data file"},
				&cli.StringFlag{Name: "device_id", Usage: "device id"},
			},
			Action: licgen,
		},
		&cli.Command{
			Name:    "verify",
			Aliases: []string{"v"},
			Usage:   "Verify a license",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "license", Usage: "license file", Required: true},
				&cli.StringFlag{Name: "public_key", Usage: "public key"},
				&cli.StringFlag{Name: "device_id", Usage: "device id"},
			},
			Action: licver,
		},
		&cli.Command{
			Name:    "getid",
			Aliases: []string{"g"},
			Usage:   "Generate a device id",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "macs", Usage: "macs for generate device id"},
			},
			Action: genid,
		},
	)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

func licgen(c *cli.Context) error {
	fp := c.String("private_key")
	prikey, err := ioutil.ReadFile(fp)
	if err != nil {
		log.Fatal(err)
	}
	fp = c.String("data")
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		log.Fatal(err)
	}
	devid := c.String("device_id")
	lic, err := lictool.Sign(prikey, data, devid)
	fmt.Println(lic)
	return err
}

func licver(c *cli.Context) error {
	licfile := c.String("license")
	lic, err := ioutil.ReadFile(licfile)
	if err != nil {
		log.Fatal(err)
	}
	pkfile := c.String("public_key")
	pk, err := ioutil.ReadFile(pkfile)
	if err != nil {
		log.Fatal(err)
	}

	devid := c.String("device_id")
	content, err := lictool.Verify(lic, pk, devid)
	fmt.Println(string(content))
	return err
}

func genid(c *cli.Context) error {
	macs := c.String("macs")
	fmt.Println(string(lictool.GenDevId(macs)))
	return nil
}
