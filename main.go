package main

import (
	"github.com/idprm/go-yellowclinic/src/cmd"
	"github.com/idprm/go-yellowclinic/src/database"
)

func init() {
	database.Connect()

	// helper.WriteLog()
}

func main() {
	/**
	 * SETUP COBRA
	 */
	cmd.Execute()
}
