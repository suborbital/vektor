
build/gapi/test:
	GO111MODULE=on go build ./gapi/test

run/gapi/test: build/gapi/test
	./test