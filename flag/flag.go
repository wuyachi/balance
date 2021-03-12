package flag

import (
	"github.com/urfave/cli"
)

// ConfFlag ...
var ConfFlag = cli.StringFlag{
	Name:  "conf",
	Usage: "specify config file",
	Value: "config.json",
}
