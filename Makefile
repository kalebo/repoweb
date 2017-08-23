UPX=$(shell find . -iname upx*)

.PHONY: build compress compress_linux compress_win clean

compress: compress_linux compress_win

compress_linux: repoweb
	$(UPX) --brute repoweb

compress_win: repoweb.exe
	$(UPX) --brute repoweb.exe

build: repoweb repoweb.exe

repoweb: main.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"

repoweb.exe: main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w"

clean:
	rm repoweb.exe repoweb