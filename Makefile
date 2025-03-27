build:
	go build -o cmd/kubectl-nginx_supportpkg

debugger:
	docker buildx build --build-context project=nginx-debugger --platform linux/amd64 -t nginx-debugger -f nginx-debugger/Dockerfile .
#	docker tag nginx-debugger:latest mrajagopal/f5-utils:latest
#	docker push mrajagopal/f5-utils:latest

install: build
	sudo cp cmd/kubectl-nginx_supportpkg /usr/local/bin