package luasocket

import (
  "net"
  "time"

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
  proto := L.CheckString(1)
  address := L.CheckString(2)
  timeout := L.CheckAny(3)
  var t time.Duration

  switch lv := timeout.(type) {
    case *lua.LNilType:
      t = time.Duration(60)
    case lua.LNumber:
      t = time.Duration(int64(lua.LVAsNumber(lv)))
  }

  dialer := net.Dialer{Timeout: t * time.Second}
  c, err := dialer.Dial(proto, address)
  if err != nil {
    L.Push(lua.LNil)
    L.Push(lua.LString(err.Error()))
    return 2
  }

  var _s = map[string]lua.LGFunction{
    "close": func(L *lua.LState) int {
      c.Close()
      return 0
    },
  }

  ms := L.SetFuncs(L.NewTable(), _s)
  L.Push(ms)
  return 1
}
