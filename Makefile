build:
	@go build -o packer-builder-libvirt main.go
test:
	@go test ./...
