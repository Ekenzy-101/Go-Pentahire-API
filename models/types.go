package models

type SQLResponse struct {
	StatusCode int
	Body       interface{}
}

type SQLOption struct {
	// Contains additional clauses after the table name e.g ... WHERE id = 1
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
}
