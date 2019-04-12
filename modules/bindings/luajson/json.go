package luajson

import (
  "encoding/json"

  "github.com/yuin/gopher-lua"
)

var m = map[string]lua.LGFunction{
  "decode": decode,
}

func Loader(L *lua.LState) int {
  module := L.SetFuncs(L.NewTable(), m)
  L.Push(module)
  return 1
}

func Map2LValue(L *lua.LState, obj interface{}) lua.LValue {
  switch lv := obj.(type) {
    case bool:
      return lua.LBool(lv)
    case float64:
      return lua.LNumber(lv)
    case string:
      return lua.LString(lv)
    case []interface{}:
      t := L.CreateTable(len(lv), 0)
      for _, v := range lv {
        t.Append(Map2LValue(L, v))
      }
      return t
    case map[string]interface{}:
      t := L.CreateTable(0, len(lv))
      for k, v := range lv {
        t.RawSetH(lua.LString(k), Map2LValue(L, v))
      }
      return t
    case nil:
      return lua.LNil
  }

  return lua.LNil
}

func decode(L *lua.LState) int {
  jsonString := L.CheckString(1)
  var j map[string]interface{}

  err := json.Unmarshal([]byte(jsonString), &j)
  if err != nil {
    L.Push(lua.LNil)
    L.Push(lua.LString(err.Error()))
    return 2
  }

  lv := Map2LValue(L, j)
  L.Push(lv)

  return 1
}
