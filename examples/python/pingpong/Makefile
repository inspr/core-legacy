build:
	docker build --no-cache -t gcr.io/insprlabs/inspr/example/python/ping:latest -f ping/ping.Dockerfile .
	docker push gcr.io/insprlabs/inspr/example/python/ping:latest

	docker build --no-cache -t gcr.io/insprlabs/inspr/example/python/pong:latest -f pong/pong.Dockerfile .
	docker push gcr.io/insprlabs/inspr/example/python/pong:latest