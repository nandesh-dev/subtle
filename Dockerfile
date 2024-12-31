FROM ubuntu:latest AS build-stage


# Build tools and dependencies
RUN apt-get update
RUN apt-get install -y \
  g++-10 autoconf make git golang libtool pkg-config wget xz-utils libpng-dev \
  tesseract-ocr-eng \
  protobuf-compiler

# FFmpeg static build
WORKDIR /build/ffmpeg
RUN wget https://johnvansickle.com/ffmpeg/builds/ffmpeg-git-amd64-static.tar.xz
RUN tar -xvf ffmpeg-git-amd64-static.tar.xz
RUN mkdir bin
RUN cp ffmpeg-git-*-amd64-static/ffmpeg bin/ffmpeg
RUN cp ffmpeg-git-*-amd64-static/ffprobe bin/ffprobe


# Leptonica static build ( required for tesseract )
WORKDIR /build/leptonica
RUN git clone --depth 1 https://github.com/DanBloomberg/leptonica.git .
RUN ./autogen.sh
RUN ./configure '--with-pic' '--disable-shared' '--without-zlib' '--without-jpeg' '--without-libtiff' '--without-giflib' '--without-libwebp' '--without-libwebpmux' '--without-libopenjpeg' '--disable-programs' 'CXX=g++-10' 'CFLAGS=-D DEFAULT_SEVERITY=L_SEVERITY_ERROR -g0 -O3'
RUN make 
RUN make install


# Tesseract static build
WORKDIR /build/tesseract
RUN git clone --depth 1 https://github.com/tesseract-ocr/tesseract.git .
RUN ./autogen.sh
RUN ./configure '--with-pic' '--disable-shared' '--disable-legacy' '--disable-graphics' '--disable-openmp' '--without-curl' '--without-archive' '--disable-doc' 'CXX=g++-10' 'CXXFLAGS=-DTESS_EXPORTS -g0 -O3 -ffast-math'
RUN make
RUN make install


# Subtle backend static build
WORKDIR /build/subtle
COPY . .

ENV GOPATH=$HOME/go
ENV PATH=$PATH:$GOPATH/bin

RUN sh ./scripts/entgen.sh
RUN sh ./scripts/buf/gen.sh

RUN CGO_ENABLED=1 GOOS=linux \
    go build  -a -tags netgo -ldflags '-extldflags "-static -L/usr/local/lib -ltesseract -lleptonica -lpng -lz"' ./cmd/subtle

# Subtle frontend build
WORKDIR /build/subtle/web

RUN wget -qO- https://get.pnpm.io/install.sh | bash -
RUN . /root/.bashrc && \
  pnpm env use --global lts && \
  pnpm install && \
  npx buf generate && \
  pnpm run build


# Empty volume mount points
RUN mkdir /volumes
RUN mkdir /volumes/media
RUN mkdir /volumes/config



FROM alpine:latest AS user-stage

# Setup user and group
ENV UID=1000
ENV GID=1000

RUN addgroup -g $GID docker
RUN adduser -S -u $UID -G docker subtle



FROM scratch

# User and group
COPY --from=user-stage /etc/passwd /etc/passwd
COPY --from=user-stage /etc/group /etc/group

USER subtle:docker


# Binaries
COPY --from=build-stage /build/subtle/subtle /subtle
COPY --from=build-stage /build/subtle/web/dist /public

COPY --from=build-stage /build/ffmpeg/bin/ffmpeg /usr/local/bin/ffmpeg
COPY --from=build-stage /build/ffmpeg/bin/ffprobe /usr/local/bin/ffprobe


# OCR Language data
COPY --from=build-stage /usr/share/tesseract-ocr/5/tessdata/eng.traineddata /usr/local/share/tessdata/eng.traineddata


# Volume mounts
COPY --from=build-stage /volumes/media /media
COPY --from=build-stage /volumes/config /config

CMD ["/subtle"]
