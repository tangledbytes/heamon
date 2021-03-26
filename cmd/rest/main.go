package main

import (
	"github.com/utkarsh-pro/heamon/pkg/entrypoint/rest"
	"github.com/utkarsh-pro/heamon/pkg/utils"
)

var (
	version = "dev"
	commit  = "none"
	date    = "NA"
)

func main() {
	utils.PrintInfo(version, commit, date)
	rest.Run()
}
