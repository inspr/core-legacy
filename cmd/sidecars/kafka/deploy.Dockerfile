# copies builded kafka sidecar to deploy folder
FROM gcr.io/red-inspr/inspr AS build-env
WORKDIR /deploy
COPY --from=build-env /inspr/cmd/sidecars/kafka kafka
ENTRYPOINT ./kafka
