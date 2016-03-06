using UnityEngine;
using System.Collections;

public class CoinSpawnerController : MonoBehaviour
{

	public GameObject coinPrefab;

	public int startSpawn = 100;
	public float startSpawnScale = 0.8f;

	public int shootPiece = 1;
	public float shootX = 0f;
	public float shootY = 100f;
	public float shootZ = 300f;
	public float shootCourseRandomness = 10f;

	void Start ()
	{
		for (var i = 0; i < startSpawn; i++) {
			Instantiate (coinPrefab, transform.position + new Vector3 ((Random.value - 0.5f) * startSpawnScale, 0, ((Random.value - 1.0f) * 0.5f) * startSpawnScale), transform.rotation);
		}
	}

	void Update ()
	{
		if (Input.GetMouseButtonDown (0) || Input.GetButtonDown ("Jump")) {
			for (var i = 0; i < shootPiece; i++) {
				var coin = Instantiate (coinPrefab, transform.position, transform.rotation) as GameObject;
				var coinRigid = coin.GetComponent<Rigidbody> ();
				var shootForce = new Vector3 (shootX, shootY, shootZ);
				coinRigid.AddForce (shootForce);
				coin.transform.rotation = Random.rotation;
			}
			Score.Unlock ();
		}
	}
}
