install:
	cd ./vs-uman && sudo npm install;

build:
	go build -o umanlsp main.go;