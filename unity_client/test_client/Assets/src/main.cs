using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using Google.Protobuf;

public class main : MonoBehaviour
{
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
    }

    // Update is called once per frame
    void Update()
    {
        
    }
}
