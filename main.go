package main

import (
  "fmt"
  "os"

  "metagit.org/fnordpipe/luado/modules/bindings/luahttp"
  "metagit.org/fnordpipe/luado/modules/logger"
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
  if err := L.DoFile(os.Args[1]); err != nil {
    logger.Error("Cannot parse lua script")
    logger.Debug(err.Error())
    os.Exit(2)
  }
}
