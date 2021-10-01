package models

import (
	"fmt"
	"strings"
)

const (
	DeleteStatement SQLStatement = iota
	InsertStatement
	SelectStatement
	UpdateStatement
)

type SQLStatement int

type SQLResponse struct {
	StatusCode int
	Body       interface{}
}

type SQLOptions struct {
	// Contains additional clauses after the table name e.g ... WHERE id = 1, ... SET name = 'test'
	AfterTableClauses string

	// Values to replace the index parameters in the query e.g $1
	Arguments []interface{}

	// Variable pointers to scan values into when query a single row  e.g &User.Email
	Destination []interface{}

	// List of column names when inserting a row into a table
	InsertColumns []string

	// List of column names when selecting a row from a table or
	// returning a row from an insert, update or delete
	ReturnColumns []string

	Statement SQLStatement

	TableName string
}

func buildQuery(options SQLOptions) string {
	switch options.Statement {
	case DeleteStatement:
		return ""
	case InsertStatement:
		params := []string{}
		for i := 0; i < len(options.InsertColumns); i++ {
			params = append(params, fmt.Sprintf("$%v", i+1))
		}
		paramString := strings.Join(params, ", ")

		returnColumns := ""
		if len(options.ReturnColumns) != 0 {
			returnColumns = "RETURNING "
		}
		returnColumns += strings.Join(options.ReturnColumns, ", ")

		insertColumns := strings.Join(options.InsertColumns, ", ")
		format := "INSERT INTO %v (%v) VALUES (%v) %v %v"
		args := []interface{}{options.TableName, insertColumns, paramString, options.AfterTableClauses, returnColumns}
		return fmt.Sprintf(format, args...)
	case UpdateStatement:
		returnColumns := ""
		if len(options.ReturnColumns) != 0 {
			returnColumns = "RETURNING "
		}
		returnColumns += strings.Join(options.ReturnColumns, ", ")

		format := "UPDATE %v %v %v"
		args := []interface{}{options.TableName, options.AfterTableClauses, returnColumns}
		return fmt.Sprintf(format, args...)
	case SelectStatement:
		selectColumns := strings.Join(options.ReturnColumns, ", ")
		format := "SELECT %v FROM %v %v"
		args := []interface{}{selectColumns, options.TableName, options.AfterTableClauses}
		return fmt.Sprintf(format, args...)
	default:
		panic("invalid sql statement")
	}
}
