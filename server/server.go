package server

import (
	"fmt"
	"strings"
	"github.com/gorilla/mux"
	"github.com/labbcb/brave/mongo"
	"github.com/labbcb/brave/search"
	"github.com/labbcb/brave/variant"
)

// Server contains required dependencies.
type Server struct {
	DB       *mongo.DB
	Router   *mux.Router
	Username string
	Password string
}

// New creates a BraVE server.
func New(mongoClient *mongo.DB, username, password string) *Server {
	s := &Server{
		Router:   mux.NewRouter(),
		DB:       mongoClient,
		Username: username,
		Password: password,
	}
	s.register()
	return s
}

func (s *Server) Search(input *search.Input) (*search.Response, error) {
	return s.DB.Search(input)
}

// InsertVariant generates an ID dataset-assembly-reference-start-ref-alt and saves into database.
func (s *Server) InsertVariant(v *variant.Variant) error {
	v.ID = fmt.Sprintf("%s-%s-%s-%d-%s-%s", v.DatasetID, v.AssemblyID, v.ReferenceName, v.Start, v.ReferenceBases, strings.Join(v.AlternateBases[:], "_"))
	return s.DB.Save(v)
}

func (s *Server) RemoveVariants(datasetID, assemblyID string) error {
	return s.DB.Remove(datasetID, assemblyID)
}
