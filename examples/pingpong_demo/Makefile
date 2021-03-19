build:
	docker build -t gcr.io/red-inspr/inspr/example/ping:latest -f ping/Dockerfile ../..
	docker push gcr.io/red-inspr/inspr/example/ping:latest
	docker build -t gcr.io/red-inspr/inspr/example/pong:latest -f pong/Dockerfile ../..
	docker push gcr.io/red-inspr/inspr/example/pong:latest
