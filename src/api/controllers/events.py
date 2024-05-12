from src.models.log import Log
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.future import select
from src.utils.db import engine, instance_to_dict


async def get_all_logs():
    async with AsyncSession(engine) as async_session:
        res = await async_session.execute(select(Log))
        return [instance_to_dict(log) for log in res.scalars().all()]


async def get_log(log_id):
    async with AsyncSession(engine) as async_session:
        res = await async_session.execute(select(Log).where(Log.id == log_id))
        return (instance_to_dict(res.scalars().first()))
