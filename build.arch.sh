#!/bin/bash

cd $(dirname $0)

mkdir -p data

if [ ! -e data/arch.tar.gz ]; then
	curl -L http://os.archlinuxarm.org/os/ArchLinuxARM-rpi-2-latest.tar.gz -o data/arch.tar.gz
fi

if ! docker images fopina/misc:bubbles-arm-builder-arch-base | grep -q base; then
	gunzip -c data/arch.tar.gz | docker import - fopina/misc:bubbles-arm-builder-arch-base
fi

docker build -t fopina/misc:bubbles-arm-builder-arch -f Dockerfile.arch .
