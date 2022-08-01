using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class Mapblock : MonoBehaviour
{
    public bool Block
    {
        get { return b_block_; }
        set { b_block_ = value; }
    }

    public int X
    {
        get { return pos_x_; }
        set { pos_x_ = value; }
    }
    public int Y
    {
        get { return pos_y_; }
        set { pos_y_ = value; }
    }

    private bool b_block_ = false;
    private int pos_x_ = 0;
    private int pos_y_ = 0;
    // Start is called before the first frame update
    void Start()
    {
        
    }

    // Update is called once per frame
    void Update()
    {
        
    }
}
