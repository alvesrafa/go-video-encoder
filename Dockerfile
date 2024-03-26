FROM golang:1.22-alpine3.19

ENV PATH="$PATH:/bin/bash" \
    BENTO4_BIN="/opt/bento4/bin" \
    PATH="$PATH:/opt/bento4/bin" \
    CGO_CFLAGS="-D_LARGEFILE64_SOURCE"

# FFMPEG
RUN apk add --update ffmpeg bash make

# Install Bento
WORKDIR /tmp/bento4

RUN apk add --update --upgrade python3 git cmake unzip bash gcc g++ scons
RUN git clone https://github.com/axiomatic-systems/Bento4.git \
&& cd Bento4 && mkdir cmakebuild && cd cmakebuild && cmake -DCMAKE_BUILD_TYPE=Release ..

WORKDIR /go/src

ENTRYPOINT ["tail", "-f", "/dev/null"]