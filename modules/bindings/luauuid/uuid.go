package luauuid

import (
  "github.com/google/uuid"
  "github.com/yuin/gopher-lua"
)

var m = map[string]lua.LGFunction{
  "create": create,
}

func Loader(L *lua.LState) int {
  module := L.SetFuncs(L.NewTable(), m)
  L.Push(module)
  return 1
}

func create(L *lua.LState) int {
  t, err := uuid.NewRandom()
  if err != nil {
    L.Push(lua.LNil)
    L.Push(lua.LString(err.Error()))
    return 2
  }

  token := t.String()
  L.Push(lua.LString(token))
  return 1
}
