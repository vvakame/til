using UnityEngine;
using UnityEngine.UI;
using System.Collections;

public class CoinRemoverController : MonoBehaviour {
	public int scoreScale;
	public Text scoreText;

	private int count;

	void OnTriggerEnter(Collider other) {
		if (other.CompareTag("Coin") ){
			Destroy(other.gameObject );
			count += scoreScale;
			SetCountText ();
		}
	}

	void SetCountText()
	{
		if (scoreText == null) {
			return;
		}

		scoreText.text = "Score: " + count.ToString ();
	}
}
