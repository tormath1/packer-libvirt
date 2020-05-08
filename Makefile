build:
	@go build -o packer-builder-libvirt main.go
test:
	@go test ./...

install: build
	@mv packer-builder-libvirt ${HOME}/.packer.d/plugins

validate: install
	@packer validate testdata/packer.json
