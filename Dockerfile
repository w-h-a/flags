FROM golang AS build
WORKDIR /flags
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/flags ./

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=build /go/bin/flags /bin/flags
ENTRYPOINT [ "/bin/flags" ]