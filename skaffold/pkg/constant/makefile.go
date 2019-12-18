package constant

const Makefile = `
ver=$(shell git tag|tail -n1)
.PHONY: build
build:
	docker build --no-cache -t {{ .CompanyName }}/{{ .SrvName }}:$(ver) .
	docker image prune -f --filter label=stage=builder
	docker push {{ .CompanyName }}/{{ .SrvName }}:$(ver)
	docker tag {{ .CompanyName }}/{{ .SrvName }}:$(ver) {{ .CompanyName }}/{{ .SrvName }}:latest
	docker push {{ .CompanyName }}/{{ .SrvName }}:latest
`
