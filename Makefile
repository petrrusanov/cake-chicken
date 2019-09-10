OS=linux
ARCH=amd64

lint:
	gometalinter --disable-all --deadline=300s --vendor --enable=vet --enable=vetshadow --enable=golint \
    --enable=staticcheck --enable=ineffassign --enable=goconst --enable=errcheck --enable=unconvert \
    --enable=deadcode --enable=gosimple --enable=gosec --exclude=test --exclude=mock --exclude=vendor ./...
