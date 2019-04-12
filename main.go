package main

import (
  "fmt"
  "os"

  "metagit.org/fnordpipe/luagoesweb/modules/bindings/luahttp"
  "metagit.org/fnordpipe/luagoesweb/modules/bindings/luajson"
  "metagit.org/fnordpipe/luagoesweb/modules/logger"
  "github.com/yuin/gopher-lua"
)

var L *lua.LState

func main() {
  if len(os.Args) != 2 {
    logger.Stdout(fmt.Sprintf("USAGE: %s <lua>", os.Args[0]))
    os.Exit(1)
  }

  L := lua.NewState()
  defer L.Close()

  L.PreloadModule("http", luahttp.Loader)
  L.PreloadModule("json", luajson.Loader)
  if err := L.DoFile(os.Args[1]); err != nil {
    logger.Error("Cannot parse lua script")
    logger.Debug(err.Error())
    os.Exit(2)
  }
}
