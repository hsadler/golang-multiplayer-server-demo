using System;

[Serializable]
public class ClientMessagePlayerEnter
{

    public string messageType = Constants.CLIENT_MESSAGE_TYPE_PLAYER_ENTER;
    public Player player;

    public ClientMessagePlayerEnter(Player playerModel)
    {
        this.player = playerModel;
    }

}

[Serializable]
public class ClientMessagePlayerExit
{

    public string messageType = Constants.CLIENT_MESSAGE_TYPE_PLAYER_EXIT;
    public string playerId;

    public ClientMessagePlayerExit(string playerId)
    {
        this.playerId = playerId;
    }

}

[Serializable]
public class ClientMessagePlayerPosition
{

    public string messageType = Constants.CLIENT_MESSAGE_TYPE_PLAYER_POSITION;
    public string playerId;
    public Position position;

    public ClientMessagePlayerPosition(string playerId, Position position)
    {
        this.playerId = playerId;
        this.position = position;
    }

}

[Serializable]
public class ClientMessagePlayerEatFood
{

    public string messageType = Constants.CLIENT_MESSAGE_TYPE_PLAYER_EAT_FOOD;
    public string playerId;
    public string foodId;

    public ClientMessagePlayerEatFood(string playerId, string foodId)
    {
        this.playerId = playerId;
        this.foodId = foodId;
    }

}

[Serializable]
public class ClientMessagePlayerEatPlayer
{

    public string messageType = Constants.CLIENT_MESSAGE_TYPE_PLAYER_EAT_PLAYER;
    public string playerId;
    public string otherPlayerId;

    public ClientMessagePlayerEatPlayer(string playerId, string otherPlayerId)
    {
        this.playerId = playerId;
        this.otherPlayerId = otherPlayerId;
    }

}

[Serializable]
public class ClientMessageMineDamagePlayer
{

    public string messageType = Constants.CLIENT_MESSAGE_TYPE_MINE_DAMAGE_PLAYER;
    public string playerId;
    public string mineId;

    public ClientMessageMineDamagePlayer(string playerId, string mineId)
    {
        this.playerId = playerId;
        this.mineId = mineId;
    }

}
