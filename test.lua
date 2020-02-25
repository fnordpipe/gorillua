http = require("http")
json = require("json")

logger = require("logger")

logger.set_level(0) -- 0:debug, 1:error, 2:info

function getreq(w, r)
  --
  logger.info("http: request parse")
  logger.debug(r.get_header("User-Agent"))
  --

  --
  logger.info("json: decode")
  jsonobj, e = json.decode('{"hello":"world","foo":["a", "b", "c"]}')
  if not e then
    for k, v in pairs(jsonobj) do
      if type(v) == "table" then
        for i, j in pairs(v) do logger.debug("" .. i .. " : " .. j) end end
      if type(v) ~= "table" then logger.debug("" .. k .. " : " .. v) end
    end
  end
  --

  --
  logger.info("json: encode")
  jsonobj.tablevar = {}
  jsonobj.nilvar = nil
  jsonobj.tablevar[1] = 23
  jsonobj.tablevar[4] = 42
  logger.debug(json.encode(jsonobj))
  --

  --
  logger.info("http: response write")
  w.add_header("X-Server", "luado")
  w.set_status(200)
  w.write("hello world\n")
  --
end

http.serve("127.0.0.1:5558", {
  { method = "GET", context = "/", callback = getreq }
})
