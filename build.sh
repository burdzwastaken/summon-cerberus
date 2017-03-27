#!/usr/bin/env sh

docker run -it --rm -v $PWD:/usr/src/summon-cerberus -w /usr/src/summon-cerberus golang bash -c '
goos="darwin linux windows"
arch="amd64"
export VERSION=$(grep -Po -m1 "(?<=## \[)[0-9.]+" CHANGELOG.md)

rm -r ./pkg/dist
mkdir ./pkg/dist

go get
for GOOS in $goos; do
  export GOOS
  for GOARCH in $arch; do
    export GOARCH
    go build -v -o pkg/${GOOS}_${GOARCH}/summon-cerberus
    pushd pkg/${GOOS}_${GOARCH}/
    tar -czvf ../dist/summon-cerberus_v${VERSION}_${GOOS}_${GOARCH}.tar.gz ./summon-cerberus
    popd
  done
done
'
