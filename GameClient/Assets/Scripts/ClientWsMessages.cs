using System;

[Serializable]
public class ClientMessagePlayerEnter
{

    public string messageType = "CLIENT_MESSAGE_TYPE_PLAYER_ENTER";
    public Player player;

    public ClientMessagePlayerEnter(Player playerModel)
    {
        this.player = playerModel;
    }

}

[Serializable]
public class ClientMessagePlayerPosition
{

    public string messageType = "CLIENT_MESSAGE_TYPE_PLAYER_POSITION";
    public Player player;

    public ClientMessagePlayerPosition(Player playerModel)
    {
        this.player = playerModel;
    }

}
