using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class MineScript : MonoBehaviour
{

    public Mine mineModel;

    // UNITY HOOKS

    void Start() { }

    void Update() { }

    // INTERFACE METHODS

    public void UpdateFromMineModel(Mine m)
    {
        // stub
        Debug.Log("UpdateFromMineModel...");
    }

}
