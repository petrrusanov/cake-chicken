FROM golang:1.11-alpine as build-backend

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOMETALINTER=2.0.11

RUN \
    apk add --no-cache --update tzdata git bash curl && \
    rm -rf /var/cache/apk/*

RUN \
    go version && \
    go get -u -v github.com/alecthomas/gometalinter && \
    cd /go/src/github.com/alecthomas/gometalinter && \
    git checkout v${GOMETALINTER} && \
    go install github.com/alecthomas/gometalinter && \
    gometalinter --install && \
    go get -u -v github.com/securego/gosec/cmd/gosec && \
    go get -u -v github.com/golang/dep/cmd/dep && \
    go get -u -v github.com/kardianos/govendor && \
    go get -u -v github.com/mattn/goveralls && \
    go get -u -v github.com/jteeuwen/go-bindata/... && \
    go get -u -v github.com/stretchr/testify && \
    go get -u -v github.com/vektra/mockery/.../

WORKDIR /go/src/github.com/boilerplate/backend

ADD app /go/src/github.com/boilerplate/backend/app
ADD vendor /go/src/github.com/boilerplate/backend/vendor
ADD .git /go/src/github.com/boilerplate/backend/.git
ADD git-rev.sh /script/git-rev.sh

RUN chmod +x /script/git-rev.sh

RUN cd app && go test ./...

RUN gometalinter --disable-all --deadline=300s --vendor --enable=vet --enable=vetshadow --enable=golint \
    --enable=staticcheck --enable=ineffassign --enable=goconst --enable=errcheck --enable=unconvert \
    --enable=deadcode --enable=gosimple --enable=gosec --exclude=test --exclude=mock --exclude=vendor ./...

RUN \
    version=$(/script/git-rev.sh) && \
    echo "version $version" && \
    go build -o app -ldflags "-X main.revision=${version} -s -w" ./app

FROM umputun/baseimage:app-latest

ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY --from=build-backend /go/src/github.com/boilerplate/backend/backend /srv/backend
RUN chown -R app:app /srv
RUN ln -s /srv/backend /usr/bin/backend

EXPOSE 3000
EXPOSE 4000

HEALTHCHECK --interval=30s --timeout=3s CMD curl --fail http://localhost:3000/ping || exit 1

CMD ["server"]
ENTRYPOINT ["/entrypoint.sh"]