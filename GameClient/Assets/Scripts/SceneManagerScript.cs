using System.Collections.Generic;
using UnityEngine;
using WebSocketSharp;

public class SceneManagerScript : MonoBehaviour
{

    public GameObject playerPrefab;
    public GameObject wallPrefab;

    public GameObject giveNameUI;

    private WebSocket ws;
    private string gameServerUrl = "ws://localhost:5000"; // local game server
    //private string gameServerUrl = "ws://golang-multiplayer-server-demo.herokuapp.com/:80"; // heroku game server

    private bool gameStateInitialized = false;
    private GameState gameState;
    private Player mainPlayerModel;
    private GameObject mainPlayerGO;

    private IDictionary<string, GameObject> playerIdToOtherPlayerGO =
            new Dictionary<string, GameObject>();

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
        // create player game object
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
        var mainPlayerScript = this.mainPlayerGO.GetComponent<PlayerScript>();
        mainPlayerScript.sceneManager = this;
        mainPlayerScript.isMainPlayer = true;
        // create player model
        this.mainPlayerModel = new Player(
            id: System.Guid.NewGuid().ToString(),
            active: true,
            name: playerName,
            position: new Position(randStartPos.x, randStartPos.y),
            size: 1
        );
        // send "player enter" message to server
        var playerEnterMessage = new ClientMessagePlayerEnter(this.mainPlayerModel);
        this.SendWebsocketClientMessage(JsonUtility.ToJson(playerEnterMessage));
    }

    public void SyncPlayerState(GameObject playerGO)
    {
        // send "player position" message to server
        this.mainPlayerModel.position = new Position(
            playerGO.transform.position.x,
            playerGO.transform.position.y
        );
        var playerPositionMessage = new ClientMessagePlayerPosition(
            this.mainPlayerModel.id,
            this.mainPlayerModel.position
        );
        this.SendWebsocketClientMessage(JsonUtility.ToJson(playerPositionMessage));
    }

    // IMPLEMENTATION METHODS

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
        if (!this.gameStateInitialized)
        {
            this.gameState = gameStateMessage.gameState;
            this.InitGameState();
            this.gameStateInitialized = true;
        }
        else {
            this.ProcessGameStateUpdate(gameStateMessage.gameState);
        }
    }

    private void HandlePlayerEnterServerMessage(string messageJSON)
    {
        var playerEnterMessage = JsonUtility.FromJson<ServerMessagePlayerEnter>(messageJSON);
        this.AddOtherPlayerFromPlayerModel(playerEnterMessage.player);
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
        if (this.playerIdToOtherPlayerGO.ContainsKey(playerModel.id))
        {
            var newPosition = new Vector3(
                playerModel.position.x,
                playerModel.position.y,
                0
            );
            this.playerIdToOtherPlayerGO[playerModel.id].transform.position = newPosition;
        }
    }

    private void HandleFoodStateUpdateServerMessage(string messageJSON)
    {
        // stub
    }

    private void HandleMineStateUpdateServerMessage(string messageJSON)
    {
        // stub
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

    private void InitGameState() {
        // make add name UI visible
        this.giveNameUI.SetActive(true);
        // add walls to scene
        this.CreateWalls();
        // add other players to scene
        foreach (Player player in this.gameState.players)
        {
            this.AddOtherPlayerFromPlayerModel(player);
        }
    }

    private void CreateWalls()
    {
        var wallTop = Instantiate(
            this.wallPrefab,
            new Vector3(0, Functions.GetBound(this.gameState, Vector3.up), 0),
            Quaternion.identity
        );
        wallTop.transform.localScale = new Vector3(this.gameState.mapWidth, 1, 0);
        var wallBottom = Instantiate(
            this.wallPrefab,
            new Vector3(0, Functions.GetBound(this.gameState, Vector3.down), 0),
            Quaternion.identity
        );
        wallBottom.transform.localScale = new Vector3(this.gameState.mapWidth, 1, 0);
        var wallLeft = Instantiate(
            this.wallPrefab,
            new Vector3(Functions.GetBound(this.gameState, Vector3.left), 0, 0),
            Quaternion.identity
        );
        wallLeft.transform.localScale = new Vector3(1, this.gameState.mapHeight, 0);
        var wallRight = Instantiate(
            this.wallPrefab,
            new Vector3(Functions.GetBound(this.gameState, Vector3.right), 0, 0),
            Quaternion.identity
        );
        wallRight.transform.localScale = new Vector3(1, this.gameState.mapHeight, 0);
    }

    private void ProcessGameStateUpdate(GameState gameState) {
        // stub
    }

    private void AddOtherPlayerFromPlayerModel(Player otherPlayerModel)
    {
        // player is not main player and player is not currently tracked
        bool isMainPlayer = (this.mainPlayerModel != null && otherPlayerModel.id == this.mainPlayerModel.id);
        if (!isMainPlayer && !this.playerIdToOtherPlayerGO.ContainsKey(otherPlayerModel.id))
        {
            Debug.Log("adding other player: " + otherPlayerModel.id.ToString());
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
            otherPlayerScript.sceneManager = this;
            otherPlayerScript.isMainPlayer = false;
            this.playerIdToOtherPlayerGO.Add(otherPlayerModel.id, otherPlayerGO);
        }
    }

    private void SendWebsocketClientMessage(string messageJson)
    {
        //Debug.Log("Client message sent: " + messageJson);
        this.ws.Send(messageJson);
    }

}
