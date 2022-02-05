/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cli

import (
	"github.com/obarbier/awesome-crypto/pkg/user/adapters"
	"github.com/obarbier/awesome-crypto/pkg/user/domain"
	"github.com/obarbier/awesome-crypto/pkg/user/handler"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		interruptChan := make(chan os.Signal, 1)
		mongoRepo, err := adapters.NewMongoRepository()
		if err != nil {
			log.Fatalf("error occurred while creating MongoRepository: %v", err)
			return
		}

		userService := domain.NewUserService(mongoRepo)
		handlerProps := handler.NewPropertiesHandler(userService, false)

		err = http.ListenAndServe("0.0.0.0:2023", handler.Handler(handlerProps))
		if err != nil {
			return
		}

		signal.Notify(interruptChan, syscall.SIGINT, syscall.SIGTERM)

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
