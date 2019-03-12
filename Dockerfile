FROM golang:1.12-stretch

RUN apt update && \
	apt install -y --no-install-recommends libsdl2-gfx-dev && \
	rm -rf /var/lib/apt/lists/*

RUN go get -v github.com/veandco/go-sdl2/sdl && \
	rm -fr /go/src/github.com/veandco/go-sdl2/.go-sdl2-examples && \
	rm -fr /go/src/github.com/veandco/go-sdl2/.go-sdl2-libs && \
	rm -fr /go/src/github.com/veandco/go-sdl2/.git
