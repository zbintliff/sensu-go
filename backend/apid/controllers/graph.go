package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	graphqlast "github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/relay"
	"github.com/sensu/sensu-go/backend/store"
	"github.com/sensu/sensu-go/types"
	"golang.org/x/net/context"
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
	ctx := r.Context()
	ctx = context.WithValue(ctx, types.OrganizationKey, "*")
	ctx = context.WithValue(ctx, types.EnvironmentKey, "*")
	ctx = context.WithValue(ctx, types.StoreKey, c.Store)

	queryStr := r.URL.Query().Get("query")
	res := execQuery(ctx, queryStr)
	if len(res.Errors) > 0 {
		logger.
			WithField("errors", res.Errors).
			Errorf("failed to execute graphql operation")

		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	rJSON, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", rJSON)
}

func execQuery(ctx context.Context, queryStr string) *graphql.Result {
	params := graphql.Params{
		Schema:        GraphqlSchema,
		RequestString: queryStr,
		Context:       ctx,
	}

	if logger.Logger.Level >= logrus.DebugLevel {
		re := regexp.MustCompile(`\s+`) // TODO move to init
		formattedQuery := queryStr
		formattedQuery = strings.Replace(formattedQuery, "\n", " ", -1)
		formattedQuery = re.ReplaceAllLiteralString(formattedQuery, " ")
		logger.WithField("query", formattedQuery).Debug("executing GraphQL query")
	}

	return graphql.Do(params)
}

// GraphqlSchema ...
var GraphqlSchema graphql.Schema

