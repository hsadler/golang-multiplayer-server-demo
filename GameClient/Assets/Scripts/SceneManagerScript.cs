using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using WebSocketSharp;

public class SceneManagerScript : MonoBehaviour
{

    // prefabs
    public GameObject playerPrefab;
    public GameObject foodPrefab;
    public GameObject minePrefab;
    public GameObject wallPrefab;

    // game object refs
    private GameObject mainPlayerGO;
    public GameObject mainCameraGO;
    public GameObject giveNameUI;

    private WebSocket ws;
    // local game server
    private string gameServerUrl = "ws://localhost:5000";
    // heroku game server
    //private string gameServerUrl = "ws://golang-multiplayer-server-demo.herokuapp.com/:80";

    // game state
    private GameState gameState;
    private bool gameStateInitialized = false;

    // main player state
    private Player mainPlayerModel;

    // env refs
    private List<GameObject> wallGOs = new List<GameObject>();
    private IDictionary<string, GameObject> playerIdToOtherPlayerGO;
    private IDictionary<string, GameObject> foodIdToFoodGO;
    private IDictionary<string, GameObject> mineIdToMineGO;

    private Queue<string> gameServerMessageQueue = new Queue<string>();

    // the static reference to the singleton instance
    public static SceneManagerScript instance { get; private set; }


    // UNITY HOOKS

    void Awake()
    {
        if (instance == null)
        {
            instance = this;
        }
        else
        {
            Destroy(gameObject);
        }
        this.InitEnvRefs();
    }

    private void Start()
    {
        this.InitWebSocketClient();
    }

    private void Update()
    {
        // process all queued server messages
        while (this.gameServerMessageQueue.Count > 0)
        {
            this.HandleServerMessage(this.gameServerMessageQueue.Dequeue());
        }
    }

    private void OnDestroy()
    {
        // close websocket connection
        this.ws.Close(CloseStatusCode.Normal);
    }

    // INTERFACE METHODS

    public void InitMainPlayer(string playerName)
    {
        // create player game object and give random start position
        float randX = Random.Range(
            Functions.GetBound(this.gameState, Vector3.down),
            Functions.GetBound(this.gameState, Vector3.up)
        );
        float randY = Random.Range(
            Functions.GetBound(this.gameState, Vector3.left),
            Functions.GetBound(this.gameState, Vector3.right)
        );
        var randStartPos = new Vector3(randX, randY, 0);
        this.mainPlayerGO = Instantiate(this.playerPrefab, randStartPos, Quaternion.identity);
        // create player model
        this.mainPlayerModel = new Player(
            id: System.Guid.NewGuid().ToString(),
            active: true,
            name: playerName,
            position: new Position(randStartPos.x, randStartPos.y),
            size: 1
        );
        var mainPlayerScript = this.mainPlayerGO.GetComponent<PlayerScript>();
        mainPlayerScript.playerModel = this.mainPlayerModel;
        mainPlayerScript.isMainPlayer = true;
        // send "player enter" message to server
        var m = new ClientMessagePlayerEnter(this.mainPlayerModel);
        this.SendWebsocketClientMessage(JsonUtility.ToJson(m));
    }

    // client-to-server sync methods

    public void SyncPlayerState(GameObject playerGO)
    {
        this.mainPlayerModel.position = new Position(
            playerGO.transform.position.x,
            playerGO.transform.position.y
        );
        var m = new ClientMessagePlayerPosition(
            this.mainPlayerModel.id,
            this.mainPlayerModel.position
        );
        this.SendWebsocketClientMessage(JsonUtility.ToJson(m));
    }

    public void SyncPlayerEatFood(Food foodModel)
    {
        var m = new ClientMessagePlayerEatFood(
            this.mainPlayerModel.id,
            foodModel.id
        );
        this.SendWebsocketClientMessage(JsonUtility.ToJson(m));
    }

    public void SyncPlayerHitMine(Mine mindModel)
    {
        var m = new ClientMessageMineDamagePlayer(
            this.mainPlayerModel.id,
            mindModel.id
        );
        this.SendWebsocketClientMessage(JsonUtility.ToJson(m));
    }

