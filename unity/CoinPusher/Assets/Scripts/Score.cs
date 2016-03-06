using UnityEngine;
using System.Collections;

public class Score : MonoBehaviour
{

	private	static   int score;
	public static  bool scoreLock = true;

	void Start ()
	{
		score = 0;
	}

	public static void Unlock ()
	{
		scoreLock = false;
	}

	public	static void AddScore (int v)
	{
		if (scoreLock) {
			return;
		}
		score += v;
	}

	public	static int GetScore ()
	{
		return score;
	}
}
