
PROJ_SRC = ../../.

# docker images

build: ## Build the container
	# generator
	docker build -f ./generator/Dockerfile -t gcr.io/red-inspr/inspr/examples/primes/generator:latest $(PROJ_SRC)
	docker push gcr.io/red-inspr/inspr/examples/primes/generator:latest 
	# printer
	docker build -f ./printer/Dockerfile -t gcr.io/red-inspr/inspr/examples/primes/printer:latest $(PROJ_SRC)
	docker push gcr.io/red-inspr/inspr/examples/primes/printer:latest 

clean:
	docker image prune -f
