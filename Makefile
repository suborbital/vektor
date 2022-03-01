
vk/tester:
	go build ./vk/test/cmd

vk/tester/run: vk/test/cmd
	./cmd

test:
	go test -v -count=1 ./...

deps:
	go get -u -d ./...

mocks:
	mockery --name=RouterWrapperTester --dir=./vk/test --output=./vk/test/mocks

.PHONY: vk/tester vk/tester/run test deps
