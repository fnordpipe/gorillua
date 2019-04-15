package luajson

import (
  "encoding/json"

  "github.com/yuin/gopher-lua"
)

var m = map[string]lua.LGFunction{
  "encode": encode,
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

func LValue2Map(obj lua.LValue) interface{} {
  switch lv := obj.(type) {
    case lua.LBool:
      return bool(lv)
    case lua.LNumber:
      if float64(lv) == float64(int64(lv)) {
        return int64(lv)
      } else {
        return float64(lv)
      }
    case lua.LString:
      return string(lv)
    case *lua.LNilType:
      return "null"
    case *lua.LTable:
      m := make(map[string]interface{})
      s := make([]interface{}, 0, lv.Len())

      lv.ForEach(func(k, v lua.LValue) {
        switch k.Type() {
          case lua.LTNumber:
            s = append(s, LValue2Map(v))
          case lua.LTString:
            m[k.String()] = LValue2Map(v)
        }
      })

      if len(m) > 0 {
        return m
      }

      return s
  }

  return nil
}

func encode(L *lua.LState) int {
  jsonTable := L.CheckTable(1)
  j := LValue2Map(jsonTable)

  jsonString, err := json.Marshal(j)
  if err != nil {
    L.Push(lua.LNil)
    L.Push(lua.LString(err.Error()))
    return 2
  }

  L.Push(lua.LString(jsonString))
  return 1
}

func decode(L *lua.LState) int {
  jsonString := L.CheckString(1)
  var m map[string]interface{}
  var s []interface{}

  err := json.Unmarshal([]byte(jsonString), &m)
  if err != nil {
    err = json.Unmarshal([]byte(jsonString), &s)
  }
  if err != nil {
    L.Push(lua.LNil)
    L.Push(lua.LString(err.Error()))
    return 2
  }

  var lv lua.LValue
  if len(m) > 0 {
    lv = Map2LValue(L, m)
  } else {
    lv = Map2LValue(L, s)
  }
  L.Push(lv)

  return 1
}
