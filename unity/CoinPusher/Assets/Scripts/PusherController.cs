using UnityEngine;
using System.Collections;

public class PusherController : MonoBehaviour
{

	public float cycleSec = 1; // 何秒かけて台が一周期動くか
	public float maxDelta = 0.5f;
	private Vector3 defaultPos;
	private int currentCycle = 0;

	void Start ()
	{
		defaultPos = transform.localPosition;
	}

	void Update ()
	{
		// http://yaseino.hatenablog.com/entry/2016/02/23/234652 を参考に
		float translation = Mathf.Cos (Time.time * 6 / cycleSec) * (maxDelta/2) - (maxDelta/2); // ゲームスタート時点が一番板が引かれてる状態にしたい
		transform.localPosition = defaultPos;
		transform.Translate (0, 0, translation);
	}
}
