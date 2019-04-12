http = require("http")
json = require("json")

function hello(w, r)
  foo, e = json.decode('{"hello":"world","foo":["a", "b", "c"]}')
  if not e then
    for k, v in pairs(foo) do
      if type(v) == "table" then 
        for i, j in pairs(v) do print(i .. " fgh " .. j) end end
      if type(v) ~= "table" then print(k .. " asd " .. v) end
    end
  end

  print("lua encode")
  foo.lol = {}
  foo.rofl = {}
  foo.nilnil = nil
  foo.rofl[1] = 23
  foo.rofl[4] = 42
  print(json.encode(foo))
  print(r.getHeader("User-Agent") .. "\n")
  w.addHeader("X-Server", "luado")
  w.setStatus(404)
  w.write("hello world")
end

http.serve("127.0.0.1:5558", {
  { method = "GET", context = "/", callback = hello }
})
