from aiogram import Router ,F ,types
from aiogram.types import Message
from aiogram.fsm.context import FSMContext
from states import ProfileCreationState ,ProfileUpdateState,user_ids
from pb2s.profile_pb2 import Profile
from aiogram.filters import Command
import kb
from server_calls import CreateRequest ,UpdateRequest,ReadRequest
CRU_router = Router()
from funcs import show_profile




#Create
@CRU_router.message(F.text == "Create profile.")
async def profile_creattion_handler(msg:Message,state:FSMContext):
     await state.update_data(is_updating=False)
     await create_profile_form(msg,state)
     
     return

async def create_profile_form(msg:Message,state: FSMContext):
     await msg.answer("Enter your name:",reply_markup=types.ReplyKeyboardRemove())
     await state.set_state(ProfileCreationState.name)

@CRU_router.message(ProfileCreationState.name)
async def process_name(msg: Message, state: FSMContext):
    await state.update_data(name=msg.text)
    await msg.answer("Enter your age:")
    await state.set_state(ProfileCreationState.age)

@CRU_router.message(ProfileCreationState.age)
async def process_age(msg: Message, state: FSMContext):
     try:
          int(msg.text)
     except:
          await msg.answer("Enter number! Try again:")
          return
     
     if int(msg.text) < 16 or int(msg.text) > 100:
          await msg.answer("Age must be number between 16 and 100! Try again:")
          return
     
     await state.update_data(age=msg.text)
     await msg.answer("Enter your description:")
     await state.set_state(ProfileCreationState.description)

@CRU_router.message(ProfileCreationState.description)
async def process_description(msg: Message, state: FSMContext):
     await state.update_data(description=msg.text)
     await msg.answer("Send your photo:")
     await state.set_state(ProfileCreationState.pfp_id)

@CRU_router.message(ProfileCreationState.pfp_id)
async def procces_pfp(msg:Message, state: FSMContext):
     pfp = msg.photo[-1]
     await state.update_data(pfp_id = pfp.file_id)
     await msg.answer("Enter your location:")
     await state.set_state(ProfileCreationState.location)


@CRU_router.message(ProfileCreationState.location)
async def process_location(msg: Message, state: FSMContext):
    await state.update_data(location=msg.text)
    user_data = await state.get_data()
    profile = Profile(
                    ID=msg.chat.id, 
                    name=user_data['name'],
                    age=int(user_data['age']),
                    description=user_data['description'],
                    location =user_data['location'],
                    pfp_id=user_data["pfp_id"],
                    )

    try:
        if user_data['is_updating']:
            response = UpdateRequest(profile)  # Update the profile
            await msg.answer("Profile updated successfully.")
        else:
          response = CreateRequest(profile) 
          user_ids.add(msg.chat.id)
          await msg.answer("Profile created successfully.")
          await msg.answer("Your profile: ")
          await show_profile(msg,profile, kb.get_continue_keyboard )
     
    except Exception as e:
        await msg.answer("Something went wrong")
        print(e)
    
    await state.clear()





#Update
@CRU_router.message(F.text == "Change profile.")
async def profile_update_handler(msg:Message,state:FSMContext):
     await update_profile(msg)

async def update_profile(msg:Message):
     await msg.answer(
'''What you want to change?
1.Name
2.Age
3.Description
4.Photo
5.Location
''' , reply_markup=kb.get_account_update_keyboard())
     

@CRU_router.message(F.text == "1.")
async def update_name(msg:Message ,state:FSMContext):
     await msg.answer("Enter your new name:")
     await state.set_state(ProfileUpdateState.new_name)
@CRU_router.message(ProfileUpdateState.new_name)
async def update_name_handler(msg:Message ,state:FSMContext):
     new_name = msg.text
     await state.clear()
     profile = ReadRequest(msg.chat.id)
     profile.name = new_name
     try:
          updateResponse = UpdateRequest(profile)
          await msg.answer("Profile name updated successfully.")
          await msg.answer("Your profile: ")
          await show_profile(msg,profile, kb.get_continue_keyboard )
     except Exception as e:
          await msg.answer("Something went wrong")
          print(e)


