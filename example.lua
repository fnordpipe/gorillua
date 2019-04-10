http = require("http")

function hello(w, r)
  print(r.getHeader("User-Agent") .. "\n")
  w.addHeader("X-Server", "luado")
  w.setStatus(404)
  w.write("hello world")
end

http.serve("127.0.0.1:5558", {
  { method = "GET", context = "/", callback = hello }
})
