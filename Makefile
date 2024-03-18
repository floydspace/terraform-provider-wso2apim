default: install

generate:
	go generate ./...
	npx cdktf-registry-docs convert \
		--files='**/*' \
		--languages='typescript,python' \
		--provider-from-registry="floydspace/wso2apim"

install:
	go install .

test:
	go test -count=1 -parallel=4 ./...

testacc:
	TF_ACC=1 go test -count=1 -parallel=4 -timeout 10m -v ./...
