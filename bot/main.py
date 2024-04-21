import asyncio
import logging
from CRU_handler import CRU_router
from menu_handler import menu_router
from aiogram import Bot, Dispatcher
from aiogram.enums.parse_mode import ParseMode
from aiogram.fsm.storage.memory import MemoryStorage
from aiogram.fsm.context import FSMContext,StorageKey
import config
from funcs import show_likes,get_all_likes ,get_users_likes
from states import get_user_ids
#message gets the giveuser
#message is with profile getuser and not giveuser

async def get_likes_timer(bot:Bot ,dp:Dispatcher):
    user_ids = get_user_ids()
    while True:
        likes = await get_all_likes()
        if likes:
            for user_id in user_ids:
                key = StorageKey(bot_id=bot.id, chat_id=user_id, user_id=user_id)
                state = FSMContext(dp.storage, key)
                users_likes = await get_users_likes(user_id,likes)
                await show_likes(bot,user_id,state, users_likes)
        print("Likes iterated")
        await asyncio.sleep(30)


async def main():
    bot = Bot(token=config.BOT_TOKEN, parse_mode=ParseMode.HTML)
    dp = Dispatcher(storage=MemoryStorage())
    dp.include_router(menu_router)
    dp.include_router(CRU_router)
    await bot.delete_webhook(drop_pending_updates=True)
    asyncio.create_task(get_likes_timer(bot,dp))
    await dp.start_polling(bot, allowed_updates=dp.resolve_used_update_types())



if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    asyncio.run(main())
