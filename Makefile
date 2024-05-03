build:
	go build -o cmd/kubectl-nic-supportpkg

install: build
	sudo cp cmd/kubectl-nic-supportpkg /usr/local/bin