using UnityEngine;

public static class Functions
{

    // ALL PURE UTIL FUNCTIONS

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

    public static Vector3 GetRandomGamePosition(GameState gameState) {
        float randX = Random.Range(
            Functions.GetBound(gameState, Vector3.down),
            Functions.GetBound(gameState, Vector3.up)
        );
        float randY = Random.Range(
            Functions.GetBound(gameState, Vector3.left),
            Functions.GetBound(gameState, Vector3.right)
        );
        return new Vector3(randX, randY, 0);

    }
    
}
