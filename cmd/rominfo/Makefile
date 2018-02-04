SOURCES = $(wildcard *.go)

all: rominfo

rominfo: $(SOURCES)
	go build -x -o $@ $^

clean:
	rm -f rominfo
