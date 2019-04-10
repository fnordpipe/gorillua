http = require("http")

function hello(w) w.write("hello world") end

http.serve("127.0.0.1:5558", {
  { method = "GET", context = "/", callback = hello }
})
