all:
	gcc -o fork -g fork.c -lpthread
	gcc -o fork_noquit -g fork_noquit.c -lpthread
	GOOS=linux GOARCH=amd64 go build -gcflags 'all=-N -l' -o loop loop.go
	GOOS=linux GOARCH=amd64 go build -gcflags 'all=-N -l' -o loop2 loop2.go
	GOOS=linux GOARCH=amd64 go build -gcflags 'all=-N -l' -o loop3 loop3.go
	GOOS=linux GOARCH=amd64 go build -gcflags 'all=-N -l' -o forloopprint forloopprint.go

clean:
	rm -rf fork fork_noquit loop loop2 loop3 forloopprint *.dSYM

PHONY: clean
