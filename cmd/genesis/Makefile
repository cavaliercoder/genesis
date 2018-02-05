SOURCES = $(wildcard *.go)

all: genesis

genesis: $(SOURCES)
	go build -x -o $@ $^

clean:
	rm -f genesis
