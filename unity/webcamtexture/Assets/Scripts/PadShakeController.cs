using UnityEngine;
using System.Collections;

public class PadShakeController : MonoBehaviour
{

	public float ButtonDelta = 1;
	public float RotateDelta = 10;

	private Vector3 defaultPos;
	private Vector3 defaultAngle;

	// Use this for initialization
	void Start ()
	{
		defaultPos = transform.position;
		defaultAngle = transform.eulerAngles;
	}
	
	// Update is called once per frame
	void Update ()
	{

		ResetGamePadPosition ();

		var vertical = Input.GetAxis ("Vertical");
		XAxis (vertical);
		ButtonX (vertical);

		var horizontal = Input.GetAxis ("Horizontal");
		YAxis (horizontal);
		ButtonY (horizontal);

		if (Input.GetButton ("Fire1") || Input.GetButton ("Fire2") || Input.GetButton ("Fire3")) {
			ButtonX (-1);
		}
	}

	private void ResetGamePadPosition ()
	{
		transform.position = defaultPos;
		transform.eulerAngles = defaultAngle;
	}

	private void ButtonX (float delta)
	{
		transform.Translate (new Vector3 (0, delta, 0) * ButtonDelta, Space.World);
	}

	private void ButtonY (float delta)
	{
		transform.Translate (new Vector3 (delta, 0, 0) * ButtonDelta, Space.World);
	}

	private void XAxis (float delta)
	{
		transform.Rotate (new Vector3 (delta, 0, 0) * RotateDelta, Space.Self);
	}

	private void YAxis (float delta)
	{
		transform.Rotate (new Vector3 (0, delta, 0) * RotateDelta, Space.Self);
	}

	private void OnGUI ()
	{
		if (Debug.isDebugBuild) {
			GUI.Label (
				new Rect (0f, 0f, Screen.width, Screen.height), string.Format ("post   = x:{0}, y:{1}, z:{2}\nangles = x:{3}, y:{4}, z:{5}\nv: {6}\nh: {7}", transform.position.x, transform.position.y, transform.position.z, transform.eulerAngles.x, transform.eulerAngles.y, transform.eulerAngles.z, Input.GetAxis ("Horizontal"), Input.GetAxis ("Vertical")));
		}
	}
}
