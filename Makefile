
build/vk/test:
	GO111MODULE=on go build ./vk/test

run/vk/test: build/vk/test
	./test

test:
	go test -v -count=1 ./...

deps:
	go get -u -d ./...