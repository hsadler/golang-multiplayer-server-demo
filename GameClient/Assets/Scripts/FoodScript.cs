using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class FoodScript : MonoBehaviour
{

    public Food foodModel;

    // UNITY HOOKS

    void Start() { }

    void Update() { }

    // INTERFACE METHODS

    public void UpdateFromFoodModel(Food fModel)
    {
        this.transform.position = new Vector3(fModel.position.x, fModel.position.y, 0);
        this.transform.localScale = new Vector3(fModel.size, fModel.size, 1);
        this.gameObject.SetActive(fModel.active);
    }


}
