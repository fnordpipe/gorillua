package luacron

import (
  "time"

  "github.com/yuin/gopher-lua"
  "metagit.org/fnordpipe/gorillua/modules/logger"
)

var m = map[string]lua.LGFunction{
  "run": run,
}

func Loader(L *lua.LState) int {
  module := L.SetFuncs(L.NewTable(), m)
  L.Push(module)
  return 1
}

func run(L *lua.LState) int {
  interval := L.CheckNumber(1)
  function := L.CheckFunction(2)

  nL, cF := L.NewThread()
  var _m = map[string]lua.LGFunction{
    "stop": func(L *lua.LState) int {
      cF()
      return 0
    },
  }
  module := L.SetFuncs(L.NewTable(), _m)
  L.Push(module)

  f := func() {
    nL.Push(function)
    err := nL.PCall(0, 0, nil)
    if err != nil {
      logger.Error(err.Error())
      return
    }
    t := time.Duration(interval) * time.Second
    for range time.Tick(t) {
      nL.Push(function)
      err := nL.PCall(0, 0, nil)
      if err != nil {
        logger.Error(err.Error())
        return
      }
    }
  }

  go f()

  return 1
}
