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
	public float maxRotateY = 40;

	private Vector3 defaultRotate;

	void Start ()
	{
		defaultRotate = transform.eulerAngles;
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
			rotateArrow (-1);
		} else if (0 < moveHorizontal) {
			rotateArrow (1);
		}
	}

	void rotateArrow (float direction)
	{
		// 左 0 より小さくなると 0-360 の範囲に補正される
		// 例: 0から+40と-40の範囲について。-40 は 320 に補正される。 つまり、0-40 と 320-360 の範囲にあればよい。
		//     350から+40と-40の範囲について。 370は30に補正される。
		var min1 = defaultRotate.y - maxRotateY;
		var max1 = defaultRotate.y + maxRotateY;
		var min2 = min1;
		var max2 = max1;
		if (min1 < 0) {
			min2 += 360;
			max2 += 360;
		} else if (360 < max1) {
			min1 -= 360;
			max1 -= 360; 
		}

		var newY = transform.eulerAngles.y + moveSpeed * direction;
		while (newY < 0) {
			newY += 360;
		}
		while (360 < newY) {
			newY -= 360;
		}
		if (min1 <= newY && newY <= max1) {
			// OK
		} else if (min2 <= newY && newY <= max2) {
			// OK
		} else if (newY < min1 && newY < max1) {
			newY = min1;
		} else if (min2 < newY && max2 < newY) {
			newY = max2;
		} else if (newY - max1 < min2 - newY) {
			newY = max1;
		} else {
			newY = min2;
		}

		transform.rotation = Quaternion.Euler (0, newY, 0);
	}
}
