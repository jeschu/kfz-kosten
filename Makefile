default: show-config

show-config:
	@echo ">> show-config <<"
	@echo "=== ~/Library/Application\ Support/kfzs.yaml ==="
	@cat $(HOME)/Library/Application\ Support/kfzs.yaml
	@echo "--- ~/Library/Application\ Support/kfzs.yaml ---"

edit-config:
	@code $(HOME)/Library/Application\ Support/kfzs.yaml

clean:
	@echo ">> clean <<"
	@rm -rfv version.go build

generate:
	@echo ">> generate <<"
	@go generate

build: generate
	@echo ">> build <<"
	@mkdir -p build
	@go build -o "build" ./...

install: clean build
	@echo ">> install <<"
	@go install ./...