BINARY?=bubbles

all: deps bubbles

generate:
	go get github.com/shurcooL/vfsgen
	go generate

deps: generate
	go get -d

bubbles:
	go build -tags static -ldflags "-s -w" -o ${BINARY}

slim:
	go build -ldflags "-s -w" -o ${BINARY}	

clean:
	rm ${BINARY}
