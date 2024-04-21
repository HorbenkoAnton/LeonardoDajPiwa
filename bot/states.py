from aiogram.filters.state import StatesGroup, State

class ProfileCreationState(StatesGroup):
    name = State()
    age = State()
    location = State()
    description = State()
    pfp_id = State()

class ProfileUpdateState(StatesGroup):
    new_name = State()
    new_age = State()
    new_location = State()
    new_description = State()
    new_pfp = State()


class CurrentDisplayedProfile(StatesGroup):
    profile = State()

class Likes(StatesGroup):
    likes = State()
    likesDisplayed = State()

user_ids = set()

def get_user_ids():
    return user_ids