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

    public void UpdateFromFoodModel(Food f)
    {
        // stub
        Debug.Log("UpdateFromFoodModel...");
    }


}
