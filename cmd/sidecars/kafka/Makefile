build:
	docker build -t gcr.io/red-inspr/inspr/sidecar/kafka:build -f build.Dockerfile ../../..
    
	docker push gcr.io/red-inspr/inspr/sidecar/kafka:build
    
	docker build -t gcr.io/red-inspr/inspr/sidecar/kafka:deploy -f deploy.Dockerfile ../../..
    
	docker push gcr.io/red-inspr/inspr/sidecar/kafka:deploy
    
	docker build -t gcr.io/red-inspr/inspr/sidecar/kafka:latest -f Dockerfile ../../..
    
	docker push gcr.io/red-inspr/inspr/sidecar/kafka:latest