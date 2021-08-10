from inspr import *
import sys

PONG_INPUT_CHANNEL = "ponginput"
PONG_OUTPUT_CHANNEL = "pongoutput"

def main():
    client = Client()
    msg = "Pong!"

    def readPingAndSendPong(data):
        print("data =", data, file=sys.stderr)
        print("Ping!", file=sys.stderr)
        
        try:
            client.writeMessage(PONG_OUTPUT_CHANNEL, msg)
        except:
            raise Exception


    try:
        client.writeMessage(PONG_OUTPUT_CHANNEL, msg)
    except:
        print("An error has occured")
        return

    client.handleChannel(PONG_INPUT_CHANNEL, readPingAndSendPong)
    client.run()

if __name__ == "__main__":
    main()