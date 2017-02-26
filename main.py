import asyncio
from time import sleep

import uvloop
import zmq

import config
from predictions import fetch_all_stops


def main():
    context = zmq.Context()
    socket = context.socket(zmq.PUB)
    socket.bind("tcp://*:7350")

    asyncio.set_event_loop_policy(uvloop.EventLoopPolicy())
    loop = asyncio.get_event_loop()

    while True:
        try:
            all_predictions = loop.run_until_complete(fetch_all_stops(loop=loop))
        except Exception:
            break
        else:
            for p in all_predictions:
                socket.send_json(p)
            sleep(config.data['polling_interval'])

    loop.close()


if __name__ == "__main__":
    main()