    public void SyncPlayerEatPlayer(Player otherPlayerModel)
    {
        var m = new ClientMessagePlayerEatPlayer(
            this.mainPlayerModel.id,
            otherPlayerModel.id
        );
        this.SendWebsocketClientMessage(JsonUtility.ToJson(m));
    }

    // IMPLEMENTATION METHODS

    // game server message routing and handlers

    private void HandleServerMessage(string messageJSON)
    {
        // parse message type
        string messageType = JsonUtility.FromJson<ServerMessageGeneric>(messageJSON).messageType;
        // route message to handler based on message type
        switch (messageType)
        {
            case Constants.SERVER_MESSAGE_TYPE_GAME_STATE:
                this.HandleGameStateServerMessage(messageJSON);
                break;
            case Constants.SERVER_MESSAGE_TYPE_PLAYER_ENTER:
                this.HandlePlayerEnterServerMessage(messageJSON);
                break;
            case Constants.SERVER_MESSAGE_TYPE_PLAYER_EXIT:
                this.HandlePlayerExitServerMessage(messageJSON);
                break;
            case Constants.SERVER_MESSAGE_TYPE_PLAYER_STATE_UPDATE:
                this.HandlePlayerStateUpdateServerMessage(messageJSON);
                break;
            case Constants.SERVER_MESSAGE_TYPE_FOOD_STATE_UPDATE:
                this.HandleFoodStateUpdateServerMessage(messageJSON);
                break;
            case Constants.SERVER_MESSAGE_TYPE_MINE_STATE_UPDATE:
                this.HandleMineStateUpdateServerMessage(messageJSON);
                break;
            case Constants.SERVER_MESSAGE_TYPE_SECONDS_TO_NEXT_ROUND_START:
                this.HandleSecondsToNextRoundStartServerMessage(messageJSON);
                break;
            case Constants.SERVER_MESSAGE_TYPE_SECONDS_TO_CURRENT_ROUND_END:
                this.HandleSecondsToCurrentRoundEndServerMessage(messageJSON);
                break;
            case Constants.SERVER_MESSAGE_TYPE_ROUND_RESULT:
                this.HandleRoundResultServerMessage(messageJSON);
                break;
            default:
                Debug.LogWarning("Server message not processed: " + messageJSON);
                break;
        }
    }

    private void HandleGameStateServerMessage(string messageJSON)
    {
        var gameStateMessage = JsonUtility.FromJson<ServerMessageGameState>(messageJSON);
        this.gameState = gameStateMessage.gameState;
        // sequence of biz-logic for every full game state update
        this.DeleteEnvGameObjects();
        this.InitEnvRefs();
        this.UpdateFromGameState();
        // if game not yet initialized, do so here 
        if (!this.gameStateInitialized)
        {
            this.gameStateInitialized = true;
            // since scene is initialized, make add name UI visible so that player can join the game
            this.giveNameUI.SetActive(true);
        }
    }

    private void HandlePlayerEnterServerMessage(string messageJSON)
    {
        var playerEnterMessage = JsonUtility.FromJson<ServerMessagePlayerEnter>(messageJSON);
        if (!this.PlayerIsMainPlayer(playerEnterMessage.player))
        {
            this.AddOtherPlayerFromPlayerModel(playerEnterMessage.player);
        }
    }

    private void HandlePlayerExitServerMessage(string messageJSON)
    {
        var playerExitMessage = JsonUtility.FromJson<ServerMessagePlayerExit>(messageJSON);
        if (this.playerIdToOtherPlayerGO.ContainsKey(playerExitMessage.playerId))
        {
            Object.Destroy(this.playerIdToOtherPlayerGO[playerExitMessage.playerId]);
            this.playerIdToOtherPlayerGO.Remove(playerExitMessage.playerId);
        }
    }

