FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /inventory .

FROM gcr.io/distroless/static-debian12
COPY --from=build /inventory /inventory
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/inventory"]
