### Step 1 - build executable
FROM golang:alpine as exec-builder

COPY . /app
WORKDIR /app 

# install git (for dependencies), get dependencies & build
# do it in one RUN cmd to lessen image layers
RUN apk update && apk add --no-cache git && \
    go get -d -v && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app


### Step 2 - build a barebones image
FROM scratch

COPY --from=exec-builder /go/bin/app /go/bin/app
# for file server to work
COPY --from=exec-builder /app/pkg/static/* /pkg/static/
COPY --from=exec-builder /app/.env /.env
ENV ENVIRONMENT=DEV

ENTRYPOINT ["/go/bin/app"]