using System.Collections.Generic;
using UnityEngine;

public class PlayerScript : MonoBehaviour
{

    public Player playerModel;
    public bool isMainPlayer;

    // movement
    private float moveSpeed = 5f;
    // TESTING:
    // use for testing if you want to connect multiple players to the server and
    // see them moving
    private bool autopilotOn = false;
    // autopilot movement for testing
    private List<Vector3> moveDirections = new List<Vector3> {
        Vector3.up,
        Vector3.right,
        Vector3.down,
        Vector3.left
    };
    private int currMoveDirIndex = 0;

    // camera
    private GameObject camGO;
    private Camera cam;
    private float camZoomSpeed = 10f;
    private float minCamZoom = 2f;
    private float maxCamZoom = 30f;

    // UNITY HOOKS

    void Start()
    {
        if (this.isMainPlayer)
        {
            this.camGO = SceneManagerScript.instance.mainCameraGO;
            this.cam = this.camGO.GetComponent<Camera>();
            this.MoveCameraToPlayer();
            // autopilot movement for testing
            InvokeRepeating("SetNextMoveDirectionIndex", 0f, 1f);
        }
        else
        {
            this.gameObject.GetComponent<SpriteRenderer>().color = Color.red;
        }
    }

    void Update()
    {
        if (this.isMainPlayer)
        {
            this.HandleMovement();
            this.HandleCameraZoom();
        }
    }

    // IMPLEMENTATION METHODS

    private void HandleMovement()
    {
        var targetPos = this.transform.position;
        if (Input.anyKey)
        {
            // left
            if (Input.GetKey(KeyCode.A))
            {
                targetPos += Vector3.left;
            }
            // right
            if (Input.GetKey(KeyCode.D))
            {
                targetPos += Vector3.right;
            }
            // up
            if (Input.GetKey(KeyCode.W))
            {
                targetPos += Vector3.up;
            }
            // down
            if (Input.GetKey(KeyCode.S))
            {
                targetPos += Vector3.down;
            }
        }
        else if(this.autopilotOn) {
            targetPos += this.moveDirections[this.currMoveDirIndex];
        }
        if (targetPos != this.transform.position)
        {
            this.transform.position = Vector3.MoveTowards(
                this.transform.position,
                targetPos,
                Time.deltaTime * this.moveSpeed
            );
            this.MoveCameraToPlayer();
            SceneManagerScript.instance.SyncPlayerState(this.gameObject);
        }
    }

    private void HandleCameraZoom()
    {
        int direction = 0;
        if (Input.GetKey(KeyCode.UpArrow))
        {
            direction = -1;
        }
        else if (Input.GetKey(KeyCode.DownArrow))
        {
            direction = 1;
        }
        if (direction != 0) {
            float zoomChange = direction * this.camZoomSpeed * Time.deltaTime;
            float newCameraSize = this.cam.orthographicSize + zoomChange;
            if (newCameraSize > this.minCamZoom && newCameraSize < this.maxCamZoom) {
                this.cam.orthographicSize = newCameraSize;
            }
        }
    }

    private void MoveCameraToPlayer() {
        this.camGO.transform.position = new Vector3(
            this.transform.position.x,
            this.transform.position.y,
            this.camGO.transform.position.z
        );
    }

    // autopilot movement for testing
    private void SetNextMoveDirectionIndex() {
        if (this.currMoveDirIndex == 3)
        {
            this.currMoveDirIndex = 0;
        }
        else
        {
            this.currMoveDirIndex += 1;
        }
    }

}
