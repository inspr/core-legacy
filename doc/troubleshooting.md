# Troubleshooting

## I can't connect to Insprd

- Check if your CLI is correctly configured.
- Repeat the steps in [the getting started tutorial](./Getting Started)

## My dApp only reads one message from a channel

- Check if there is no type conflict between the dApp and the channel
- Check if you are commiting the message after reading it
- Check if your dApp is reading messages continuously

## I can't apply a folder because the CLI complains that there is no channel named X, when it is defined inside the folder

This happens because the Inspr CLI tries to apply dapps first, so that you have your contexts already created.

To fix it, either apply your channels and channel types before your dApps or create a base dApp that contains all channels and channel types, and instantiate your other dApps inside the context defined previously.
