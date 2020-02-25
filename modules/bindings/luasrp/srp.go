package luasrp

import (
  "github.com/yuin/gopher-lua"
  "github.com/blizzlike-org/wowpasswd/srp"
)

var m = map[string]lua.LGFunction{
  "hash": _hash,
  "new": _new,
}

func Loader(L *lua.LState) int {
  module := L.SetFuncs(L.NewTable(), m)
  L.Push(module)
  return 1
}

func _hash(L *lua.LState) int {
  username := L.CheckString(1)
  password := L.CheckString(2)
  hash := srp.Hash(username, password)
  L.Push(lua.LString(hash))
  return 1
}

func _new(L *lua.LState) int {
  _srp := srp.New()

  var _m = map[string]lua.LGFunction{
    "compute_verifier": func(L *lua.LState) int {
      identifier := L.CheckString(1)
      _srp.ComputeVerifier(identifier)
      return 0
    },
    "generate_salt": func(L *lua.LState) int {
      _srp.GenerateSalt()
      return 0
    },
    "get_salt": func(L *lua.LState) int {
      L.Push(lua.LString(_srp.GetSalt()))
      return 1
    },
    "get_verifier": func(L *lua.LState) int {
      L.Push(lua.LString(_srp.GetVerifier()))
      return 1
    },
    "proof_verifier": func(L *lua.LState) int {
      v := L.CheckString(1)

      proof := _srp.ProofVerifier(v)
      L.Push(lua.LBool(proof))
      return 1
    },
    "set_salt": func(L *lua.LState) int {
      salt := L.CheckString(1)
      _srp.SetSalt(salt)
      return 0
    },
    "set_verifier": func(L *lua.LState) int {
      verifier := L.CheckString(1)
      _srp.SetVerifier(verifier)
      return 0
    },
  }

  _module := L.SetFuncs(L.NewTable(), _m)
  L.Push(_module)
  return 1
}
