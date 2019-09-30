package cmd

import (
	"log"
	"net/http"

	"github.com/labbcb/brave/mongo"
	"github.com/labbcb/brave/server"
	"github.com/spf13/cobra"
)

func init() {
	serverCmd.Flags().StringVar(&database, "database", "mongodb://localhost:27017", "URL to MongoDB")
	serverCmd.Flags().StringVar(&address, "address", ":8080", "Address to bind server.")
	serverCmd.Flags().StringVar(&username, "username", "admin", "User name.")
	serverCmd.Flags().StringVar(&password, "password", "", "Password.")
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start BraVE server",
	Long:  `BraVE server requires a MongoDB instance to store genomics variants.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := mongo.Connect(database, "brave")
		if err != nil {
			log.Fatalf("Conneting to MongoDB: %v", err)
		}

		s := server.New(db, username, password)
		log.Fatal(http.ListenAndServe(address, s.Router))
	},
}
