
test/gapi/server:
	GO111MODULE=on go build -o .bin/gapi ./gapi/gapi
	.bin/gapi