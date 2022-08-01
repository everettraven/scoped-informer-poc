FROM golang:1.18 as builder

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build -o /scoped-informer-demo

FROM gcr.io/distroless/base-debian10

COPY --from=builder /scoped-informer-demo /scoped-informer-demo

ENTRYPOINT ["/scoped-informer-demo"]