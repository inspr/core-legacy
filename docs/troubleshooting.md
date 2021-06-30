# Troubleshooting

## I can't connect to Insprd
- Check if your Inspr CLI is correctly configured. You can do so by running the following command:
```
insprctl config list
```
- Repeat the steps in [the getting started tutorial](readme.md).

## My dApp only reads one message from a Channel
- Check if there is no type conflict between the dApp and the Channel.
- Check if your dApp is reading messages continuously.

## I can't apply a folder because the CLI complains that there is no Channel named X, when it is defined inside the folder  
This happens because the Inspr CLI tries to apply dApps first, so that you have your contexts already created.

To fix it, either apply your Channels and Types before your dApps or create a base dApp that contains all channels and Types, and instantiate your other dApps inside the context defined previously.

## UID Provider k8s deployment is crashing
- Once you've installed Insprd, remember to run the command `insprctl cluster init <password>` so the admin user is configured and its token generated.
- Place the token that was prompted in your terminal in `secret.adminToken` in the file `build/uidp_helm/values.yaml`.
- If you are reinstalling Insprd, remove the secret called `redisprivatekey` from your cluster before doing so.

## Unable to login with the UID Provider
- If you have UID Provider properly installed and the `create` or `login` command aren't working try exporting the following environment variable with your cluster's address:
```zsh
export INPROV_PROVIDER_URL="http://<CLUSTER-ADDRESS>"
```

## Avro Schema file not being resolved when applying a Type
When the `schema` field of a Type is a path to a `.avsc` (or `.schema`) file, that path is resolved based on the user's current directory. Basically, if you have the following structure:
```
pingpong
└── yamls
    └── types
        ├── type.yaml
        └── schema.avsc
```
And the `schema` field of `type.yaml` is `yamls/types/schema.avsc`, you'd have to be in the folder `/pingpong` so the schema is properly resolved. If you are within `/pinpong/yamls` and try to apply the Type, its schema won't be resolved and will be set as "yamls/types/schema.avsc".