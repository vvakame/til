using UnityEngine;
using System.Collections;

public class MarkerController : MonoBehaviour {

	// Use this for initialization
	void Start () {
        Debug.Log("MarkerController#Start");
	}
    
    void OnTriggerEnter(Collider other)
    {
        Debug.Log(other.name);
        gameObject.GetComponent<Renderer>().material.color = Color.red;
    }

    void OnTriggerExit(Collider other)
    {
        Debug.Log(other.name);
        gameObject.GetComponent<Renderer>().material.color = Color.white;
    }
}
