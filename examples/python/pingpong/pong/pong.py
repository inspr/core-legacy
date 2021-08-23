from inspr import *
import sys

PONG_INPUT_CHANNEL = "ponginput"
PONG_OUTPUT_CHANNEL = "pongoutput"

def main():
    client = Client()
    msg = "Pong!"

    @client.handle_channel(PONG_INPUT_CHANNEL)
    def read_ping_and_send_pong(data):
        if data == 'Ping!':
            print(data, file=sys.stderr)
        else:
            print('Not received Ping', file=sys.stderr)

        try:
            client.write_message(PONG_OUTPUT_CHANNEL, msg)
            return Response(status=200)
        except:
            raise Exception
        

    try:
        client.write_message(PONG_OUTPUT_CHANNEL, msg)
    except:
        print("An error has occured", file=sys.stderr)
        return

    client.run()

if __name__ == "__main__":
    main()
