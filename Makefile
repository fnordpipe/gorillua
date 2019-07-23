LUA_PATH?=./?.lua;/usr/share/luagoesweb/?.lua

all: gorillua

gorillua:
	go get github.com/gorilla/mux
	go get github.com/google/uuid
	go get github.com/yuin/gopher-lua
	go get github.com/go-sql-driver/mysql
	go get metagit.org/blizzlike/wowpasswd
	go build -ldflags "-X main._LUA_PATH=${LUA_PATH}"

clean:
	rm -rf gorillua
