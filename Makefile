build:
	go build -o cmd/kubectl-nginx_supportpkg

install: build
	sudo cp cmd/kubectl-nginx_supportpkg /usr/local/bin