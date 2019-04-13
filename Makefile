LUA_PATH?=./?.lua;/usr/share/luagoesweb/?.lua

all: luagoesweb

luagoesweb:
	go get github.com/gorilla/mux
	go get gopkg.in/ini.v1
	go get github.com/yuin/gopher-lua
	go get github.com/go-sql-driver/mysql
	go build -ldflags "-X main._LUA_PATH=${LUA_PATH}"

clean:
	rm -rf luagoesweb
