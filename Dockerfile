FROM oven/bun:1.0.21-debian as bun

WORKDIR /app

COPY tailwind.config.js style.css ./

COPY ./internals/templates/*.templ ./internals/templates/

RUN bunx tailwindcss@latest -i ./style.css -o ./dist/style.css



FROM golang:1.21.5-bullseye as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

RUN useradd -u 1001 crocoder

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@e98db353f87ebedea804cb3dc3200a826afb8904

RUN templ generate

RUN go build \
  -ldflags="-linkmode external -extldflags -static" \
  -o web \
  ./cmd/web/main.go



FROM scratch

WORKDIR /

COPY --from=bun /app/dist/style.css /dist/style.css

COPY --from=builder /app/web /web

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /etc/passwd /etc/passwd

USER crocoder

EXPOSE 3000

CMD ["/web"]
