package main

import (
  "fmt"
  "os"

  "metagit.org/fnordpipe/gorillua/modules/bindings/luabase64"
  "metagit.org/fnordpipe/gorillua/modules/bindings/luahttp"
  "metagit.org/fnordpipe/gorillua/modules/bindings/luajson"
  "metagit.org/fnordpipe/gorillua/modules/bindings/luarequest"
  "metagit.org/fnordpipe/gorillua/modules/logger"
  "github.com/yuin/gopher-lua"
)

var L *lua.LState

var _LUA_PATH string

func main() {
  if len(os.Args) < 2 {
    logger.Stdout(fmt.Sprintf("USAGE: %s <lua> [...]", os.Args[0]))
    os.Exit(1)
  }

  lua.LuaPathDefault = _LUA_PATH

  L := lua.NewState()
  defer L.Close()

  t := L.CreateTable(0, len(os.Args))
  for _, v := range os.Args {
    t.Append(lua.LString(v))
  }
  L.SetGlobal("arg", t)

  L.PreloadModule("base64", luabase64.Loader)
  L.PreloadModule("http", luahttp.Loader)
  L.PreloadModule("json", luajson.Loader)
  L.PreloadModule("request", luarequest.Loader)

  if err := L.DoFile(os.Args[1]); err != nil {
    logger.Error("Cannot parse lua script")
    logger.Debug(err.Error())
    os.Exit(2)
  }
}