    private void HandlePlayerStateUpdateServerMessage(string messageJSON)
    {
        var playerUpdateMessage = JsonUtility.FromJson<ServerMessagePlayerUpdate>(messageJSON);
        Player playerModel = playerUpdateMessage.player;
        if (playerModel.id == this.mainPlayerModel.id)
        {
            this.mainPlayerGO.GetComponent<PlayerScript>()
                .UpdateFromPlayerModel(playerModel);
        }
        else if (this.playerIdToOtherPlayerGO.ContainsKey(playerModel.id))
        {
            this.playerIdToOtherPlayerGO[playerModel.id]
                .GetComponent<PlayerScript>()
                .UpdateFromPlayerModel(playerModel);
        }
        else
        {
            Debug.LogWarning("Received update message for player that doesn't exist: " + playerModel.id.ToString());
        }
    }

    private void HandleFoodStateUpdateServerMessage(string messageJSON)
    {
        var foodUpdateMessage = JsonUtility.FromJson<ServerMessageFoodUpdate>(messageJSON);
        Food foodModel = foodUpdateMessage.food;
        if(this.foodIdToFoodGO.ContainsKey(foodModel.id))
        {
            this.foodIdToFoodGO[foodModel.id].GetComponent<FoodScript>()
                .UpdateFromFoodModel(foodModel);
        }
        else
        {
            Debug.LogWarning("Received update message for food that doesn't exist: " + foodModel.id.ToString());
        }
    }

    private void HandleMineStateUpdateServerMessage(string messageJSON)
    {
        var mineUpdateMessage = JsonUtility.FromJson<ServerMessageMineUpdate>(messageJSON);
        Mine mineModel = mineUpdateMessage.mine;
        if (this.mineIdToMineGO.ContainsKey(mineModel.id))
        {
            this.mineIdToMineGO[mineModel.id].GetComponent<MineScript>()
                .UpdateFromMineModel(mineModel);
        }
        else
        {
            Debug.LogWarning("Received update message for mine that doesn't exist: " + mineModel.id.ToString());
        }
    }

    private void HandleSecondsToNextRoundStartServerMessage(string messageJSON)
    {
        // stub
    }

    private void HandleSecondsToCurrentRoundEndServerMessage(string messageJSON)
    {
        // stub
    }

    private void HandleRoundResultServerMessage(string messageJSON)
    {
        // stub
    }

    // game data management

    private void InitEnvRefs()
    {
        this.wallGOs = new List<GameObject>();
        this.playerIdToOtherPlayerGO = new Dictionary<string, GameObject>();
        this.foodIdToFoodGO = new Dictionary<string, GameObject>();
        this.mineIdToMineGO = new Dictionary<string, GameObject>();
    }

    private void DeleteEnvGameObjects()
    {
        // delete any exising walls
        foreach (GameObject wallGO in this.wallGOs)
        {
            Object.Destroy(wallGO);
        }
        // delete existing other players
        foreach (GameObject otherPlayerGO in this.playerIdToOtherPlayerGO.Values)
        {
            Object.Destroy(otherPlayerGO);
        }
        // delete existing food
        foreach (GameObject foodGO in this.foodIdToFoodGO.Values)
        {
            Object.Destroy(foodGO);
        }
        // delete existing mines
        foreach (GameObject mineGO in this.mineIdToMineGO.Values)
        {
            Object.Destroy(mineGO);
        }
    }

    private void UpdateFromGameState()
    {
        // add walls to scene
        this.CreateWalls();
        foreach (Player player in this.gameState.players)
        {
            // do main player update
            if (this.PlayerIsMainPlayer(player))
            {
                this.mainPlayerGO.GetComponent<PlayerScript>()
                    .UpdateFromPlayerModel(player, forceMainPlayerPosition: true);
            }
            // add other players to scene
            else
            {
                this.AddOtherPlayerFromPlayerModel(player);
            }
        }
        // add food to scene
        foreach (Food food in this.gameState.foods)
        {
            this.AddFoodFromFoodModel(food);
        }
        // add mines to scene
        foreach (Mine mine in this.gameState.mines)
        {
            this.AddMineFromMineModel(mine);
        }
    }

