package luabase64

import (
  "encoding/base64"

  "github.com/yuin/gopher-lua"
)

var m = map[string]lua.LGFunction{
  "decode": decode,
  "encode": encode,
}

func Loader(L *lua.LState) int {
  module := L.SetFuncs(L.NewTable(), m)
  L.Push(module)
  return 1
}

func encode(L *lua.LState) int {
  s := L.CheckString(1)
  b := base64.StdEncoding.EncodeToString([]byte(s))
  L.Push(lua.LString(b))
  return 1
}

func decode(L *lua.LState) int {
  b := L.CheckString(1)
  s, err := base64.StdEncoding.DecodeString(b)
  if err != nil {
    L.Push(lua.LNil)
    L.Push(lua.LString(err.Error()))
    return 2
  }

  L.Push(lua.LString(s))
  return 1
}
