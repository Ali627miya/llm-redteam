# Build: docker build -t llm-redteam ./redteam
# Run CLI: docker run --rm llm-redteam redteam version
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /redteam ./cmd/redteam \
  && CGO_ENABLED=0 go build -ldflags="-s -w" -o /mocktarget ./cmd/mocktarget

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /redteam /usr/local/bin/redteam
COPY --from=build /mocktarget /usr/local/bin/mocktarget
COPY examples /app/examples
ENTRYPOINT ["redteam"]
CMD ["version"]
