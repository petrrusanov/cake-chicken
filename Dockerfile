FROM golang:1.13-alpine as build-backend

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

WORKDIR /go/src/github.com/petrrusanov/cake-chicken

ADD app /go/src/github.com/petrrusanov/cake-chicken/app
ADD vendor /go/src/github.com/petrrusanov/cake-chicken/vendor
ADD .git /go/src/github.com/petrrusanov/cake-chicken/.git
ADD git-rev.sh /script/git-rev.sh

RUN chmod +x /script/git-rev.sh

RUN cd app && go test ./...

RUN gometalinter --disable-all --deadline=300s --vendor --enable=vet --enable=vetshadow --enable=golint \
    --enable=staticcheck --enable=ineffassign  --enable=errcheck --enable=unconvert \
    --enable=deadcode --enable=gosimple --exclude=test --exclude=mock --exclude=vendor ./...

RUN \
    version=$(/script/git-rev.sh) && \
    echo "version $version" && \
    go build -o backend -ldflags "-X main.revision=${version} -s -w" ./app

FROM umputun/baseimage:app-latest

ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY --from=build-backend /go/src/github.com/petrrusanov/cake-chicken/backend /srv/backend
RUN chown -R app:app /srv
RUN ln -s /srv/backend /usr/bin/backend

RUN mkdir -p /data/db

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=3s CMD curl --fail http://localhost:3000/ping || exit 1

CMD ["server"]
ENTRYPOINT ["/entrypoint.sh"]