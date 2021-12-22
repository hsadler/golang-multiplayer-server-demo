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

[Serializable]
public class ServerMessageSecondsToNextRoundStart
{
    public string messageType;
    public int seconds;
}

[Serializable]
public class ServerMessageSecondsToCurrentRoundEnd
{
    public string messageType;
    public int seconds;
}
