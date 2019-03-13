FROM golang:1.12-stretch

RUN apt update && \
	apt install -y --no-install-recommends libsdl2-gfx-dev && \
	rm -rf /var/lib/apt/lists/*
