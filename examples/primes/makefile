
PROJ_SRC = ../../.

# docker images

build: ## Build the container
	# generator
	docker build -f ./generator/Dockerfile -t gcr.io/insprlabs/inspr/examples/primes/generator:latest $(PROJ_SRC)
	docker push gcr.io/insprlabs/inspr/examples/primes/generator:latest 

	# filter
	docker build -f  ./filter/Dockerfile -t gcr.io/insprlabs/inspr/examples/primes/filter:latest $(PROJ_SRC)
	docker push gcr.io/insprlabs/inspr/examples/primes/filter:latest 

	# printer
	docker build -f ./printer/Dockerfile -t gcr.io/insprlabs/inspr/examples/primes/printer:latest $(PROJ_SRC)
	docker push gcr.io/insprlabs/inspr/examples/primes/printer:latest 

clean:
	docker image prune -f
