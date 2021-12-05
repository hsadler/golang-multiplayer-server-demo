using System;
using System.Collections.Generic;

[Serializable]
public class GameState
{

    public int mapHeight;
    public int mapWidth;
    public List<Player> players;
    public List<Food> foods;
    public List<Mine> mines;

}

[Serializable]
public class Player
{

    public string id;
    public Position position;

    public Player(string id, Position position)
    {
        this.id = id;
        this.position = position;
    }

}

[Serializable]
public class Food
{

    public string id;
    public Position position;

    public Food(string id, Position position)
    {
        this.id = id;
        this.position = position;
    }

}

[Serializable]
public class Mine
{

    public string id;
    public Position position;

    public Mine(string id, Position position)
    {
        this.id = id;
        this.position = position;
    }

}

[Serializable]
public class Position
{

    public float x;
    public float y;

    public Position(float x, float y)
    {
        this.x = x;
        this.y = y;
    }

}