    private void CreateWalls()
    {
        var wallTop = Instantiate(
            this.wallPrefab,
            new Vector3(0, Functions.GetBound(this.gameState, Vector3.up) + 1, 0),
            Quaternion.identity
        );
        wallTop.transform.localScale = new Vector3(this.gameState.mapWidth + 3, 1, 0);
        this.wallGOs.Add(wallTop);
        var wallBottom = Instantiate(
            this.wallPrefab,
            new Vector3(0, Functions.GetBound(this.gameState, Vector3.down) - 1, 0),
            Quaternion.identity
        );
        wallBottom.transform.localScale = new Vector3(this.gameState.mapWidth + 3, 1, 0);
        this.wallGOs.Add(wallBottom);
        var wallLeft = Instantiate(
            this.wallPrefab,
            new Vector3(Functions.GetBound(this.gameState, Vector3.left) - 1, 0, 0),
            Quaternion.identity
        );
        wallLeft.transform.localScale = new Vector3(1, this.gameState.mapHeight + 3, 0);
        this.wallGOs.Add(wallLeft);
        var wallRight = Instantiate(
            this.wallPrefab,
            new Vector3(Functions.GetBound(this.gameState, Vector3.right) + 1, 0, 0),
            Quaternion.identity
        );
        wallRight.transform.localScale = new Vector3(1, this.gameState.mapHeight + 3, 0);
        this.wallGOs.Add(wallRight);
    }

    private void AddOtherPlayerFromPlayerModel(Player otherPlayerModel)
    {
        // player is not currently tracked
        if (!this.playerIdToOtherPlayerGO.ContainsKey(otherPlayerModel.id))
        {
            //Debug.Log("adding other player: " + otherPlayerModel.id.ToString());
            var otherPlayerPosition = new Vector3(
                otherPlayerModel.position.x,
                otherPlayerModel.position.y,
                0
            );
            GameObject otherPlayerGO = Instantiate(
                this.playerPrefab,
                otherPlayerPosition,
                Quaternion.identity
            );
            var otherPlayerScript = otherPlayerGO.GetComponent<PlayerScript>();
            otherPlayerScript.playerModel = otherPlayerModel;
            otherPlayerScript.isMainPlayer = false;
            this.playerIdToOtherPlayerGO.Add(otherPlayerModel.id, otherPlayerGO);
        }
    }

    private void AddFoodFromFoodModel(Food food)
    {
        var foodPosition = new Vector3(
            food.position.x,
            food.position.y,
            0
        );
        GameObject foodGO = Instantiate(
            this.foodPrefab,
            foodPosition,
            Quaternion.identity
        );
        this.foodIdToFoodGO.Add(food.id, foodGO);
        foodGO.GetComponent<FoodScript>().foodModel = food;
    }

    private void AddMineFromMineModel(Mine mine)
    {
        var minePosition = new Vector3(
            mine.position.x,
            mine.position.y,
            0
        );
        GameObject mineGO = Instantiate(
            this.minePrefab,
            minePosition,
            Quaternion.identity
        );
        this.mineIdToMineGO.Add(mine.id, mineGO);
        mineGO.GetComponent<MineScript>().mineModel = mine;
    }

    // websocket helpers

    private void InitWebSocketClient()
    {
        // create websocket connection
        this.ws = new WebSocket(this.gameServerUrl);
        this.ws.Connect();
        // add message handler callback
        this.ws.OnMessage += this.QueueServerMessage;
    }

    private void QueueServerMessage(object sender, MessageEventArgs e)
    {
        //Debug.Log("Server message received: " + e.Data);
        this.gameServerMessageQueue.Enqueue(e.Data);
    }

    private void SendWebsocketClientMessage(string messageJson)
    {
        //Debug.Log("Client message sent: " + messageJson);
        this.ws.Send(messageJson);
    }

    // misc helpers

    private bool PlayerIsMainPlayer(Player p)
    {
        return this.mainPlayerModel != null && p.id == this.mainPlayerModel.id;
    }

}
