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
  port := L.CheckString(4)
  database := L.CheckString(5)

  db, err := sql.Open("mysql", fmt.Sprintf(
    "%s:%s@tcp(%s:%d)/%s",
    username, password,
    address, port,
    database))
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
      var nargs int
      query := L.CheckString(1)

      for i := 2;; i++ {
        a := L.CheckAny(i)
        switch lv := a.(type) {
          case *lua.LNilType:
            nargs = i - 1
            break;
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
      if nargs > 0 {
        p := make([]interface{}, nargs)
	for i := 0; i < nargs; i++ {
          p[i] = &args[i]
        }

        rows, err = stmt.Query(p...)
      } else {
        rows, err = stmt.Query()
      }

      if err != nil {
        L.Push(lua.LNil)
        L.Push(lua.LString(err.Error()))
        return 2
      }

      cols, _ := rows.Columns()
      var result []lua.LTable
      for rows.Next() {
        columns := make([]interface{}, len(cols))
        cp := make([]interface{}, len(cols))
	for i, _ := range columns {
          cp[i] = &columns[i]
        }

	err := rows.Scan(cp...)
        if err != nil {
          L.Push(lua.LNil)
          L.Push(lua.LString(err.Error()))
          return 2
        }

        t := L.CreateTable(0, len(cols))
        for k, _ := range cols {
          switch lv := cp[k].(type) {
            case bool:
              t.RawSetH(lua.LString(k), lua.LBool(lv))
            case float64:
              t.RawSetH(lua.LString(k), lua.LNumber(lv))
            case int:
              t.RawSetH(lua.LString(k), lua.LNumber(float64(lv)))
            case string:
              t.RawSetH(lua.LString(k), lua.LString(lv))
            case nil:
              t.RawSetH(lua.LString(k), lua.LNil)
          }
        }

        result = append(result, *t)
      }

      tables := L.CreateTable(len(result), 0)
      for _, v := range result {
        tables.Append(&v)
      }

      L.Push(tables)
      return 1
    },
  }

  module := L.SetFuncs(L.NewTable(), d)
  L.Push(module)
  return 1
}
