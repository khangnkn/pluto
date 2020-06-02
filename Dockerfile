FROM golang:alpine AS builder
WORKDIR /src
ADD ./go.mod ./go.sum ./
RUN go mod download
ADD . ./
RUN CGO_ENABLED=0 go build -tags static_all -o pluto ./cmd/*.go
RUN ls -la


FROM scratch
WORKDIR /root/
COPY --from=builder /src/config/pluto.yaml .
COPY --from=builder /src/pluto .
EXPOSE 8083
ENTRYPOINT ["./pluto"]
