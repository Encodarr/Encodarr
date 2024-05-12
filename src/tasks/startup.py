import asyncio
from datetime import datetime
import json
import os
from fastapi.staticfiles import StaticFiles
from src.api.controllers.settings_controller import get_all_settings
from src.api.controllers.system_controller import set_system
from src.services.encode_service import encode_service
from src.services.logging_service import start_logger
from src.services.metadata_service import metadata_service
from src.services.scan_service import scan_service
from src.services.watchdog_service import start_watchdog
from src.utils.folders import get_root_folder, verify_folders
from fastapi.openapi.utils import get_openapi


def custom_openapi():
    from src.main import app
    if app.openapi_schema:
        return app.openapi_schema
    openapi_schema = get_openapi(
        title="Transfigurr API",
        version="V1",
        routes=app.routes,
    )
    app.openapi_schema = openapi_schema
    return app.openapi_schema


async def write_api():
    openapi_schema = custom_openapi()
    openapi_path = 'src/Transfigurr.API.V1'
    if not os.path.exists(openapi_path):
        os.makedirs(openapi_path)
    with open(openapi_path + '/openapi.json', 'w') as file:
        json.dump(openapi_schema, file)


async def write_uptime_to_db():
    uptime_date = datetime.now()
    await set_system({'id': 'start_time', 'value': uptime_date})


async def startup():
    from src.main import app
    await verify_folders()
    app.mount("/static", StaticFiles(directory="frontend/dist"), name="static")
    log_level = (await get_all_settings()).get('log_level', '')
    start_logger(log_level)
    await write_uptime_to_db()
    await write_api()
    asyncio.create_task(scan_service.enqueue_all())
    asyncio.create_task(scan_service.process())
    asyncio.create_task(metadata_service.process())
    asyncio.create_task(encode_service.process())
    start_watchdog(await get_root_folder() + '/series', 'series')
    start_watchdog(await get_root_folder() + '/movies', 'movies')
