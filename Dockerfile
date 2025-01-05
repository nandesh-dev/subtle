FROM ubuntu:latest AS backend-build-stage


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


# Subtle static build
WORKDIR /build/subtle
COPY . .

ENV GOPATH=$HOME/go
ENV PATH=$PATH:$GOPATH/bin

RUN sh ./scripts/entgen.sh
RUN sh ./scripts/buf/gen.sh

RUN CGO_ENABLED=1 GOOS=linux \
    go build  -a -tags netgo -ldflags '-extldflags "-static -L/usr/local/lib -ltesseract -lleptonica -lpng -lz"' ./cmd/subtle


# Empty volume mount points
RUN mkdir /volumes
RUN mkdir /volumes/media
RUN mkdir /volumes/config



FROM node:22-alpine AS frontend-build-stage


RUN corepack enable pnpm

WORKDIR /build/subtle
COPY ./web .

# Subtle build
RUN pnpm i
RUN npx buf generate
RUN pnpm run build



FROM scratch


# Subtle Backend binary
COPY --from=backend-build-stage /build/subtle/subtle /subtle

# Subtle Frontend files
COPY --from=frontend-build-stage /build/subtle/dist /public

# FFMpeg binaries
COPY --from=backend-build-stage /build/ffmpeg/bin/ffmpeg /usr/local/bin/ffmpeg
COPY --from=backend-build-stage /build/ffmpeg/bin/ffprobe /usr/local/bin/ffprobe

# OCR Language data
COPY --from=backend-build-stage /usr/share/tesseract-ocr/5/tessdata/eng.traineddata /usr/local/share/tessdata/eng.traineddata


# Volume mounts
COPY --from=backend-build-stage /volumes/media /media
COPY --from=backend-build-stage /volumes/config /config

CMD ["/subtle"]
