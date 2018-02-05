all: check rominfo genesis

check:
	go test -v ./...

rominfo:
	cd cmd/rominfo && make

genesis:
	cd cmd/genesis && make

clean:
	cd cmd/rominfo && make clean
	cd cmd/genesis && make clean
