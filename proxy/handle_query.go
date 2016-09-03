package proxy

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/berkaroad/saashard/net/mysql"
	"github.com/berkaroad/saashard/sqlparser"

	"github.com/berkaroad/saashard/utils/golog"
)

func (c *ClientConn) handleQuery(sql string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			golog.OutputSql("Error", "err:%v,sql:%s", e, sql)

			if err, ok := e.(error); ok {
				const size = 4096
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]

				golog.Error("ClientConn", "handleQuery",
					err.Error(), 0,
					"stack", string(buf), "sql", sql)
			}
			return
		}
	}()
	sql = strings.TrimSpace(sql)
	println(sql)
	var stmt sqlparser.Statement
	stmt, err = sqlparser.Parse(sql)
	if err != nil {
		golog.Error("proxy", "handleQuery", err.Error(), 0, "sql", sql)
		err = nil
	}

	result := new(mysql.Result)
	result.Status = mysql.SERVER_STATUS_AUTOCOMMIT

	switch v := stmt.(type) {
	case *sqlparser.SimpleSelect:
		println("stmt's type is SimpleSelect")
		switch strings.ToUpper(sqlparser.String(v.SelectExprs[0])) {
		case "CURRENT_USER()":
			result.Resultset = new(mysql.Resultset)
			result.Resultset.Fields = make([]*mysql.Field, 1)
			result.Resultset.Fields[0] = &mysql.Field{Schema: []byte(""),
				Table:        []byte(""),
				OrgTable:     []byte(""),
				Name:         []byte("current_user()"),
				OrgName:      []byte(""),
				Charset:      uint16(mysql.DEFAULT_COLLATION_ID),
				ColumnLength: 423,
				ColumnType:   mysql.MYSQL_TYPE_VAR_STRING,
				Flags:        mysql.NOT_NULL_FLAG,
				Decimals:     31}

			result.Rows = make([]*mysql.Row, 1)
			row := mysql.NewTextRow(result.Resultset.Fields)
			row.AppendStringValue(c.user)
			result.Rows[0] = row
			return c.pkg.WriteResultSet(c.capability, c.status, result)

		case "VERSION()":
			result.Resultset = new(mysql.Resultset)
			result.Resultset.Fields = make([]*mysql.Field, 1)
			result.Resultset.Fields[0] = &mysql.Field{Schema: []byte(""),
				Table:        []byte(""),
				OrgTable:     []byte(""),
				Name:         []byte("version()"),
				OrgName:      []byte(""),
				Charset:      uint16(mysql.DEFAULT_COLLATION_ID),
				ColumnLength: 72,
				ColumnType:   mysql.MYSQL_TYPE_VAR_STRING,
				Flags:        mysql.NOT_NULL_FLAG,
				Decimals:     31}

			result.Rows = make([]*mysql.Row, 1)
			row := mysql.NewTextRow(result.Resultset.Fields)
			row.AppendStringValue(mysql.ServerVersion)
			result.Rows[0] = row
			return c.pkg.WriteResultSet(c.capability, c.status, result)

		case "CONNECTION_ID()":
			result.Resultset = new(mysql.Resultset)
			result.Resultset.Fields = make([]*mysql.Field, 1)
			result.Resultset.Fields[0] = &mysql.Field{Schema: []byte(""),
				Table:        []byte(""),
				OrgTable:     []byte(""),
				Name:         []byte("CONNECTION_ID()"),
				OrgName:      []byte(""),
				Charset:      uint16(mysql.DEFAULT_COLLATION_ID),
				ColumnLength: 10,
				ColumnType:   mysql.MYSQL_TYPE_LONGLONG,
				Flags:        mysql.NOT_NULL_FLAG | mysql.BINARY_FLAG,
				Decimals:     0}

			result.Rows = make([]*mysql.Row, 1)
			row := mysql.NewTextRow(result.Resultset.Fields)
			row.AppendUIntValue(uint64(c.connectionID))
			result.Rows[0] = row
			return c.pkg.WriteResultSet(c.capability, c.status, result)

		default:
			println("select.default=", strings.ToUpper(sqlparser.String(v.SelectExprs[0])))
		}

	case *sqlparser.Show:
		switch strings.ToUpper(v.Section) {
		case "DATABASES":
			result.Resultset = new(mysql.Resultset)
			result.Resultset.Fields = make([]*mysql.Field, 1)
			result.Resultset.Fields[0] = &mysql.Field{Schema: []byte("information_schema"),
				Table:        []byte("SCHEMATA"),
				OrgTable:     []byte("SCHEMATA"),
				Name:         []byte("Database"),
				OrgName:      []byte("SCHEMA_NAME"),
				Charset:      uint16(mysql.DEFAULT_COLLATION_ID),
				ColumnLength: 192,
				ColumnType:   mysql.MYSQL_TYPE_VAR_STRING,
				Flags:        mysql.NOT_NULL_FLAG,
				Decimals:     0}

			result.Rows = make([]*mysql.Row, 0, len(c.schemas))
			for name := range c.schemas {
				row := mysql.NewTextRow(result.Resultset.Fields)
				row.AppendStringValue(name)
				result.Rows = append(result.Rows, row)
			}
			return c.pkg.WriteResultSet(c.capability, c.status, result)

		case "TABLES":
			if isMatched, err := c.ShowTables(sql); isMatched {
				return err
			}

		default:
			if isMatched, err := c.ShowVariables(sql); isMatched {
				return err
			}
			if isMatched, err := c.ShowRoutines(sql); isMatched {
				return err
			}
			if isMatched, err := c.ShowStatus(sql); isMatched {
				return err
			}
			if isMatched, err := c.ShowTables(sql); isMatched {
				return err
			}
			if isMatched, err := c.ShowFullTables(sql); isMatched {
				return err
			}
			if isMatched, err := c.ShowTables(sql); isMatched {
				return err
			}
			if isMatched, err := c.ShowColumns(sql); isMatched {
				return err
			}
			if isMatched, err := c.ShowFullColumns(sql); isMatched {
				return err
			}
			if isMatched, err := c.ShowIndexes(sql); isMatched {
				return err
			}
			return c.pkg.WriteError(c.capability, mysql.NewError(mysql.ER_TABLEACCESS_DENIED_ERROR, fmt.Sprintf("SELECT command denied to user '%s'@'%v' for '%s'", c.user, c.clientIP, sql)))
		}

	case *sqlparser.Set:
		return c.pkg.WriteOK(c.capability, c.status, result)
	}

	if strings.HasPrefix(strings.ToUpper(sql), "SHOW") ||
		strings.HasPrefix(strings.ToUpper(sql), "SELECT") ||
		strings.HasPrefix(strings.ToUpper(sql), "EXPLAIN") {

		println(sql)
		return c.pkg.WriteError(c.capability, mysql.NewError(mysql.ER_TABLEACCESS_DENIED_ERROR, fmt.Sprintf("SELECT command denied to user '%s'@'%v' for '%s'", c.user, c.clientIP, sql)))
	}
	return c.pkg.WriteOK(c.capability, c.status, result)

}
