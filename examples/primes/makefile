# docker images
build: ## Build the container
	# generator
	docker build -f generator/Dockerfile -t gcr.io/red-inspr/inspr/examples/primes/generator:latest ../..
	docker push gcr.io/red-inspr/inspr/examples/primes/generator:latest 

	# filter
	docker build -f  filter/Dockerfile -t gcr.io/red-inspr/inspr/examples/primes/filter:latest ../..
	docker push gcr.io/red-inspr/inspr/examples/primes/filter:latest 

	# printer
	docker build -f printer/Dockerfile -t gcr.io/red-inspr/inspr/examples/primes/printer:latest ../..
	docker push gcr.io/red-inspr/inspr/examples/primes/printer:latest 

clean:
	docker image prune -f
