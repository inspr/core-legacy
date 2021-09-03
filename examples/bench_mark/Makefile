build:
	docker build -t gcr.io/insprlabs/inspr/example/sender:latest -f sender/sender.Dockerfile ../..
	docker push gcr.io/insprlabs/inspr/example/sender:latest
	docker build -t gcr.io/insprlabs/inspr/example/receiver:latest -f receiver/receiver.Dockerfile ../..
	docker push gcr.io/insprlabs/inspr/example/receiver:latest