package cmd

import (
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"log"
	"net/http"

	"github.com/labbcb/brave/mongo"
	"github.com/labbcb/brave/server"
	"github.com/spf13/cobra"
)

func init() {
	serverCmd.Flags().String("database", "mongodb://localhost:27017", "URL to MongoDB")
	viper.BindPFlag("database", serverCmd.Flags().Lookup("database"))

	serverCmd.Flags().String("address", ":8080", "Address to bind server.")
	viper.BindPFlag("address", serverCmd.Flags().Lookup("address"))

	serverCmd.Flags().String("username", "admin", "User name.")
	viper.BindPFlag("username", serverCmd.Flags().Lookup("username"))

	serverCmd.Flags().String("password", "", "Password.")
	viper.BindPFlag("password", serverCmd.Flags().Lookup("password"))

	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start BraVE server",
	Long:  `BraVE server requires a MongoDB instance to store genomics variants.`,
	Run: func(cmd *cobra.Command, args []string) {
		database := viper.GetString("database")
		db, err := mongo.Connect(database, "brave")
		if err != nil {
			log.Fatalf("Conneting to MongoDB: %v", err)
		}

		username := viper.GetString("username")
		password := viper.GetString("password")
		s := server.New(db, username, password)

		address := viper.GetString("address")
		log.Fatal(http.ListenAndServe(address, cors.Default().Handler(s.Router)))
	},
}
