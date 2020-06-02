FROM golang:alpine AS builder
WORKDIR /src
ADD ./go.mod ./go.sum ./
RUN go mod download
ADD . ./
RUN CGO_ENABLED=0 go build -o pluto ./cmd/*.go
RUN ls -la


FROM alpine 
WORKDIR /root/
COPY --from=builder /src/pluto .
RUN ls -la
EXPOSE 80
ENTRYPOINT ["./pluto"]
