FROM alpine:latest

RUN apk add --no-cache \
    go \
    tesseract-ocr \
    tesseract-ocr-dev \
    leptonica-dev \
    g++ \
    tesseract-ocr-data-eng \
    ffmpeg

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /bin/subtle ./cmd/subtle

RUN apk del go

WORKDIR /media

ARG UID=1000
ARG GID=1000

RUN addgroup -g ${GID} docker && \
  adduser -D -u ${UID} -G docker subtle

USER subtle:docker

CMD ["/bin/subtle"]
