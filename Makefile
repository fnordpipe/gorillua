all: cmangos-api

cmangos-api:
	go get github.com/gorilla/mux
	go get gopkg.in/ini.v1
	go get github.com/yuin/gopher-lua
	go get github.com/go-sql-driver/mysql
	go build

clean:
	rm -rf luado
