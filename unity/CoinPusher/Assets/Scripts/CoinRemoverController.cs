using UnityEngine;
using UnityEngine.UI;
using System.Collections;

public class CoinRemoverController : MonoBehaviour
{
	public int scoreScale;
	public Text scoreText;

	void OnTriggerEnter (Collider other)
	{
		if (other.CompareTag ("Coin")) {
			Destroy (other.gameObject);
			Score.AddScore (scoreScale);
			SetCountText ();
		}
	}

	void SetCountText ()
	{
		if (scoreText == null) {
			return;
		}

		scoreText.text = "Score: " + Score.GetScore ();
	}
}
