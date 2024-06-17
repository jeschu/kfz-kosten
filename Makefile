default: show-config

show-config:
	@echo ">> show-config <<"
	@echo "=== ~/Library/Application\ Support/kfzs.yml ==="
	@cat ~/Library/Application\ Support/kfzs.yml
	@echo "--- ~/Library/Application\ Support/kfzs.yml ---"

clean:
	@echo ">> clean <<"
	@rm -rfv version.go build

generate:
	@echo ">> generate <<"
	@go generate

build: generate
	@echo ">> build <<"
	@go build -o "build/kfz" ./...

install: clean build
	@echo ">> install <<"
	@go install ./...