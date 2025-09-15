.PHONY: test test-html clean

COVER := coverage.out

test:
	@go test ./... -coverprofile=$(COVER)
	@go tool cover -func=$(COVER)

test-html:
	@go test ./... -coverprofile=$(COVER)
	@go tool cover -html=$(COVER)

cleaning-test:
	@rm -f $(COVER)

generate:
	@go generate ./...
