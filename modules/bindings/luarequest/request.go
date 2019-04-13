package luarequest

import (
  "bytes"
  "io/ioutil"
  "net/http"
  "strings"

  "github.com/yuin/gopher-lua"
)

type RouterInfo struct {
  Method string
  Context string
  Callback lua.LValue
}

var m = map[string]lua.LGFunction{
  "send": send,
}

func Loader(L *lua.LState) int {
  module := L.SetFuncs(L.NewTable(), m)
  L.Push(module)
  return 1
}

func send(L *lua.LState) int {
  method := L.CheckString(1)
  url := L.CheckString(2)
  body := L.CheckAny(3)
  header := L.CheckAny(4)

  var req *http.Request
  var err error
  switch bt := body.(type) {
    case *lua.LNilType:
      req, err = http.NewRequest(method, url, nil)
    case lua.LString:
      req, err = http.NewRequest(method, url,
        bytes.NewBuffer([]byte(bt.String())))
  }

  if err != nil {
    L.Push(lua.LNil)
    L.Push(lua.LNil)
    L.Push(lua.LNil)
    L.Push(lua.LString(err.Error()))
    return 4
  }

  switch ht := header.(type) {
    case *lua.LTable:
      ht.ForEach(func(k, v lua.LValue) {
        req.Header.Set(k.String(), v.String())
      })
  }

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    L.Push(lua.LNil)
    L.Push(lua.LNil)
    L.Push(lua.LNil)
    L.Push(lua.LString(err.Error()))
    return 4
  }
  defer resp.Body.Close()

  respbody, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    L.Push(lua.LNil)
    L.Push(lua.LNil)
    L.Push(lua.LNil)
    L.Push(lua.LString(err.Error()))
    return 4
  }

  respheader := L.CreateTable(0, len(resp.Header))
  for k, v := range resp.Header {
    respheader.RawSetH(lua.LString(k), lua.LString(strings.Join(v, "")))
  }

  L.Push(lua.LNumber(resp.StatusCode))
  L.Push(lua.LString(respbody))
  L.Push(respheader)

  return 3
}
