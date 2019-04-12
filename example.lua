http = require("http")
json = require("json")

function hello(w, r)
  foo, e = json.decode('{"hello":"world","foo":["a", "b", "c"]}')
  if not e then
    for k, v in pairs(foo) do
      if type(v) == "table" then print(k) end
      if type(v) ~= "table" then print(k .. " " .. v) end
    end
  end
  print(r.getHeader("User-Agent") .. "\n")
  w.addHeader("X-Server", "luado")
  w.setStatus(404)
  w.write("hello world")
end

http.serve("127.0.0.1:5558", {
  { method = "GET", context = "/", callback = hello }
})
