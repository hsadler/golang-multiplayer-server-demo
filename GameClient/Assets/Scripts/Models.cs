﻿using System;
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
    public bool active;
    public string name;
    public Position position;
    public int size;

    public Player(string id, bool active, string name, Position position, int size)
    {
        this.id = id;
        this.active = active;
        this.name = name;
        this.position = position;
        this.size = size;
    }

}

[Serializable]
public class Food
{

    public string id;
    public bool active;
    public Position position;
    public int size;

    public Food(string id, bool active, Position position, int size)
    {
        this.id = id;
        this.active = active;
        this.position = position;
        this.size = size;
    }

}

[Serializable]
public class Mine
{

    public string id;
    public bool active;
    public Position position;
    public int size;

    public Mine(string id, bool active, Position position, int size)
    {
        this.id = id;
        this.active = active;
        this.position = position;
        this.size = size;
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
