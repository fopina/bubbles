BINARY?=bubbles

all: deps bubbles

deps:
	go get -d

bubbles:
	go build -tags static -ldflags "-s -w" -o ${BINARY}

slim:
	go build -ldflags "-s -w" -o ${BINARY}	

clean:
	rm ${BINARY}
