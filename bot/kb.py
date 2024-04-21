from aiogram.types import  ReplyKeyboardMarkup,KeyboardButton

def get_greet_keyboard():
    keyboard = [[KeyboardButton(text = "Create profile.")]]
    return ReplyKeyboardMarkup(keyboard=keyboard ,resize_keyboard=True, one_time_keyboard=True)

def get_main_keyboard():
    keyboard = [
        [KeyboardButton(text = "🍻"),
         KeyboardButton(text = "👎"),
         KeyboardButton(text = "😴")
         ]
    ]
    return ReplyKeyboardMarkup(keyboard=keyboard,resize_keyboard=True, one_time_keyboard=True)


def get_continue_keyboard():
    keyboard = [
        [KeyboardButton(text="Start searching."),
         KeyboardButton(text="Change profile.")]
    ]
    return ReplyKeyboardMarkup(keyboard=keyboard,resize_keyboard=True, one_time_keyboard=True)

def get_likes_keyboard():
    keyboard = [
        [KeyboardButton(text="See likes.")]
    ]
    return ReplyKeyboardMarkup(keyboard=keyboard,resize_keyboard=True, one_time_keyboard=True)

def get_account_update_keyboard():
    keyboard = [
        [
        KeyboardButton(text="1."),
        KeyboardButton(text="2."),
        KeyboardButton(text="3."),
        KeyboardButton(text="4."),
        KeyboardButton(text="5."),
        ]
    ]
    return ReplyKeyboardMarkup(keyboard=keyboard,resize_keyboard=True, one_time_keyboard=True)
