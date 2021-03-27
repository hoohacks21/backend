build:
	DOCKER_BUILDKIT=1 docker build -t hoohack-backend .

pull:
	docker pull alphakilo07/hoohack-backend
push:
	docker tag hoohack-backend alphakilo07/hoohack-backend
	docker push alphakilo07/hoohack-backend

cloud:
	docker push gcr.io/lassondeathoohacks/backend

run:
	docker run  --rm -d -p 8081:8081 -e PORT='8081' \
		--name hoohack-backend hoohack-backend

kill:
	docker kill hoohack-backend
	