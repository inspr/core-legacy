VERSION=$(git describe --always)

gsutil rsync bin gs://inspr-cli/$VERSION

release-cli create --name "Release $VERSION" \
        --tag-name $VERSION \
        --ref "$CI_COMMIT_SHA" \
        --assets-link "{\"name\":\"insprcli-linux-arm64\",
                        \"url\":\"https://console.cloud.google.com/storage/browser/inspr-cli/bin/insprcli-linux-arm64-$VERSION\"}" \

        --assets-link "{\"name\":\"insprcli-linux-amd64\",
                        \"url\":\"https://console.cloud.google.com/storage/browser/inspr-cli/bin/insprcli-linux-amd64-$VERSION\"}" \

        --assets-link "{\"name\":\"insprcli-windows-amd64\",
                        \"url\":\"https://console.cloud.google.com/storage/browser/inspr-cli/bin/insprcli-windows-amd64-$VERSION\"}" \

        --assets-link "{\"name\":\"insprcli-windows-386\",
                        \"url\":\"https://console.cloud.google.com/storage/browser/inspr-cli/bin/insprcli-windows-386-$VERSION\"}" \
                        
        --assets-link "{\"name\":\"insprcli-darwin-amd64\",
                        \"url\":\"https://console.cloud.google.com/storage/browser/inspr-cli/bin/insprcli-darwin-amd64-$VERSION\"}"