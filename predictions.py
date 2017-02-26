import asyncio

import aiohttp
import async_timeout

import config


async def fetch_stop(session, stop_id):
    params = {
        'key': config.data['api_key'],
        'format': config.data['api_format'],
        'stpid': stop_id
    }
    with async_timeout.timeout(config.data['polling_timeout']):
        async with session.get(config.data['api_route'], params=params) as r:
            return await r.json()


async def fetch_all_stops(loop):
    tasks = []
    async with aiohttp.ClientSession(loop=loop) as session:
        for stop_id in config.data['stop_ids']:
            tasks.append(asyncio.ensure_future(fetch_stop(session, stop_id)))
        return await asyncio.gather(*tasks)
