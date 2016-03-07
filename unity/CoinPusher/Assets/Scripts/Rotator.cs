using UnityEngine;
using System.Collections;

public class Rotator : MonoBehaviour
{

	public int xSpeed = 30;
	public int ySpeed = 30;
	public int zSpeed = 30;

	void Update ()
	{
		transform.Rotate (new Vector3 (xSpeed, ySpeed, zSpeed) * Time.deltaTime);
	}
}
