from inspr import *
import sys

PONG_INPUT_CHANNEL = "ponginput"
PONG_OUTPUT_CHANNEL = "pongoutput"

def main():
    client = Client()
    msg = "Pong!"

    def readPingAndSendPong(data):

        if data == 'Ping!':
            print(data, file=sys.stderr)
        else:
            print('Not received Ping', file=sys.stderr)
        
        try:
            client.writeMessage(PONG_OUTPUT_CHANNEL, msg)
        except:
            raise Exception


    try:
        client.writeMessage(PONG_OUTPUT_CHANNEL, msg)
    except:
        print("An error has occured", file=sys.stderr)
        return

    client.handleChannel(PONG_INPUT_CHANNEL, readPingAndSendPong)
    client.run()

if __name__ == "__main__":
    main()