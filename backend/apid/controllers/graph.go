package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/sensu/sensu-go/backend/store"
)

// GraphController defines the fields required by GraphController.
type GraphController struct {
	Store store.Store
}

// Register should define an association between HTTP routes and their
// respective handlers defined within this Controller.
func (c *GraphController) Register(r *mux.Router) {
	r.HandleFunc("/graphql", c.query).Methods(http.MethodGet)
}

// many handles requests to /info
func (c *GraphController) query(w http.ResponseWriter, r *http.Request) {
	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		logger.WithError(err).Fatal("failed to create new schema")
	}

	query := r.URL.Query().Get("query")
	params := graphql.Params{Schema: schema, RequestString: query}
	logger.WithField("query", query).Info("got params")

	res := graphql.Do(params)
	if len(res.Errors) > 0 {
		logger.
			WithField("errors", res.Errors).
			Errorf("failed to execute graphql operation")
	}

	rJSON, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", rJSON)
}
