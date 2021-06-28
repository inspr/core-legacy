build:
	docker build -t gcr.io/insprlabs/inspr/example/ping:latest -f ping/ping.Dockerfile ../..
	docker push gcr.io/insprlabs/inspr/example/ping:latest
	docker build -t gcr.io/insprlabs/inspr/example/pong:latest -f pong/pong.Dockerfile ../..
	docker push gcr.io/insprlabs/inspr/example/pong:latest
