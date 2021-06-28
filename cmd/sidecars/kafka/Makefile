build:
	docker build -t gcr.io/insprlabs/inspr/sidecar/kafka:build -f build.Dockerfile ../../..
    
	docker push gcr.io/insprlabs/inspr/sidecar/kafka:build
    
	docker build -t gcr.io/insprlabs/inspr/sidecar/kafka:deploy -f deploy.Dockerfile ../../..
    
	docker push gcr.io/insprlabs/inspr/sidecar/kafka:deploy
    
	docker build -t gcr.io/insprlabs/inspr/sidecar/kafka:latest -f Dockerfile ../../..
    
	docker push gcr.io/insprlabs/inspr/sidecar/kafka:latest