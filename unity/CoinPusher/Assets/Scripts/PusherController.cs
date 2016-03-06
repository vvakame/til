using UnityEngine;
using System.Collections;

public class PusherController : MonoBehaviour
{
	// 何秒かけて台が一周期動くか
	public float cycleSec = 1;
	// 最大台幅の何割分動くか
	public float maxDelta = 0.5f;

	private Vector3 defaultPos;
	private int currentCycle = 0;

	void Start ()
	{
		defaultPos = transform.position;
	}

	void Update ()
	{
		// http://yaseino.hatenablog.com/entry/2016/02/23/234652 を参考に
		Vector3 offset = new Vector3 (0, 0, -1 * (Mathf.Cos (Time.time * 6 / cycleSec) * (maxDelta / 2) - (maxDelta / 2)));
		GetComponent<Rigidbody> ().MovePosition (defaultPos + offset);
	}
}
