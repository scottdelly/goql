package schemas

import (
	"strconv"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"

	"github.com/scottdelly/goql/models"
)

var ModelIdScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "ModelId",
	Description: "The `ModelId` scalar type represents an Id Object.",
	// Serialize serializes `CustomID` to string.
	Serialize: func(value interface{}) interface{} {
		switch value := value.(type) {
		case models.ModelId:
			return int(value)
		case *models.ModelId:
			v := *value
			return int(v)
		default:
			return nil
		}
	},
	// ParseValue parses GraphQL variables from `string` to `CustomID`.
	ParseValue: func(value interface{}) interface{} {
		switch value := value.(type) {
		case int:
			return models.ModelId(value)
		case *int:
			return models.ModelId(*value)
		default:
			return nil
		}
	},
	// ParseLiteral parses GraphQL AST value to `CustomID`.
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.IntValue:
			intVal, err := strconv.ParseInt(valueAST.Value, 0, 0)
			if err != nil {
				return err
			}
			return models.ModelId(intVal)
		default:
			return nil
		}
	},
})

var DurationScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Duration",
	Description: "The `Duration` scalar type represents a time duration.",
	// Serialize serializes `CustomID` to string.
	Serialize: func(value interface{}) interface{} {
		switch value := value.(type) {
		case time.Duration:
			return int(value)
		case *time.Duration:
			v := *value
			return int(v)
		default:
			return nil
		}
	},
	// ParseValue parses GraphQL variables from `string` to `CustomID`.
	ParseValue: func(value interface{}) interface{} {
		switch value := value.(type) {
		case int:
			return time.Duration(value)
		case *int:
			return time.Duration(*value)
		default:
			return nil
		}
	},
	// ParseLiteral parses GraphQL AST value to `CustomID`.
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.IntValue:
			return time.Duration(valueAST.GetValue().(int))
		default:
			return nil
		}
	},
})
