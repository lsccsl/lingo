using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class main_nav : MonoBehaviour
{
    TestClient client_;
    private float t_heart_beat_ = 0;
    Msgpacket.POS_3F cur_pos_;

    // Start is called before the first frame update
    void Start()
    {
        cur_pos_ = new Msgpacket.POS_3F();
        cur_pos_.PosX = 0;
        cur_pos_.PosY = 0;

        client_ = new TestClient();
        client_.connect("192.168.2.129", 2003);
    }

    void OnApplicationQuit()
    {
        client_.Close();
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
                case Msgpacket.MSG_TYPE.MsgNavSearchRes:
                    Debug.Log("Msgpacket.MSG_TYPE.MsgNavSearchRes");
                    break;
            }
        }

        t_heart_beat_ += Time.deltaTime;
        if (t_heart_beat_ > 30)
        {
            t_heart_beat_ = 0;
            Msgpacket.MSG_HEARTBEAT msgHB = new Msgpacket.MSG_HEARTBEAT();
            this.client_.send_msg(Msgpacket.MSG_TYPE.MsgHeartbeat, msgHB);
        }

        if (Input.GetMouseButtonDown(0))
            check_screen_ray_hit();
    }

    void check_screen_ray_hit()
    {
        Ray screen_ray = Camera.main.ScreenPointToRay(Input.mousePosition);
        RaycastHit rh;

        //LayerMask lay_terrain = (1 << LayerMask.NameToLayer("terrain"));
        bool bhit = Physics.Raycast(screen_ray, out rh, 10000);
        if (!bhit)
            return;

        Debug.Log("hit" + rh.point + " hit game obj:" + rh.collider.gameObject);

        Msgpacket.MSG_NAV_SEARCH msg = new Msgpacket.MSG_NAV_SEARCH();
        msg.PosSrc = cur_pos_;
        msg.PosDst = new Msgpacket.POS_3F();
        msg.PosDst.PosX = rh.point.x;
        msg.PosDst.PosY = rh.point.y;
        msg.PosDst.PosZ = rh.point.z;
        this.client_.send_msg(Msgpacket.MSG_TYPE.MsgNavSearch, msg);

        cur_pos_ = msg.PosDst;
    }
}
