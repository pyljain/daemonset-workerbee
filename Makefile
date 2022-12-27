all: clean build run

clean:
	rm -r ./root

build:
	go build

run:
	./workerbee --loc ./root workerbee-bucket1 workerbee-bucket2 

docker:
	docker build --platform linux/amd64  -t payaljain/workerbee:1 . && \
	docker push payaljain/workerbee:1

kubernetes:
	kubectl apply -f manifests/