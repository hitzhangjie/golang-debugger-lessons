all:
	gcc -o fork -g fork.c -lpthread
	gcc -o fork_noquit -g fork_noquit.c -lpthread
	go build -gcflags 'all=-N -l' -o loop loop.go
	go build -gcflags 'all=-N -l' -o loop2 loop2.go
	go build -gcflags 'all=-N -l' -o loop3 loop3.go
	go build -gcflags 'all=-N -l' -o forloopprint forloopprint.go

clean:
	rm -rf fork fork_noquit loop loop2 loop3 forloopprint

PHONY: clean
