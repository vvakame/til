using UnityEngine;
using System.Collections;

public class CoinShooterController : MonoBehaviour
{

	public GameObject coinPrefab;

	public int shootPiece = 1;
	public float shootX = 0f;
	public float shootY = 100f;
	public float shootZ = 300f;
	public float shootCourseRandomness = 10f;

	public float moveSpeed = 0.8f;

	private Quaternion defaultRotate;

	void Start ()
	{
		defaultRotate = transform.rotation;
	}

	void Update ()
	{
		if (Input.GetMouseButtonDown (0) || Input.GetButtonDown ("Jump")) {
			for (var i = 0; i < shootPiece; i++) {
				var coin = Instantiate (coinPrefab, transform.position, transform.rotation) as GameObject;
				var coinRigid = coin.GetComponent<Rigidbody> ();
				var shootForce = new Vector3 (shootX, shootY, shootZ);
				coinRigid.AddRelativeForce (shootForce);
				coin.transform.rotation = Random.rotation;
			}
			Score.Unlock ();
		}

		float moveHorizontal = Input.GetAxis ("Horizontal");
		if (moveHorizontal < 0) {
			transform.Rotate (Vector3.down * moveSpeed);
		} else if (0 < moveHorizontal) {
			transform.Rotate (Vector3.up * moveSpeed);
		}
	}
}
