build:
	docker build -t gcr.io/red-inspr/inspr/example/pubsubapi:latest -f api/Dockerfile ../..
	docker push gcr.io/red-inspr/inspr/example/pubsubapi:latest
	docker build -t gcr.io/red-inspr/inspr/example/discordpub:latest -f discord/Dockerfile ../..
	docker push gcr.io/red-inspr/inspr/example/discordpub:latest
	docker build -t gcr.io/red-inspr/inspr/example/slackpub:latest -f slack/Dockerfile ../..
	docker push gcr.io/red-inspr/inspr/example/slackpub:latest