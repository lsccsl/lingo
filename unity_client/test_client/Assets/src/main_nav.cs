using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class main_nav : MonoBehaviour
{
    TestClient client_;

    // Start is called before the first frame update
    void Start()
    {
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
            }
        }
    }
}
