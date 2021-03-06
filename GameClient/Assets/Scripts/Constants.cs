using System;

public static class Constants
{
    // client message types
    public const string CLIENT_MESSAGE_TYPE_PLAYER_ENTER = "CLIENT_MESSAGE_TYPE_PLAYER_ENTER";
    public const string CLIENT_MESSAGE_TYPE_PLAYER_EXIT = "CLIENT_MESSAGE_TYPE_PLAYER_EXIT";
    public const string CLIENT_MESSAGE_TYPE_PLAYER_POSITION = "CLIENT_MESSAGE_TYPE_PLAYER_POSITION";
    public const string CLIENT_MESSAGE_TYPE_PLAYER_EAT_FOOD = "CLIENT_MESSAGE_TYPE_PLAYER_EAT_FOOD";
    public const string CLIENT_MESSAGE_TYPE_PLAYER_EAT_PLAYER = "CLIENT_MESSAGE_TYPE_PLAYER_EAT_PLAYER";
    public const string CLIENT_MESSAGE_TYPE_MINE_DAMAGE_PLAYER = "CLIENT_MESSAGE_TYPE_MINE_DAMAGE_PLAYER";
    // server message types
    public const string SERVER_MESSAGE_TYPE_GAME_STATE = "SERVER_MESSAGE_TYPE_GAME_STATE";
    public const string SERVER_MESSAGE_TYPE_PLAYER_ENTER = "SERVER_MESSAGE_TYPE_PLAYER_ENTER";
    public const string SERVER_MESSAGE_TYPE_PLAYER_EXIT = "SERVER_MESSAGE_TYPE_PLAYER_EXIT";
    public const string SERVER_MESSAGE_TYPE_PLAYER_STATE_UPDATE = "SERVER_MESSAGE_TYPE_PLAYER_STATE_UPDATE";
    public const string SERVER_MESSAGE_TYPE_FOOD_STATE_UPDATE = "SERVER_MESSAGE_TYPE_FOOD_STATE_UPDATE";
    public const string SERVER_MESSAGE_TYPE_MINE_STATE_UPDATE = "SERVER_MESSAGE_TYPE_MINE_STATE_UPDATE";
    public const string SERVER_MESSAGE_TYPE_SECONDS_TO_NEXT_ROUND_START = "SERVER_MESSAGE_TYPE_SECONDS_TO_NEXT_ROUND_START";
    public const string SERVER_MESSAGE_TYPE_SECONDS_TO_CURRENT_ROUND_END = "SERVER_MESSAGE_TYPE_SECONDS_TO_CURRENT_ROUND_END";
    public const string SERVER_MESSAGE_TYPE_ROUND_RESULT = "SERVER_MESSAGE_TYPE_ROUND_RESULT";
}
