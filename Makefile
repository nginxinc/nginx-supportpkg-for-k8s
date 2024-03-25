build:
	go build -o cmd/kubectl-kic-supportpkg

install: build
	sudo cp cmd/kubectl-kic-supportpkg /usr/local/bin