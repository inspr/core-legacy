build:
	docker build -t gcr.io/red-inspr/inspr/example/chtest:latest -f chtest/Dockerfile ../..
	docker push gcr.io/red-inspr/inspr/example/chtest:latest
	docker build -t gcr.io/red-inspr/inspr/example/chcheck:latest -f chcheck/Dockerfile ../..
	docker push gcr.io/red-inspr/inspr/example/chcheck:latest