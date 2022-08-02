using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using Google.Protobuf;

public class main : MonoBehaviour
{
    public GameObject cube_block_;
    public GameObject ground_;

    public GameObject[] cube_all_;

    MapMgr map_mgr_;

    Msgpacket.POS_T cur_pos_;

    TestClient client_;

    //Dictionary<Msgpacket.POS_T, GameObject> map_path_;
    System.Collections.Generic.List<GameObject> lst_path_;

    // Start is called before the first frame update
    void Start()
    {
        lst_path_ = new List<GameObject>();
        //map_path_ = new Dictionary<Msgpacket.POS_T, GameObject>();
        cur_pos_ = new Msgpacket.POS_T();
        cur_pos_.PosX = 0;
        cur_pos_.PosY = 0;

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
                case Msgpacket.MSG_TYPE.MsgPathSearchRes:
                    Debug.Log("Msgpacket.MSG_TYPE.MsgPathSearchRes");
                    process_MsgPathSearchRes((Msgpacket.MSG_PATH_SEARCH_RES)msg.msg);
                    break;
            }
        }

        if(Input.GetMouseButtonDown(0))
            check_screen_ray_hit();
    }

    private void process_MsgPathSearchRes(Msgpacket.MSG_PATH_SEARCH_RES msg)
    {
        foreach (var it in lst_path_)
            GameObject.Destroy(it);
        lst_path_.Clear();

        var new_obj = GameObject.Instantiate(cube_block_, new Vector3(msg.PosSrc.PosX + 0.5f, 0, msg.PosSrc.PosY + 0.5f), Quaternion.identity);
        var block = new_obj.GetComponent<Mapblock>();
        block.set_clr(new Color(1,0,1));
        lst_path_.Add(new_obj);

        new_obj = GameObject.Instantiate(cube_block_, new Vector3(msg.PosDst.PosX + 0.5f, 0, msg.PosDst.PosY + 0.5f), Quaternion.identity);
        block = new_obj.GetComponent<Mapblock>();
        block.set_clr(new Color(1, 0, 1));
        lst_path_.Add(new_obj);

        foreach (var it in msg.PathPos)
        {
            new_obj = GameObject.Instantiate(cube_block_, new Vector3(it.PosX + 0.5f, 0, it.PosY + 0.5f), Quaternion.identity);
            new_obj.transform.localScale = new Vector3(1, 1, 1);
            block = new_obj.GetComponent<Mapblock>();
            block.set_clr(new Color(1, 0, 0));
            lst_path_.Add(new_obj);
        }
    }

    void check_screen_ray_hit()
    {
        Ray screen_ray = Camera.main.ScreenPointToRay(Input.mousePosition);
        RaycastHit rh;

        //LayerMask lay_terrain = (1 << LayerMask.NameToLayer("terrain"));
        bool bhit = Physics.Raycast(screen_ray, out rh, 10000);
        //Debug.Log("hit:" + bhit + " mouse pos:" + Input.mousePosition + " ray:" + screen_ray);
        if (!bhit)
            return;

        Debug.Log("hit" + rh.point + " hit game obj:" + rh.collider.gameObject);

        Msgpacket.MSG_PATH_SEARCH msg = new Msgpacket.MSG_PATH_SEARCH();
        msg.PosSrc = cur_pos_;
        msg.PosDst = new Msgpacket.POS_T();
        msg.PosDst.PosX = (int)rh.point.x;
        msg.PosDst.PosY = (int)rh.point.z;
        this.client_.send_msg(Msgpacket.MSG_TYPE.MsgPathSearch, msg);

        cur_pos_ = msg.PosDst;
    }

    void OnApplicationQuit()
    {
        client_.Close();
    }

    bool check_all_around_is_block(int x, int y)
    {
        int nx = x - 1;
        int ny = y;
        {
            if (!map_mgr_.get_block(nx, ny))
                return false;
        }

        nx = x + 1;
        ny = y;
        {
            if (!map_mgr_.get_block(nx, ny))
                return false;
        }

        nx = x;
        ny = y - 1;
        {
            if (!map_mgr_.get_block(nx, ny))
                return false;
        }

        nx = x;
        ny = y + 1;
        {
            if (!map_mgr_.get_block(nx, ny))
                return false;
        }

        return true;
    }
    void process_map_load_msg(IMessage msg)
    {
        map_mgr_.load_map((Msgpacket.MSG_GET_MAP_RES)msg);

        cube_all_ = new GameObject[map_mgr_.hei * map_mgr_.wid];

        for (int y = 0; y < map_mgr_.hei; y++)
        {
            for (int x = 0; x < map_mgr_.wid; x++)
            {
                if (!map_mgr_.get_block(x, y))
                    continue;

                if (check_all_around_is_block(x, y))
                    continue;

                var idx = y * map_mgr_.wid + x;
                //var pos = new Vector3(x, 0, map_mgr_.hei - 1 - y);
                var pos = new Vector3(x + 0.5f, 0, y + 0.5f);

                var new_obj = GameObject.Instantiate(cube_block_, pos, Quaternion.identity);
                cube_all_[idx] = new_obj;

                var block = new_obj.GetComponent<Mapblock>();
                block.X = x;
                block.Y = y;
            }
        }

        var old_scale = ground_.transform.localScale;
        ground_.transform.localScale = new Vector3(map_mgr_.wid, old_scale.y, map_mgr_.hei);
        var old_pos = ground_.transform.position;
        ground_.transform.position = new Vector3(map_mgr_.wid/2, old_pos.y, map_mgr_.hei/2);
    }
}
