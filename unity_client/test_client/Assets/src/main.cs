using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using Google.Protobuf;

public class main : MonoBehaviour
{
    public GameObject cube_block_;
    public GameObject cube_no_block_;

    public GameObject[] cube_all_;

    MapMgr map_mgr_;

    TestClient client_;
    // Start is called before the first frame update
    void Start()
    {
        map_mgr_ = new MapMgr();
        Debug.Log("hello");
        Msgpacket.MSG_TEST msg = new Msgpacket.MSG_TEST();
        msg.Id = 123;
        Debug.Log("msg:" + msg.ToString());
        byte[] datas = msg.ToByteArray();
        Debug.Log(datas);

        IMessage msgParse = new Msgpacket.MSG_TEST();
        Msgpacket.MSG_TEST msgNew = (Msgpacket.MSG_TEST)msgParse.Descriptor.Parser.ParseFrom(datas);
        Debug.Log("parse:" + msgNew.ToString());

        client_ = new TestClient();
        client_.connect("192.168.2.129", 2003);
    }

    // Update is called once per frame
    void Update()
    {
        if (!client_.GetRecvQue().IsEmpty())
        {
            InterMsg msg = client_.GetRecvQue().Dequeue();
            Debug.Log("msgtype:" + msg.msgtype + " msg:" + msg.msg.ToString());

            switch (msg.msgtype)
            {
                case Msgpacket.MSG_TYPE.MsgLoginRes:
                    Debug.Log("Msgpacket.MSG_TYPE.MsgLoginRes");
                    break;
                case Msgpacket.MSG_TYPE.MsgGetMapRes:
                    Debug.Log("Msgpacket.MSG_TYPE.MsgGetMapRes");
                    process_map_load_msg(msg.msg);
                    break;
            }
        }

        if(Input.GetMouseButtonDown(0))
            check_screen_ray_hit();
    }

    void check_screen_ray_hit()
    {
        Ray screen_ray = Camera.main.ScreenPointToRay(Input.mousePosition);
        RaycastHit rh;

        //LayerMask lay_terrain = (1 << LayerMask.NameToLayer("terrain"));
        bool bhit = Physics.Raycast(screen_ray, out rh, 10000);
        Debug.Log("hit:" + bhit + " mouse pos:" + Input.mousePosition + " ray:" + screen_ray);
        if (bhit)
        {
            Debug.Log("hit" + rh.point);
        }
    }

    void OnApplicationQuit()
    {
        client_.Close();
    }

    void process_map_load_msg(IMessage msg)
    {
        map_mgr_.load_map((Msgpacket.MSG_GET_MAP_RES)msg);

        cube_all_ = new GameObject[map_mgr_.hei * map_mgr_.wid];

        for (int y = 0; y < map_mgr_.hei; y++)
        {
            for (int x = 0; x < map_mgr_.wid; x++)
            {
                var idx = y * map_mgr_.wid + x;
                //var pos = new Vector3(x, 0, map_mgr_.hei - 1 - y);
                var pos = new Vector3(x, 0, y);
                if (map_mgr_.get_block(x, y))
                {
                    var new_obj = GameObject.Instantiate(cube_block_, pos, Quaternion.identity);
                    cube_all_[idx] = new_obj;
                }
                else
                {
                    var new_obj = GameObject.Instantiate(cube_no_block_, pos, Quaternion.identity);
                    cube_all_[idx] = new_obj;
                }
            }
        }
    }
}