func init() {
	var err error
	var checkEventType *graphql.Object
	var metricEventType *graphql.Object
	var entityType *graphql.Object

	//
	// Relay

	nodeDefinitions := relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
		IDFetcher: func(id string, info graphql.ResolveInfo, ctx context.Context) (interface{}, error) {
			// resolve id from global id
			resolvedID := relay.FromGlobalID(id)

			// based on id and its type, return the object
			switch resolvedID.Type {
			case "CheckEvent":
				// TODO
				return types.FixtureEvent("a", "b"), nil
			case "MetricEvent":
				// TODO
				return types.FixtureEvent("a", "b"), nil
			case "Entity":
				// TODO
				return types.FixtureEntity("b"), nil
			default:
				return nil, errors.New("Unknown node type")
			}
		},
		TypeResolve: func(p graphql.ResolveTypeParams) *graphql.Object {
			// based on the type of the value, return GraphQLObjectType
			switch p.Value.(type) {
			case *types.Event:
				return checkEventType
			case *types.Entity:
				return entityType
			default:
				return nil
			}
		},
	})

	//
	// Interface

	multitenantResource := graphql.NewInterface(graphql.InterfaceConfig{
		Name:        "MultitenantResource",
		Description: "A resource that belong to an organization and environment",
		Fields: graphql.Fields{
			"environment": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The environment the resource belongs to.",
			},
			"organization": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The organization the resource belongs to.",
			},
		},
		ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
			if _, ok := p.Value.(types.Entity); ok {
				return entityType
			}
			return nil
		},
	})

	//
	// Test Scalar

	timeScalar := graphql.NewScalar(graphql.ScalarConfig{
		Name:        "Time",
		Description: "The `Time` scalar type represents an instant in time",
		Serialize:   coerceTime,
		ParseValue:  coerceTime,
		ParseLiteral: func(valueAST graphqlast.Value) interface{} {
			switch valueAST := valueAST.(type) {
			case *graphqlast.IntValue:
				if intValue, err := strconv.Atoi(valueAST.Value); err == nil {
					return time.Unix(int64(intValue), 0)
				}
			case *graphqlast.StringValue:
				// TODO: Would be nice to cover
			}
			return nil
		},
	})

	entityType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Entity",
		Interfaces: []*graphql.Interface{
			nodeDefinitions.NodeInterface,
			multitenantResource,
		},
		Fields: graphql.Fields{
			"id":               relay.GlobalIDField("Entity", nil),
			"class":            &graphql.Field{Type: graphql.String},
			"subscriptions":    &graphql.Field{Type: graphql.NewList(graphql.String)},
			"lastSeen":         &graphql.Field{Type: graphql.String},
			"deregister":       &graphql.Field{Type: graphql.Boolean},
			"keepaliveTimeout": &graphql.Field{Type: graphql.Int},
			"environment": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The environment the resource belongs to.",
			},
			"organization": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The organization the resource belongs to.",
			},
		},
		IsTypeOf: func(p graphql.IsTypeOfParams) bool {
			_, ok := p.Value.(*types.Entity)
			return ok
		},
	})

	metricEventType = graphql.NewObject(graphql.ObjectConfig{
		Name: "MetricEvent",
		Fields: graphql.Fields{
			"id":        relay.GlobalIDField("MetricEvent", nil),
			"entity":    &graphql.Field{Type: entityType},
			"timestamp": &graphql.Field{Type: timeScalar},
		},
		IsTypeOf: func(p graphql.IsTypeOfParams) bool {
			if e, ok := p.Value.(*types.Event); ok {
				return e.Metrics != nil
			}
			return false
		},
	})

	checkEventType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "CheckEvent",
		Description: "A check result",
		Interfaces: []*graphql.Interface{
			nodeDefinitions.NodeInterface,
		},
		Fields: graphql.Fields{
			"id":        relay.GlobalIDField("CheckEvent", nil),
			"timestamp": &graphql.Field{Type: timeScalar},
			"entity":    &graphql.Field{Type: entityType},
			"output":    AliasField(graphql.String, "Check", "Output"),
			"status":    AliasField(graphql.Int, "Check", "Status"),
			"issued":    AliasField(timeScalar, "Check", "Issued"),
			"executed":  AliasField(timeScalar, "Check", "Executed"),
		},
		IsTypeOf: func(p graphql.IsTypeOfParams) bool {
			if e, ok := p.Value.(*types.Event); ok {
				return e.Check != nil
			}
			return false
		},
	})

	//
	// Test Union

	eventType := graphql.NewUnion(graphql.UnionConfig{
		Name:        "Event",
		Description: "???", // TODO: ???
		Types: []*graphql.Object{
			checkEventType,
			metricEventType,
		},
		ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
			// TODO
			if event, ok := p.Value.(*types.Event); ok {
				if event.Check != nil {
					return checkEventType
				} else if event.Metrics != nil {
					return metricEventType
				}
			}

			return nil
		},
	})

	entityConnectionDef := relay.ConnectionDefinitions(relay.ConnectionConfig{
		Name:     "Entity",
		NodeType: entityType,
	})

	viewerType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Viewer",
		Description: "describes resources available to the curr1nt user",
		Fields: graphql.Fields{
			"entities": &graphql.Field{
				Type: entityConnectionDef.ConnectionType,
				Args: relay.ConnectionArgs,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					args := relay.NewConnectionArguments(p.Args)

					entities := []interface{}{
						types.FixtureEntity("a"),
						types.FixtureEntity("b"),
						types.FixtureEntity("c"),
						types.FixtureEntity("d"),
						types.FixtureEntity("e"),
						types.FixtureEntity("f"),
						types.FixtureEntity("a"),
						types.FixtureEntity("a"),
						types.FixtureEntity("a"),
						types.FixtureEntity("a"),
						types.FixtureEntity("a"),
						types.FixtureEntity("a"),
						types.FixtureEntity("a"),
						types.FixtureEntity("a"),
						types.FixtureEntity("a"),
						types.FixtureEntity("a"),
						types.FixtureEntity("a"),
						types.FixtureEntity("b"),
						types.FixtureEntity("c"),
						types.FixtureEntity("d"),
						types.FixtureEntity("e"),
						types.FixtureEntity("f"),
						types.FixtureEntity("b"),
						types.FixtureEntity("c"),
						types.FixtureEntity("d"),
						types.FixtureEntity("e"),
						types.FixtureEntity("f"),
						types.FixtureEntity("b"),
						types.FixtureEntity("c"),
						types.FixtureEntity("d"),
						types.FixtureEntity("e"),
						types.FixtureEntity("f"),
						types.FixtureEntity("b"),
						types.FixtureEntity("c"),
						types.FixtureEntity("d"),
						types.FixtureEntity("e"),
						types.FixtureEntity("f"),
						types.FixtureEntity("b"),
						types.FixtureEntity("c"),
						types.FixtureEntity("d"),
						types.FixtureEntity("e"),
						types.FixtureEntity("f"),
						types.FixtureEntity("b"),
						types.FixtureEntity("c"),
						types.FixtureEntity("d"),
						types.FixtureEntity("e"),
						types.FixtureEntity("f"),
						types.FixtureEntity("b"),
						types.FixtureEntity("c"),
						types.FixtureEntity("d"),
						types.FixtureEntity("e"),
						types.FixtureEntity("f"),
						types.FixtureEntity("b"),
						types.FixtureEntity("c"),
						types.FixtureEntity("d"),
						types.FixtureEntity("e"),
						types.FixtureEntity("f"),
						types.FixtureEntity("b"),
						types.FixtureEntity("c"),
						types.FixtureEntity("d"),
						types.FixtureEntity("e"),
						types.FixtureEntity("f"),
						types.FixtureEntity("b"),
						types.FixtureEntity("c"),
						types.FixtureEntity("d"),
						types.FixtureEntity("e"),
						types.FixtureEntity("f"),
						types.FixtureEntity("b"),
						types.FixtureEntity("c"),
						types.FixtureEntity("d"),
						types.FixtureEntity("e"),
						types.FixtureEntity("f"),
					}

					return relay.ConnectionFromArray(entities, args), nil
				},
			},
			"events": &graphql.Field{
				Type: graphql.NewList(eventType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					events := []interface{}{
						types.FixtureEvent("d", "e"),
						types.FixtureEvent("b", "f"),
						types.FixtureEvent("c", "b"),
						types.FixtureEvent("d", "e"),
						types.FixtureEvent("b", "f"),
						types.FixtureEvent("c", "b"),
						types.FixtureEvent("d", "e"),
						types.FixtureEvent("b", "f"),
						types.FixtureEvent("c", "b"),
						types.FixtureEvent("d", "e"),
					}
					return events, nil
				},
			},
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"node": nodeDefinitions.NodeField,
			"viewer": &graphql.Field{
				Type: viewerType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return 1, nil // TODO? User? Viewer warpper type?
				},
			},
		},
	})
	GraphqlSchema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
		// Mutation: mutationType,
	})
	if err != nil {
		logger.WithError(err).Fatal("failed to create new schema")
	}
}

func coerceTime(value interface{}) interface{} {
	switch value := value.(type) {
	case time.Time:
		return value.Format(time.RFC1123Z)
	case int64: // TODO: Too naive
		return coerceTime(time.Unix(value, 0))
	}

	return nil
}

// AliasField TODO: ...
func AliasField(T graphql.Output, fNames ...string) *graphql.Field {
	return &graphql.Field{
		Type: T,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			fVal := reflect.ValueOf(p.Source)
			for _, fName := range fNames {
				fVal = reflect.Indirect(fVal)
				fVal = fVal.FieldByName(fName)
			}
			return fVal.Interface(), nil
		},
	}
}
