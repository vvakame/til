using UnityEngine;
using System.Collections;

public class CoinSpawnerController : MonoBehaviour {

	public GameObject coinPrefab;

	// Update is called once per frame
	void Update () {
		if (Input.GetMouseButtonDown (0) || Input.GetButtonDown ("Jump")) {
			Instantiate (coinPrefab, transform.position, transform.rotation);
		}
	}
}
