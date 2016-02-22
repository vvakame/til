using UnityEngine;
using System.Collections;

public class GamePadController : MonoBehaviour
{

	public float upDelta = 1;

	private Vector3 defaultPos;

	// Use this for initialization
	void Start ()
	{
		defaultPos = transform.position;
	}
	
	// Update is called once per frame
	void Update ()
	{
		var vertical = Input.GetAxis ("Vertical");
		if (vertical < 0) {
			DownGamePad ();
		} else if (vertical > 0) {
			UpGamePad ();
		}

		var horizontal = Input.GetAxis ("Horizontal");
		if (horizontal < 0) {
			LeftGamePad ();
		} else if (horizontal > 0) {
			RightGamePad ();
		}

		if (vertical == 0 && horizontal == 0) {
			ResetGamePadPosition ();
		}

		Debug.Log ("Horizontal: " + Input.GetAxis ("Horizontal"));
		Debug.Log ("Vertical: " + Input.GetAxis ("Vertical"));
		Debug.Log ("Fire1: " + Input.GetButton ("Fire1"));
		Debug.Log ("Fire2: " + Input.GetButton ("Fire2"));
		Debug.Log ("Fire3: " + Input.GetButton ("Fire3"));
	}

	private void ResetGamePadPosition ()
	{
		transform.position = defaultPos;
	}

	private void UpGamePad ()
	{
		ResetGamePadPosition ();
		transform.Translate (Vector3.up * upDelta, Space.World);
	}

	private void DownGamePad ()
	{
		ResetGamePadPosition ();
		transform.Translate (Vector3.down * upDelta, Space.World);
	}

	private void LeftGamePad ()
	{
		ResetGamePadPosition ();
		transform.Translate (Vector3.left * upDelta, Space.World);
	}

	private void RightGamePad ()
	{
		ResetGamePadPosition ();
		transform.Translate (Vector3.right * upDelta, Space.World);
	}
}
