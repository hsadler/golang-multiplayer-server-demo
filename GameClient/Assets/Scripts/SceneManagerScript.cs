using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using TMPro;
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
    public GameObject respawnCountdownUI;
    public TMP_Text respawnCountdownText;
    public TMP_Text roundTimerText;

    private WebSocket ws;
    // local game server
    private string gameServerUrl = "ws://localhost:5000";
    // heroku game server
    //private string gameServerUrl = "ws://golang-multiplayer-server-demo.herokuapp.com/:80";

    // game state
    public GameState gameState;
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
        Vector3 randStartPos = Functions.GetRandomGamePosition(this.gameState);
        this.mainPlayerGO = Instantiate(this.playerPrefab, randStartPos, Quaternion.identity);
        // create player model
        this.mainPlayerModel = new Player(
            id: System.Guid.NewGuid().ToString(),
            active: true,
            name: playerName,
            position: new Position(randStartPos.x, randStartPos.y),
            size: 1,
            timeUntilRespawn: 0
        );
        var mainPlayerScript = this.mainPlayerGO.GetComponent<PlayerScript>();
        mainPlayerScript.playerModel = this.mainPlayerModel;
        mainPlayerScript.isMainPlayer = true;
        // send "player enter" message to server
        var m = new ClientMessagePlayerEnter(this.mainPlayerModel);
        this.SendWebsocketClientMessage(JsonUtility.ToJson(m));
    }

    // client-to-server sync methods

    public void SyncPlayerPosition(Player playerModel)
    {
        var m = new ClientMessagePlayerPosition(
            playerModel.id,
            playerModel.position
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
        var m = JsonUtility.FromJson<ServerMessageGameState>(messageJSON);
        this.gameState = m.gameState;
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
        var m = JsonUtility.FromJson<ServerMessagePlayerEnter>(messageJSON);
        if (!this.PlayerIsMainPlayer(m.player))
        {
            this.AddOtherPlayerFromPlayerModel(m.player);
        }
    }

    private void HandlePlayerExitServerMessage(string messageJSON)
    {
        var m = JsonUtility.FromJson<ServerMessagePlayerExit>(messageJSON);
        if (this.playerIdToOtherPlayerGO.ContainsKey(m.playerId))
        {
            Object.Destroy(this.playerIdToOtherPlayerGO[m.playerId]);
            this.playerIdToOtherPlayerGO.Remove(m.playerId);
        }
    }

    private void HandlePlayerStateUpdateServerMessage(string messageJSON)
    {
        var m = JsonUtility.FromJson<ServerMessagePlayerUpdate>(messageJSON);
        Player playerModel = m.player;
        // main player update
        if (this.mainPlayerModel != null && playerModel.id == this.mainPlayerModel.id)
        {
            if (!playerModel.active)
            {
                this.respawnCountdownUI.SetActive(true);
                if (playerModel.timeUntilRespawn > 5) {
                    this.respawnCountdownText.text = "";
                } else
                {
                    this.respawnCountdownText.text = playerModel.timeUntilRespawn.ToString();
                }
            }
            else {
                this.respawnCountdownUI.SetActive(false);
            }
            //Debug.Log("updating main-player: " + playerModel.name);
            this.mainPlayerGO.GetComponent<PlayerScript>().UpdateFromPlayerModel(playerModel);
        }
        // other players update
        else if (this.playerIdToOtherPlayerGO.ContainsKey(playerModel.id))
        {
            //Debug.Log("updating other-player: " + playerModel.name);
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
        var m = JsonUtility.FromJson<ServerMessageFoodUpdate>(messageJSON);
        Food foodModel = m.food;
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
        var m = JsonUtility.FromJson<ServerMessageMineUpdate>(messageJSON);
        Mine mineModel = m.mine;
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
        var m = JsonUtility.FromJson<ServerMessageSecondsToNextRoundStart>(messageJSON);
        this.roundTimerText.text = m.seconds.ToString();
    }

    private void HandleSecondsToCurrentRoundEndServerMessage(string messageJSON)
    {
        // stub
        var m = JsonUtility.FromJson<ServerMessageSecondsToCurrentRoundEnd>(messageJSON);
        this.roundTimerText.text = m.seconds.ToString();
    }

    private void HandleRoundResultServerMessage(string messageJSON)
    {
        // stub
        var m = JsonUtility.FromJson<ServerMessageRoundResult>(messageJSON);
        Debug.Log("HandleRoundResultServerMessage json: " + messageJSON);
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
                    .UpdateFromPlayerModel(player);
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
            GameObject otherPlayerGO = Instantiate(
                this.playerPrefab,
                Vector3.zero,
                Quaternion.identity
            );
            var otherPlayerScript = otherPlayerGO.GetComponent<PlayerScript>();
            otherPlayerScript.UpdateFromPlayerModel(otherPlayerModel);
            otherPlayerScript.isMainPlayer = false;
            this.playerIdToOtherPlayerGO.Add(otherPlayerModel.id, otherPlayerGO);
        }
    }

    private void AddFoodFromFoodModel(Food food)
    {
        GameObject foodGO = Instantiate(
            this.foodPrefab,
            Vector3.zero,
            Quaternion.identity
        );
        foodGO.GetComponent<FoodScript>().UpdateFromFoodModel(food);
        this.foodIdToFoodGO.Add(food.id, foodGO);
    }

    private void AddMineFromMineModel(Mine mine)
    {
        GameObject mineGO = Instantiate(
            this.minePrefab,
            Vector3.zero,
            Quaternion.identity
        );
        mineGO.GetComponent<MineScript>().UpdateFromMineModel(mine);
        this.mineIdToMineGO.Add(mine.id, mineGO);
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
