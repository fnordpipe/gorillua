> **go**ril**lua**

# doc

    make LUA_PATH="./?.lua;/usr/share/gorillua/?.lua"
    ./gorillua ./example.lua [...]

## bindings

### base64

    local base64 = require("base64")

    b = base64.encode("hello world")
    s = base64.decode(b)

### cron

    local cron = require("cron")

    -- 1st argument is the interval in seconds while
    -- 2nd is a function that will be called
    c = cron.run(10, function() print("hello world") end)
    c.stop()

### http

start serving requests

    local http = require("http")

    -- the 3rd argument is either nil or a path to a directory containing static files
    http.serve("127.0.0.1:5558", {
      { method = "GET", context = "/", callback = func(w, r) w.write("hello world") end }
      { method = "POST", context = "/", callback = func(w, r) w.write("hello post world") end }
    }, nil)

define callback functions for http requests

    function example(w, r)
      -- r is the request
      useragent = r.get_header("User-Agent")

      cookie = r.get_cookie("key")

      params = r.parse_form()
      print(params.username)

      -- w is the response writer
      w.add_header("X-Header-Foo", "example")
      w.set_status(200)

      w.set_cookie("key", "value", "/", 86400, true)

      w.write("this is the response body")
    end

### json

    local json = require("json")

    foo = json.decode('{"hello":"world"}')
    print(foo.hello)
    foo.bla = 3
    json.encode(foo)

### mariadb

    local mariadb = require("mariadb")

    db = mariadb.open("user", "password", "127.0.0.1:3306", "database")
    result = db.query("SELECT * FROM table WHERE id = ?", 1)
    for k, v in pairs(result) do print(v.id) end
    db.close()

### request

    local request = require("request")

    -- 3rd argument contains the request body
    -- 4th argument is a table of request headers
    code, body, header, err = request.send(
      "GET", "http://example.org", nil, nil)

### socket

    local socket = require("socket")

    -- 3rd argument is optional and defines a timeout in seconds
    c = socket.open("tcp", "127.0.0.1:8080", nil)
    c.close()
