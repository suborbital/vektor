
vk/tester:
	GO111MODULE=on go build ./vk/test

vk/tester/run: vk/tester
	./test

test:
	go test -v -count=1 ./...

deps:
	go get -u -d ./...

.PHONY: vk/tester vk/tester/run test deps