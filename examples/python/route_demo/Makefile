build:
	docker build --no-cache -t gcr.io/insprlabs/inspr/example/python/client:latest -f client/client.Dockerfile .
	docker push gcr.io/insprlabs/inspr/example/python/client:latest

	docker build --no-cache -t gcr.io/insprlabs/inspr/example/python/api:latest -f api/api.Dockerfile .
	docker push gcr.io/insprlabs/inspr/example/python/api:latest