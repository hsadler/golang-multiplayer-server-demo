using System.Collections;
using System.Collections.Generic;
using TMPro;
using UnityEngine;

public class GiveNameUIScript : MonoBehaviour
{

    public TMP_InputField tmpInput;

    void Start()
    {
        tmpInput.onSubmit.AddListener(this.AddPlayer);
    }

    private void AddPlayer(string playerName) {
        SceneManagerScript.instance.InitMainPlayer(playerName);
        this.gameObject.SetActive(false);
    }

}
