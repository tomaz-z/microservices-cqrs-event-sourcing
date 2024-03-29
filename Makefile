.PHONY: logs run stop run/% stop/% restart/% swagger/validate/% swagger/remove/% swagger/generate/%

logs:
	docker-compose \
		-p eventsourcing \
		logs --follow

run:
	docker-compose \
		-p eventsourcing \
		up --detach

stop:
	docker-compose \
		-p eventsourcing \
		down --remove-orphans \
		-v

run/%:
	docker-compose \
		-p eventsourcing \
		up --detach \
		$*

stop/%:
	docker-compose \
		-p eventsourcing \
		stop \
		$*

restart/%:
	docker-compose \
		-p eventsourcing \
		restart \
		$*

swagger/validate/%:
	docker run \
		--rm \
		-t \
		--env GOPATH=/go \
		-v $(shell pwd):/go/src \
		-w /go/src \
		quay.io/goswagger/swagger \
		validate \
		./services/$*/api/specifications.yml

swagger/remove/%:
	rm -rf ./services/$*/api/models
	rm -rf ./services/$*/api/restapi

swagger/generate/%: swagger/validate/% swagger/remove/%
	docker run \
		--rm \
		-t \
		--env GOPATH=/go \
		-v $(shell pwd):/go/src \
		-w /go/src \
		quay.io/goswagger/swagger \
		generate server \
		--spec=./services/$*/api/specifications.yml \
		--target=./services/$*/api \
		--exclude-main
		
swagger/client/generate/%: swagger/validate/%
	rm -rf ./internal/api_clients/$*
	mkdir ./internal/api_clients/$*

	docker run \
		--rm \
		-t \
		--env GOPATH=/go \
		-v $(shell pwd):/go/src \
		-w /go/src \
		quay.io/goswagger/swagger \
		generate client \
		--spec=./services/$*/api/specifications.yml \
		--target=./internal/api_clients/$*
		--name=$*
		