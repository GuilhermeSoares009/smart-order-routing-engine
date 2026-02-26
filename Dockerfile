FROM golang:1.22 AS build
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/sor ./cmd/sor

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /out/sor /sor
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/sor"]
