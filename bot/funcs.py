from aiogram.fsm.context import FSMContext
from aiogram.types import Message
import asyncio
from server_calls import GetNextProfileRequest , LikeRequest ,GetLikesRequest,ReadRequest
from aiogram import types
from states import CurrentDisplayedProfile, Likes
import kb
import asyncio


async def show_profile(msg: Message, profile, get_keyboard):
    caption = f"{profile.name}, {profile.location}, {profile.age} - {profile.description}"
    await msg.bot.send_photo(
        chat_id=msg.chat.id,
        photo=profile.pfp_id,
        caption=caption,
        reply_markup=get_keyboard()
    )


async def like(msg:Message,tgID):
     try:
          responce = LikeRequest(msg.chat.id, tgID)
     except:
          await msg.answer("Something went wrong while like request!" ,
                           reply_markup= types.ReplyKeyboardRemove())
          return
     await msg.answer("Like sent!",reply_markup= types.ReplyKeyboardRemove())



async def show_next_profile(msg: Message, state: FSMContext):
    data = await state.get_data()
    if data.get('likesDisplayed') == True:
        likes = data.get('likes', [])
        
        if likes:
            first_like = likes[0]
            like_profile = ReadRequest(first_like)
            await state.update_data(profile = like_profile)
            await state.set_state(CurrentDisplayedProfile.profile)  
            await show_profile(msg,like_profile,kb.get_main_keyboard)
            likes.pop(0)
            await state.update_data(likes=likes) 
                
        else:
               await state.update_data(likesDisplayed=False)
               await msg.answer("There are no likes to display.", reply_markup=kb.get_continue_keyboard())
    else:
        try:
            next_profile = GetNextProfileRequest(msg.chat.id)
            await state.update_data(profile=next_profile)
            await state.set_state(CurrentDisplayedProfile.profile)  
            await show_profile(msg,next_profile, kb.get_main_keyboard)
        except Exception as e:
            await msg.answer(f"Something went wrong while getting the next profile! {str(e)}")


async def get_all_likes():
    try:
        response = GetLikesRequest()
    except Exception as e:
        print(e)
    return response.likes


async def get_users_likes(UID,allLikes):
    users_likes = []
    for like in allLikes:
        #server side is wrong
        #user_id is liked user
        #and liked_user_id is the person who liked
        if like.user_id == UID and like.liked_user_id not in users_likes:
            users_likes.append(like.liked_user_id)
    print([UID, users_likes])
    return users_likes

async def show_likes(bot , user_id, state: FSMContext,likes):
        if not likes:
            return 
        else:
            await state.set_state(Likes.likes)
            await state.update_data(likes=likes, likesDisplayed=True)
            await bot.send_message(chat_id=user_id, text=f"You have {len(likes)} likes.", reply_markup=kb.get_likes_keyboard())


async def match(msg:Message, liked_id):
    await msg.answer(f"Have a nice chat :) <a href='tg://user?id={liked_id}'>click here</a>")
    profile = ReadRequest(msg.chat.id)
    caption1 = f"{profile.name}, {profile.location}, {profile.age} - {profile.description}"
    await msg.bot.send_photo(
        chat_id=liked_id,
        photo=profile.pfp_id,
        caption=caption1,
    )
    caption2 =  f"This person liked u back, have a nice chat :)  <a href='tg://user?id={msg.chat.id}'>click here</a>"
    await msg.bot.send_message(chat_id=liked_id, text=caption2)