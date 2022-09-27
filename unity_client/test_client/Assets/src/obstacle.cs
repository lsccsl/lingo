using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class obstacle : MonoBehaviour
{
    public uint obstacle_id = 0;

    // Start is called before the first frame update
    void Start()
    {
        
    }

    // Update is called once per frame
    void Update()
    {
        
    }

    public void set_scale(Vector3 v_scale)
    {
        this.gameObject.transform.localScale = v_scale;
    }
}
