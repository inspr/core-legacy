build:
	docker build -t gcr.io/insprlabs/inspr/example/chtest:latest -f chtest/Dockerfile ../..
	docker push gcr.io/insprlabs/inspr/example/chtest:latest
	docker build -t gcr.io/insprlabs/inspr/example/chcheck:latest -f chcheck/Dockerfile ../..
	docker push gcr.io/insprlabs/inspr/example/chcheck:latest