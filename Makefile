build:
	go build -o cmd/kubectl-nginx_supportpkg

debugger:
	docker buildx build --platform linux/amd64 -t nginx-debugger -f nginx-debugger/Dockerfile .

install: build
	sudo cp cmd/kubectl-nginx_supportpkg /usr/local/bin