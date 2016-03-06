using UnityEngine;
using System.Collections;

public class CoinSpawnerController : MonoBehaviour
{

	public GameObject coinPrefab;
	public int startSpawn = 100;
	public float startSpawnScale = 0.8f;

	void Start ()
	{
		for (var i = 0; i < startSpawn; i++) {
			Instantiate (coinPrefab, transform.position + new Vector3 ((Random.value - 0.5f) * startSpawnScale, 0, ((Random.value - 1.0f) * 0.5f) * startSpawnScale), transform.rotation);
		}
	}

	void Update ()
	{
		if (Input.GetMouseButtonDown (0) || Input.GetButtonDown ("Jump")) {
			Instantiate (coinPrefab, transform.position, transform.rotation);
			Score.Unlock ();
		}
	}
}
