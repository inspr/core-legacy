build:
	docker build -t gcr.io/insprlabs/inspr/example/controller:latest -f Dockerfile ../..
	docker push gcr.io/insprlabs/inspr/example/controller:latest

	insprctl apply -f yamls/dapp.yaml
