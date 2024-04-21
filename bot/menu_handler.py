from aiogram import types, F, Router
from aiogram.fsm.context import FSMContext
from aiogram.types import Message
from aiogram.filters import Command
from server_calls import ReadRequest
from funcs import like ,show_next_profile ,show_profile,match
import kb


menu_router = Router()


@menu_router.message(Command("start"))
async def start_handler(msg:Message):
     try:
          readResponce = ReadRequest(msg.chat.id)
     except:
          await msg.answer(
'''It seems that you are new here.
You can create your profile. Press the button below.''',
                           reply_markup=kb.get_greet_keyboard())
     else:
          await msg.answer("Welcome back")
          await msg.answer("Here is your profile:")
          await show_profile(msg,readResponce,kb.get_continue_keyboard)



@menu_router.message(F.text == "Start searching." )
async def start_searching_handler(msg:Message,state:FSMContext):
     await show_next_profile(msg,state)


@menu_router.message(F.text =="🍻")
async def like_handler(msg:Message,state:FSMContext):
     data = await state.get_data()
     if data.get('likesDisplayed') == True:
          await match(msg,data['profile'].ID)
     else:
          await like(msg, data['profile'].ID)
     await show_next_profile(msg,state)
     # await like(msg, data['profile'].ID)
     # await show_next_profile(msg,state)

@menu_router.message(F.text =="👎",)
async def dislike_handler(msg:Message,state:FSMContext):
     await show_next_profile(msg,state)

@menu_router.message(F.text =="😴")
async def pause_searching_handler(msg:Message):
     await msg.answer("Pause searching people...",reply_markup= kb.get_continue_keyboard())

@menu_router.message(F.text == "See likes.")
async def see_likes_handler(msg:Message, state:FSMContext):
     await show_next_profile(msg,state)
