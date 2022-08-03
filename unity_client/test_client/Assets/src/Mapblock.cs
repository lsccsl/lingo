using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

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
    private Color org_clr_;
    // Start is called before the first frame update
    void Start()
    {
        org_clr_ = this.gameObject.GetComponent<MeshRenderer>().material.color;
    }

    // Update is called once per frame
    void Update()
    {
        
    }

    public void set_clr(Color clr)
    {
        this.gameObject.GetComponent<MeshRenderer>().material.color = clr;
    }

    public void set_path_flag()
    {
        this.gameObject.GetComponent<MeshRenderer>().material.color = new Color(1, 0, 0);
    }
    public void set_start_end_flag()
    {
        this.gameObject.GetComponent<MeshRenderer>().material.color = new Color(1, 0, 1);
    }
    public void reset_flag()
    {
        this.gameObject.GetComponent<MeshRenderer>().material.color = org_clr_;
    }
}
