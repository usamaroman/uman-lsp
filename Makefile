install:
	cd ./vs-uman && sudo npm install;

build:
	go build -o ./bin/umanlsp ./cmd/umanlsp/main.go;