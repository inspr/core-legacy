build:
	docker build -t gcr.io/insprlabs/inspr/example/kafkasa/create:latest -f topic/create/create.Dockerfile ../..
	docker push gcr.io/insprlabs/inspr/example/kafkasa/create:latest

	docker build -t gcr.io/insprlabs/inspr/example/kafkasa/delete:latest -f topic/delete/delete.Dockerfile ../..
	docker push gcr.io/insprlabs/inspr/example/kafkasa/delete:latest
	
	docker build -t gcr.io/insprlabs/inspr/example/kafkasa/writer:latest -f writer/writer.Dockerfile ../..
	docker push gcr.io/insprlabs/inspr/example/kafkasa/writer:latest
