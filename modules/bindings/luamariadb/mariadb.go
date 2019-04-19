package luamariadb

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "fmt"

  "github.com/yuin/gopher-lua"
)

var m = map[string]lua.LGFunction{
  "open": open,
}

func Loader(L *lua.LState) int {
  module := L.SetFuncs(L.NewTable(), m)
  L.Push(module)
  return 1
}

func open(L *lua.LState) int {
  username := L.CheckString(1)
  password := L.CheckString(2)
  address := L.CheckString(3)
  database := L.CheckString(4)

  db, err := sql.Open("mysql", fmt.Sprintf(
    "%s:%s@tcp(%s)/%s",
    username, password,
    address, database))
  if err != nil {
    L.Push(lua.LNil)
    L.Push(lua.LString(err.Error()))
    return 2
  }

  var d = map[string]lua.LGFunction{
    "close": func(L *lua.LState) int {
      db.Close()
      return 0
    },
    "query": func(L *lua.LState) int {
      var args []interface{}
      nargs := L.GetTop()
      query := L.CheckString(1)

      for i := 2; i <= nargs; i++ {
        a := L.CheckAny(i)
        switch lv := a.(type) {
          case *lua.LNilType:
            args = append(args, lua.LNil)
          case lua.LString:
            args = append(args, lv.String())
          case lua.LNumber:
            args = append(args, float64(lua.LVAsNumber(lv)))
        }
      }

      stmt, err := db.Prepare(query)
      if err != nil {
        L.Push(lua.LNil)
	L.Push(lua.LString(err.Error()))
        return 2
      }
      defer stmt.Close()

      var rows *sql.Rows
      if nargs > 1 {
        p := make([]interface{}, nargs - 1)
	for i := 0; i < nargs - 1; i++ {
          p[i] = &args[i]
        }

        rows, err = stmt.Query(p...)
      } else {
        rows, err = stmt.Query()
      }

      if nargs > 1 && err != nil {
        L.Push(lua.LNil)
        L.Push(lua.LString(err.Error()))
        return 2
      }

      cols, _ := rows.ColumnTypes()
      var result []lua.LTable
      for rows.Next() {
        cp := make([]interface{}, len(cols))
	for k, v := range cols {
          ct := v.DatabaseTypeName()
          if ct == "VARCHAR" || ct == "TEXT" || ct == "NVARCHAR" ||
             ct == "TIMESTAMP" {
            cp[k] = new(sql.NullString)
          }
          if ct == "DECIMAL" || ct == "INT" || ct == "BIGINT" ||
             ct == "TINYINT" {
            cp[k] = new(sql.NullInt64)
          }
          if ct == "FLOAT" {
            cp[k] = new(sql.NullFloat64)
          }
          if ct == "BOOL" {
            cp[k] = new(sql.NullBool)
          }
        }

	err := rows.Scan(cp...)
        if err != nil {
          L.Push(lua.LNil)
          L.Push(lua.LString(err.Error()))
          return 2
        }

        t := L.CreateTable(0, len(cols))
	for k, v := range cols {
          switch i := (*(&cp[k])).(type) {
            case *sql.NullString:
              if (*i).Valid {
                t.RawSetH(lua.LString(v.Name()), lua.LString((*i).String))
              } else {
                t.RawSetH(lua.LString(v.Name()), lua.LNil)
              }
            case *sql.NullInt64:
              if (*i).Valid {
                t.RawSetH(lua.LString(v.Name()), lua.LNumber((*i).Int64))
              } else {
                t.RawSetH(lua.LString(v.Name()), lua.LNil)
              }
            case *sql.NullFloat64:
              if (*i).Valid {
                t.RawSetH(lua.LString(v.Name()), lua.LNumber((*i).Float64))
              } else {
                t.RawSetH(lua.LString(v.Name()), lua.LNil)
              }
            case *sql.NullBool:
              if (*i).Valid {
                t.RawSetH(lua.LString(v.Name()), lua.LBool((*i).Bool))
              } else {
                t.RawSetH(lua.LString(v.Name()), lua.LNil)
              }
          }
        }
        result = append(result, *t)
      }

      tables := L.CreateTable(len(result), 0)
      for k, _ := range result {
        tables.Append(&result[k])
      }

      L.Push(tables)
      return 1
    },
  }

  module := L.SetFuncs(L.NewTable(), d)
  L.Push(module)
  return 1
}
