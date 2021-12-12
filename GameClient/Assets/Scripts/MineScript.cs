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

    public void UpdateFromMineModel(Mine mModel)
    {
        this.mineModel = mModel;
        this.transform.position = new Vector3(mineModel.position.x, mModel.position.y, 0);
        this.transform.localScale = new Vector3(mModel.size, mModel.size, 1);
        this.gameObject.SetActive(mModel.active);
    }

}
