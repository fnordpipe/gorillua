package lualogger

import (
  "github.com/yuin/gopher-lua"

  "metagit.org/fnordpipe/gorillua/modules/logger"
)

var m = map[string]lua.LGFunction{
  "debug": debuglog,
  "error": errorlog,
  "info": infolog,
  "set_level": setlvl,
}

func Loader(L *lua.LState) int {
  module := L.SetFuncs(L.NewTable(), m)
  L.Push(module)
  return 1
}

func debuglog(L *lua.LState) int {
  msg := L.CheckString(1)
  logger.Debug(msg)
  return 0
}

func errorlog(L *lua.LState) int {
  msg := L.CheckString(1)
  logger.Error(msg)
  return 0
}

func infolog(L *lua.LState) int {
  msg := L.CheckString(1)
  logger.Info(msg)
  return 0
}

func setlvl(L *lua.LState) int {
  lvl := L.CheckNumber(1)
  if lvl < 0 || lvl > 2 {
    L.Push(lua.LNil)
    L.Push(lua.LString("loglevel not supported"))
    return 2
  }

  logger.Verbosity = int(lvl)
  L.Push(lua.LBool(true))
  return 1
}
