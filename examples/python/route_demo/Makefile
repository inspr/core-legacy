build:
	docker build -t gcr.io/insprlabs/inspr/example/python/route/client:latest -f client/client.Dockerfile ../..
	docker push gcr.io/insprlabs/inspr/example/python/route/client:latest
	docker build -t gcr.io/insprlabs/inspr/example/python/route/api:latest -f api/api.Dockerfile ../..
	docker push gcr.io/insprlabs/inspr/example/python/route/api:latest