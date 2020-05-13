BINARY=maasapi
DOCKER_TAG=${USER}-$(shell date +%H-%M-%d-%m-%Y)
HOSTADDR=$(shell ifconfig | grep -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep -Eo '([0-9]*\.){3}[0-9]*' | grep '192.*' | awk '{print $1}')

dialog:
	docker pull registry.devops.movista.ru/maas-dialog-service:release-1-6-0.fe80a520001.20191226-1735
	docker run -it --rm --name dialog -p 5000:5000 -e MAAS=http://192.168.32.109:8080 -e SQL_HOST=192.168.32.109 registry.devops.movista.ru/maas-dialog-service:release-1-6-0.fe80a520001.20191226-1735

local-pg:
	docker run -d --name maas-postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -e POSTGRES_DB=maas mdillon/postgis:10-alpine

clean-local-pg:
	docker rm -f maas-postgres

local-env-up:
	HOSTADDR=${HOSTADDR} docker-compose up -d

local-env-down:
	docker-compose down

docker:
	docker build -t registry.devops.movista.ru/maas-api:${DOCKER_TAG} .
	docker push registry.devops.movista.ru/maas-api:${DOCKER_TAG}
	docker rmi registry.devops.movista.ru/maas-api:${DOCKER_TAG}

remove-binary:
	rm ${BINARY}

build:
	GOOS=linux GOARCH=amd64 go build -o ${BINARY}

clean:
	rm ${BINARY}

giga-local-up:
	DIALOG_BASEURL=http://localhost:5000 FRONTAPI_BASEURL=http://localhost:8081 ENVIRONMENT=dev go run main.go

dev: build docker k8s remove-binary

