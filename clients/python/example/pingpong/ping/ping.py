from inspr import *
import sys

PING_INPUT_CHANNEL = "pinginput"
PING_OUTPUT_CHANNEL = "pingoutput"

def main():
    client = Client()
    msg = "Ping!"

    def readPongAndSendPing(data):

        if data == 'Pong!':
            print(data, file=sys.stderr)
        else:
            print('Not received Pong', file=sys.stderr)

        try:
            client.writeMessage(PING_OUTPUT_CHANNEL, msg)
        except:
            raise Exception


    try:
        client.writeMessage(PING_OUTPUT_CHANNEL, msg)
    except:
        print("An error has occured", file=sys.stderr)
        return

    client.handleChannel(PING_INPUT_CHANNEL, readPongAndSendPing)
    client.run()

if __name__ == "__main__":
    main()