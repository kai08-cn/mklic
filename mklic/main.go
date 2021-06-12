package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	license "github.com/zfs123/mklic"
)

func main() {
	app := cli.NewApp()
	app.Name = "mklic"
	app.Usage = "A license generate and validate tool"

	app.Commands = append(app.Commands,
		&cli.Command{
			Name:    "sign",
			Aliases: []string{"s"},
			Usage:   "Create a license",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "prikey", Usage: "private key file"},
				&cli.StringFlag{Name: "data", Usage: "data file(json format)"},
				&cli.StringFlag{Name: "devid", Usage: "device id(if not defined then local)"},
				&cli.StringFlag{Name: "output", Usage: "output file(if not defined then stdout)"},
			},
			Action: licgen,
		},
		&cli.Command{
			Name:    "verify",
			Aliases: []string{"v"},
			Usage:   "Verify a license",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "lic", Usage: "license file", Required: true},
				&cli.StringFlag{Name: "pubkey", Usage: "public key"},
				&cli.StringFlag{Name: "devid", Usage: "device id(if not defined then local)"},
			},
			Action: licver,
		},
		&cli.Command{
			Name:    "gen",
			Aliases: []string{"g"},
			Usage:   "Generate a device id",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "macs", Usage: "macs for generate device id(if not defined then local)"},
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
	fp := c.String("prikey")
	prikey, err := ioutil.ReadFile(fp)
	if err != nil {
		log.Fatal(err)
	}
	fp = c.String("data")
	var data []byte
	if fp != "" {
		data, err = ioutil.ReadFile(fp)
		if err != nil {
			log.Fatal(err)
		}
	}
	devid := c.String("devid")
	lic, err := license.Sign(prikey, data, devid)

	output := c.String("output")
	if output != "" {
		if err := ioutil.WriteFile(output, []byte(lic), 0600); err != nil {
			log.Fatal(err)
		}
	} else {
		if _, err := os.Stdout.WriteString(lic); err != nil {
			log.Fatal(err)
		}
	}
	return err
}

func licver(c *cli.Context) error {
	licfile := c.String("lic")
	lic, err := ioutil.ReadFile(licfile)
	if err != nil {
		log.Fatal(err)
	}
	pkfile := c.String("pubkey")
	pk, err := ioutil.ReadFile(pkfile)
	if err != nil {
		log.Fatal(err)
	}
	devid := c.String("devid")
	content, err := license.Verify(lic, pk, devid)
	if err != nil {
		log.Fatal("verify license failed, " + err.Error())
	}

	output := c.String("output")
	if output != "" {
		if err := ioutil.WriteFile(output, content, 0600); err != nil {
			log.Fatal(err)
		}
	} else {
		if _, err := os.Stdout.WriteString(string(content)); err != nil {
			log.Fatal(err)
		}
	}
	return err
}

func genid(c *cli.Context) error {
	macs := c.String("macs")
	fmt.Println(string(license.GenDevId(macs)))
	return nil
}
