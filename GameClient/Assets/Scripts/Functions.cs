using System;
using UnityEngine;

public static class Functions
{

    public static float GetBound(GameState gameState, Vector3 direction)
    {
        if (direction == Vector3.up)
        {
            return gameState.mapHeight / 2;
        }
        else if (direction == Vector3.down) {
            return -gameState.mapHeight / 2;
        }
        else if (direction == Vector3.left)
        {
            return -gameState.mapWidth / 2;
        }
        else if (direction == Vector3.right)
        {
            return gameState.mapWidth / 2;
        }
        return 0;
    }
    
}
