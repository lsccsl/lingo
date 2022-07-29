using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using Google.Protobuf;

public class main : MonoBehaviour
{
    public GameObject cube_block_;
    public GameObject cube_no_block_;

    TestClient client_;
    // Start is called before the first frame update
    void Start()
    {
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
                    break;
            }
        }
    }

    void OnApplicationQuit()
    {
        client_.Close();
    }
}
