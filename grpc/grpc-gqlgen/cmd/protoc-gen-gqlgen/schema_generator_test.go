package main

import (
	"testing"

	"github.com/k0kubun/pp"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

func Test_schemaGeneratorExample(t *testing.T) {
	schema := `
		extend type Query {
		  todosA(first: Int, after: String, input: TodoListAInput!): TodoConnection!
		  todosB(first: Int, after: String, input: TodoListBInput!): TodoConnection!
		}
		
		extend type Mutation {
		  createTodo(input: CreateTodoInput!): CreateTodoPayload!
		  updateTodo(input: UpdateTodoInput!): UpdateTodoPayload!
		}
		
		type Todo {
		  id: ID!
		  text: String!
		  done: Boolean!
		  doneAt: Time
		  updatedAt: Time
		  createdAt: Time
		}
		
		input CreateTodoInput {
		  text: String!
		}
		
		type CreateTodoPayload {
		  todo: Todo!
		}
		
		input UpdateTodoInput {
		  id: ID!
		  text: String
		  done: Boolean
		}
		
		type UpdateTodoPayload {
		  todo: Todo!
		}
		
		input TodoListAInput {
		  notDone: Boolean
		}
		
		input TodoListBInput {
		  notDone: Boolean
		}
		
		type TodoConnection {
		  pageInfo: PageInfo!
		  edges: [TodoEdge!]!
		  nodes: [Todo!]!
		}
		
		type TodoEdge {
		  cursor: String
		  node: Todo!
		}
	`

	doc, gqlErr := parser.ParseSchema(&ast.Source{
		Name:    "todo_schema.graphql",
		Input:   schema,
		BuiltIn: false,
	})
	if gqlErr != nil {
		t.Fatal(gqlErr)
	}

	t.Log(pp.Sprint(doc))
}
