using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class main_nav : MonoBehaviour
{
    TestClient client_;
    private float t_heart_beat_ = 0;
    Msgpacket.PROTO_VEC_3F cur_pos_;

    public LineRenderer m_lineRenderer;

    public GameObject pref_obstacle_;

    // Start is called before the first frame update
    void Start()
    {
        m_lineRenderer = this.gameObject.GetComponent<LineRenderer>();
        m_lineRenderer.startColor = Color.blue;
        m_lineRenderer.endColor = Color.red;

        cur_pos_ = new Msgpacket.PROTO_VEC_3F();
        cur_pos_.X = 0;
        cur_pos_.Y = 0;

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
                    process_nav_search_res((Msgpacket.MSG_NAV_SEARCH_RES)msg.msg);
                    break;
                case Msgpacket.MSG_TYPE.MsgNavAddObstacleRes:
                    process_MSG_NAV_ADD_OBSTACLE_RES((Msgpacket.MSG_NAV_ADD_OBSTACLE_RES)msg.msg);
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
        if (Input.GetMouseButtonDown(1))
            check_add_obstacle();
    }

    void process_nav_search_res(Msgpacket.MSG_NAV_SEARCH_RES msg)
    {
        Debug.Log("process_nav_search_res");
        {
            m_lineRenderer.positionCount = msg.PathPos.Count;

            int idx = 0;
            foreach (var it in msg.PathPos)
            {
                m_lineRenderer.SetPosition(idx, new Vector3(it.X, it.Y, it.Z));
                idx++;
            }
        }
    }

    void process_MSG_NAV_ADD_OBSTACLE_RES(Msgpacket.MSG_NAV_ADD_OBSTACLE_RES msg)
    {
        Debug.Log("process_MSG_NAV_ADD_OBSTACLE_RES");
        var gobj_obstacle = GameObject.Instantiate(pref_obstacle_, new Vector3(msg.Center.X, msg.Center.Y, msg.Center.Z), Quaternion.EulerRotation(0,(float)(30.0/360.0 * 2.0 * 3.14),0));
        var com_obstacle = gobj_obstacle.GetComponent<obstacle>();
        com_obstacle.set_scale(new Vector3(msg.HalfExt.X, msg.HalfExt.Y, msg.HalfExt.Z) * 2);
    }

    void check_add_obstacle()
    {
        Ray screen_ray = Camera.main.ScreenPointToRay(Input.mousePosition);
        RaycastHit rh;

        bool bhit = Physics.Raycast(screen_ray, out rh, 10000);
        if (!bhit)
            return;

        Msgpacket.MSG_NAV_ADD_OBSTACLE msg = new Msgpacket.MSG_NAV_ADD_OBSTACLE();
        msg.Center = new Msgpacket.PROTO_VEC_3F();
        msg.Center.X = rh.point.x;
        msg.Center.Y = rh.point.y;
        msg.Center.Z = rh.point.z;
        msg.HalfExt = new Msgpacket.PROTO_VEC_3F();
        msg.HalfExt.X = 10;
        msg.HalfExt.Y = 10;
        msg.HalfExt.Z = 5;

        msg.YRadian = (float)(30.0 / 360.0 * 2.0 * 3.14);

        this.client_.send_msg(Msgpacket.MSG_TYPE.MsgNavAddObstacle, msg);
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
        msg.PosSrc.Y = 1.0f;
        msg.PosDst = new Msgpacket.PROTO_VEC_3F();
        msg.PosDst.X = rh.point.x;
        msg.PosDst.Y = 1.0f;// rh.point.y;
        msg.PosDst.Z = rh.point.z;

/*        msg.PosSrc = new Msgpacket.POS_3F();
        msg.PosSrc.PosX = 702.190918f;
        msg.PosSrc.PosY = 1.53082275f;
        msg.PosSrc.PosZ = 635.378662f;
        msg.PosDst = new Msgpacket.POS_3F();
        msg.PosDst.PosX = 710.805664f; 
        msg.PosDst.PosY = 1.00000000f;
        msg.PosDst.PosZ = 851.753296f;*/
        this.client_.send_msg(Msgpacket.MSG_TYPE.MsgNavSearch, msg);

        cur_pos_ = msg.PosDst;
    }
}
