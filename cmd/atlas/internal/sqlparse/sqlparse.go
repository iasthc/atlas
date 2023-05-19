// Copyright 2021-present The Atlas Authors. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package sqlparse

import (
	"sync"

	"github.com/iasthc/atlas/cmd/atlas/internal/sqlparse/myparse"
	"github.com/iasthc/atlas/cmd/atlas/internal/sqlparse/pgparse"
	"github.com/iasthc/atlas/cmd/atlas/internal/sqlparse/sqliteparse"
	"github.com/iasthc/atlas/sql/migrate"
	"github.com/iasthc/atlas/sql/mysql"
	"github.com/iasthc/atlas/sql/postgres"
	"github.com/iasthc/atlas/sql/schema"
	"github.com/iasthc/atlas/sql/sqlite"
)

// A Parser represents an SQL file parser used to fix, search and enrich schema.Changes.
type Parser interface {
	// FixChange fixes the changes according to the given statement.
	FixChange(d migrate.Driver, stmt string, changes schema.Changes) (schema.Changes, error)

	// ColumnFilledBefore checks if the column was filled with values before the given position
	// in the file. For example:
	//
	//	UPDATE <table> SET <column> = <value>
	//	UPDATE <table> SET <column> = <value> WHERE <column> IS NULL
	//
	ColumnFilledBefore(migrate.File, *schema.Table, *schema.Column, int) (bool, error)
}

// drivers specific fixers.
var drivers sync.Map

// Register a fixer with the given name.
func Register(name string, f Parser) {
	drivers.Store(name, f)
}

// ParserFor returns a ChangesFixer for the given driver.
func ParserFor(name string) Parser {
	f, ok := drivers.Load(name)
	if ok {
		return f.(Parser)
	}
	return nil
}

func init() {
	Register(mysql.DriverName, &myparse.Parser{})
	Register(postgres.DriverName, &pgparse.Parser{})
	Register(sqlite.DriverName, &sqliteparse.FileParser{})
}
