using UnityEngine;
using System.Collections;
using System.Linq;

public class WebCamScreenTexture : MonoBehaviour
{

	public string cameraName = "FaceRig Virtual Camera";

	public int Width = 1280;
	public int Height = 720;
	public int FPS = 60;
	public Material material = null;

	private WebCamTexture webcamTexture;

	// Use this for initialization
	void Start ()
	{
		WebCamTexture.devices.ToList ().ForEach (v => print (v.name));

		webcamTexture = FindWebCameraByName (cameraName);
		if (webcamTexture == null) {
			Debug.LogWarning ("Could not find " + cameraName);
			return;
		}

		if (material == null) {
			material = gameObject.GetComponent<Renderer> ().material;
		}

		material.mainTexture = webcamTexture;
		webcamTexture.Play();
	}

	private WebCamTexture FindWebCameraByName(string cameraName) {
		return WebCamTexture.devices
			.Where (v => v.name == cameraName)
			.Select(v => new WebCamTexture (v.name, Width, Height, FPS))
			.FirstOrDefault ();
	}
}
