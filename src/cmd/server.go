package cmd

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/idprm/go-yellowclinic/src/config"
	"github.com/idprm/go-yellowclinic/src/route"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		app := fiber.New()

		app.Use(cors.New())

		route.Setup(app)

		path, err := os.Getwd()

		if err != nil {
			panic(err)
		}

		app.Static("/", path+"/public")

		log.Fatal(app.Listen(":" + config.ViperEnv("APP_PORT")))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
