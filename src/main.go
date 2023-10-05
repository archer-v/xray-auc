package main

import (
	"encoding/json"
	"fmt"
	"os"
)
import "github.com/integrii/flaggy"
import "github.com/go-playground/validator/v10"
import "github.com/vrischmann/envconfig"

type config struct {
	UserRecord
	Host string `envconfig:"optional"`
}

const (
	ExitOk             = 0
	ExitPartialDone    = 1
	ExitErrUnavailable = 2
)

var gitTag, gitCommit, gitBranch, build, buildTimestamp, version string

func main() {

	if build == "" {
		build = "DEV"
		version = fmt.Sprintf("version: DEV, build: %v", build)
	} else {
		version = fmt.Sprintf("version: %v-%v-%v, build: %v,%v", gitTag, gitBranch, gitCommit, build, buildTimestamp)
	}

	var c config
	var f string

	//parse environment variables to the config struct
	if err := envconfig.InitWithPrefix(&c, "XRAY_C"); err != nil {
		panic(err)
	}

	flaggy.DefaultParser.ShowHelpOnUnexpected = true

	//flaggy.SetDescription("ddd")

	addUserOp := flaggy.NewSubcommand("addUser")
	addUserOp.Description = "add user to an inbound proxy configuration"
	addUserOp.String(&c.Proto, "p", "proto", "proxy protocol, one of: vmess, vless, trojan, shadowsocks")
	addUserOp.String(&c.Email, "e", "email", "user email (is used as a human id)")
	addUserOp.String(&c.Password, "s", "password", "user secret (password for shadowsocks or id for vless/vmess/trojan)")
	addUserOp.String(&c.Flow, "", "flow", "flow (vless proto only)")
	addUserOp.String(&c.Cipher, "c", "cipher", "cipher (shadowsocks proto only, optional)")

	delUserOp := flaggy.NewSubcommand("rmUser")
	delUserOp.Description = "remove user from an inbound proxy configuration"
	delUserOp.String(&c.Email, "e", "email", "user email (is used as human id)")

	flaggy.String(&c.Host, "a", "addr", "xray server host and port separated with a colon")
	flaggy.String(&c.Tag, "t", "tag", "proxy tag")
	flaggy.String(&f, "f", "file", "filepath to json file with array of user records in format: [ { user_options...}, .... {}]")

	flaggy.AttachSubcommand(addUserOp, 1)
	flaggy.AttachSubcommand(delUserOp, 1)

	flaggy.Parse()

	err := validator.New().Var(c.Host, "required,hostname_port")
	if err != nil {
		fmt.Printf("Host option should be defined as host:port")
		return
	}

	if !addUserOp.Used && !delUserOp.Used {
		flaggy.ShowHelp("Also environment variables could be used to configure program options. Use env variables with prefix XRAY_C, for example XRAY_C_HOST")
		return
	}

	var records []UserRecord

	if f != "" {
		content, err := os.ReadFile(f)
		if err != nil {
			fmt.Printf("Error when opening file: %v", err)
			return
		}

		// Now let's unmarshall the data into `payload`
		var payload []UserRecord
		err = json.Unmarshal(content, &payload)
		if err != nil {
			fmt.Printf("Error during JSON file parsing: %v", err)
			return
		}
		records = append(records, payload...)
	} else {
		if c.Email == "" {
			fmt.Printf("Email can't be empty")
			return
		}
		records = append(records, c.UserRecord)
	}

	var o Operation
	if addUserOp.Used {
		o = OpAdd
	} else if delUserOp.Used {
		o = OpRemove
	}

	done, err := ApplyToXray(c.Host, o, records)

	if err != nil {
		fmt.Println(err)
		if done == 0 {
			os.Exit(ExitErrUnavailable)
		}
	}

	if done == len(records) {
		fmt.Printf("All done\n")
		os.Exit(ExitOk)
	}

	fmt.Printf("Only %v from %v records applied", done, len(records))
	os.Exit(ExitPartialDone)
}
