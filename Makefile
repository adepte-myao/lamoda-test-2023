up:
	docker-compose --file ./deploy/docker-compose/docker-compose.yml --project-name lamoda-test-some-service up --build --detach

down:
	docker-compose --file ./deploy/docker-compose/docker-compose.yml --project-name lamoda-test-some-service down --volumes

test:
	go test ./...