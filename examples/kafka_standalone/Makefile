all: create delete reader writer

create:
	docker build -t gcr.io/insprlabs/inspr/example/kafkasa/create:latest -f topic/create/create.Dockerfile ../..
	docker push gcr.io/insprlabs/inspr/example/kafkasa/create:latest
delete:
	docker build -t gcr.io/insprlabs/inspr/example/kafkasa/delete:latest -f topic/delete/delete.Dockerfile ../..
	docker push gcr.io/insprlabs/inspr/example/kafkasa/delete:latest
reader:
	docker build -t gcr.io/insprlabs/inspr/example/kafkasa/reader:latest -f reader/reader.Dockerfile ../..
	docker push gcr.io/insprlabs/inspr/example/kafkasa/reader:latest
writer:
	docker build -t gcr.io/insprlabs/inspr/example/kafkasa/writer:latest -f writer/writer.Dockerfile ../..
	docker push gcr.io/insprlabs/inspr/example/kafkasa/writer:latest

.PHONY: create delete reader writer