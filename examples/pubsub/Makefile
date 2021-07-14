build:
	docker build -t gcr.io/insprlabs/inspr/examples/pubsubapi:latest -f api/Dockerfile ../..
	docker push gcr.io/insprlabs/inspr/examples/pubsubapi:latest
	docker build -t gcr.io/insprlabs/inspr/examples/discordpub:latest -f discord/Dockerfile ../..
	docker push gcr.io/insprlabs/inspr/examples/discordpub:latest
	docker build -t gcr.io/insprlabs/inspr/examples/slackpub:latest -f slack/Dockerfile ../..
	docker push gcr.io/insprlabs/inspr/examples/slackpub:latest