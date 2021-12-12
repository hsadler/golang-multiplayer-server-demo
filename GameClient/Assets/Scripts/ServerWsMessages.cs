using System;

[Serializable]
public class ServerMessageGeneric
{
    public string messageType;
}

[Serializable]
public class ServerMessageGameState
{
    public string messageType;
    public GameState gameState;
}

[Serializable]
public class ServerMessagePlayerEnter
{
    public string messageType;
    public Player player;
}

[Serializable]
public class ServerMessagePlayerExit
{
    public string messageType;
    public string playerId;
}

[Serializable]
public class ServerMessagePlayerUpdate
{
    public string messageType;
    public Player player;
}

[Serializable]
public class ServerMessageFoodUpdate
{
    public string messageType;
    public Food food;
}

[Serializable]
public class ServerMessageMineUpdate
{
    public string messageType;
    public Mine mine;
}