@CRU_router.message(F.text == "2.")
async def update_age(msg:Message ,state:FSMContext):
     await msg.answer("Enter your new age:")
     await state.set_state(ProfileUpdateState.new_age)
@CRU_router.message(ProfileUpdateState.new_age)
async def update_name_handler(msg:Message ,state:FSMContext):
     new_age = msg.text
     await state.clear()
     profile = ReadRequest(msg.chat.id)
     try:
          int(new_age)
     except:
          await msg.answer("Enter number! Try again:")
          return

     if int(new_age) < 16 or int(new_age) > 100:
          await msg.answer("Age must be number between 16 and 100! Try again:")
          return
     profile.age = int(new_age)
     try:
          updateResponse = UpdateRequest(profile)
          await msg.answer("Profile age updated successfully.")
          await msg.answer("Your profile: ")
          await show_profile(msg,profile, kb.get_continue_keyboard )
     except Exception as e:
          await msg.answer("Something went wrong")
          print(e)




@CRU_router.message(F.text == "3.")
async def update_description(msg:Message ,state:FSMContext):
     await msg.answer("Enter your new description:")
     await state.set_state(ProfileUpdateState.new_description)
@CRU_router.message(ProfileUpdateState.new_description)
async def update_description_handler(msg:Message ,state:FSMContext):
     new_description = msg.text
     await state.clear()
     profile = ReadRequest(msg.chat.id)
     profile.description = new_description
     try:
          updateResponse = UpdateRequest(profile)
          await msg.answer("Profile description updated successfully.")
          await msg.answer("Your profile: ")
          await show_profile(msg,profile, kb.get_continue_keyboard )
     except Exception as e:
          await msg.answer("Something went wrong")
          print(e)



@CRU_router.message(F.text == "4.")
async def update_photo(msg:Message ,state:FSMContext):
     await msg.answer("Send new photo:")
     await state.set_state(ProfileUpdateState.new_pfp)
@CRU_router.message(ProfileUpdateState.new_pfp)
async def update_photo_handler(msg:Message ,state:FSMContext):
     new_pfp = msg.photo[-1]
     await state.clear()
     profile = ReadRequest(msg.chat.id)
     profile.pfp_id = new_pfp.file_id
     try:
          updateResponse = UpdateRequest(profile)
          await msg.answer("Profile photo updated successfully.")
          await msg.answer("Your profile: ")
          await show_profile(msg,profile, kb.get_continue_keyboard )
     except Exception as e:
          await msg.answer("Something went wrong")
          print(e)



@CRU_router.message(F.text == "5.")
async def update_location(msg:Message ,state:FSMContext):
     await msg.answer("Enter your new location:")
     await state.set_state(ProfileUpdateState.new_location)
@CRU_router.message(ProfileUpdateState.new_location)
async def update_location_handler(msg:Message ,state:FSMContext):
     new_location = msg.text
     await state.clear()
     profile = ReadRequest(msg.chat.id)
     profile.location = new_location
     try:
          updateResponse = UpdateRequest(profile)
          await msg.answer("Profile location updated successfully.")
          await msg.answer("Your profile: ")
          await show_profile(msg,profile, kb.get_continue_keyboard )
     except Exception as e:
          await msg.answer("Something went wrong")
          print(e)









@CRU_router.message(Command("create_test_account"))
async def create_test_account(msg:Message):
     profile = Profile(
               ID=msg.chat.id, 
               name="name",
               age=19,
               description="sbojin",
               location ="ajv",
               pfp_id="AgACAgIAAxkBAAIRmmYc9PxPe7A_FZdvuOsZ00Qoe2wtAAJ41zEbGdrpSFaWZvmxD-UFAQADAgADeAADNAQ",
               )
     response = CreateRequest(profile)
     user_ids.add(msg.chat.id)
     return 