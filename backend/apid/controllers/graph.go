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
	// fields := graphql.Fields{
	// 	"hello": &graphql.Field{
	// 		Type: graphql.String,
	// 		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
	// 			return "world", nil
	// 		},
	// 	},
	// }
	// rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	// schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	// schema, err := graphql.NewSchema(schemaConfig)
	// if err != nil {
	// 	logger.WithError(err).Fatal("failed to create new schema")
	// }

	params := graphql.Params{
		Schema:        graphqlSchema,
		RequestString: r.URL.Query().Get("query"),
	}
	logger.WithField("query", params).Info("Received GraphQL query")

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

var graphqlSchema graphql.Schema

func init() {
	var err error

	viewerObj := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Viewer",
		Description: "describes resources available to the currently authorized user",
		Fields: graphql.Fields{
			"events": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "werd", nil
				},
			},
		},
	})

	rootObj := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"viewer": &graphql.Field{
				Type: viewerObj,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, nil // TODO? User? Viewer warpper type?
				},
			},
		},
	})
	graphqlSchema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: rootObj,
	})
	if err != nil {
		logger.WithError(err).Fatal("failed to create new schema")
	}
}
