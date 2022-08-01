using System;
using System.Net;
using System.Net.Sockets;
using System.Threading;
using Google.Protobuf;
using System.Runtime.InteropServices;
using System.Collections.Generic;
using UnityEngine;

public struct InterMsg
{
    public IMessage msg;
    public Msgpacket.MSG_TYPE msgtype;
}
public class TestClient
{
    [StructLayoutAttribute(LayoutKind.Sequential, CharSet = CharSet.Ansi, Pack = 1)]
    struct MSG_HEAD
    {
        public uint pack_len;
        public ushort msg_type;
        //sbyte、byte、short、ushort、int、uint、long、ulong
    }
    private Socket client_socket_ = null;

    private Thread thread_send_ = null;
    private Thread thread_recv_ = null;
    private bool b_close_ = false;


    private BlockQueue<InterMsg> send_que_ = null;
    private BlockQueue<InterMsg> recv_que_ = null;

    private Dictionary<Msgpacket.MSG_TYPE, Type> msg_parse_;

    public TestClient()
    {
        initMsgParse();

        send_que_ = new BlockQueue<InterMsg>(100);
        recv_que_ = new BlockQueue<InterMsg>(100);

        thread_send_ = new Thread(thread_send);
        thread_recv_ = new Thread(thread_recv);
    }

    public void Close()
    {
        b_close_ = true;
        client_socket_.Close();
        thread_recv_.Join();
    }

    public BlockQueue<InterMsg> GetRecvQue()
    {
        return recv_que_;
    }

    private void initMsgParse()
    {
        msg_parse_ = new Dictionary<Msgpacket.MSG_TYPE, Type>();
        msg_parse_.Add(Msgpacket.MSG_TYPE.MsgTestRes, typeof(Msgpacket.MSG_TEST));
        msg_parse_.Add(Msgpacket.MSG_TYPE.MsgLoginRes, typeof(Msgpacket.MSG_LOGIN_RES));
        msg_parse_.Add(Msgpacket.MSG_TYPE.MsgGetMapRes, typeof(Msgpacket.MSG_GET_MAP_RES));
        msg_parse_.Add(Msgpacket.MSG_TYPE.MsgPathSearchRes, typeof(Msgpacket.MSG_PATH_SEARCH_RES));
    }

    private Google.Protobuf.IMessage parseMessage(byte[] pbData, Msgpacket.MSG_TYPE msgtype)
    {
        Type msg_type;
        msg_parse_.TryGetValue(msgtype, out msg_type);
        Google.Protobuf.IMessage msgParse = (Google.Protobuf.IMessage)Activator.CreateInstance(msg_type);
        return msgParse.Descriptor.Parser.ParseFrom(pbData);
    }

    public void send_msg(Msgpacket.MSG_TYPE msgtype, IMessage proto_msg)
    {
        if (proto_msg == null)
            return;
        var inter_msg = new InterMsg();
        inter_msg.msg = proto_msg;
        inter_msg.msgtype = msgtype;

        this.send_que_.Enqueue(inter_msg);
    }

    private void thread_send()
    {
        while (!b_close_)
        {
            var msg = send_que_.Dequeue();

            send_msg_inter(msg.msgtype, msg.msg);
        }
    }

    private void thread_recv()
    {
        int headLen = Marshal.SizeOf(typeof(MSG_HEAD));

        byte[] byteRecvTmp = new byte[1024];
        byte[] Buffer = new byte[655350];        
        byte[] hdBuf = new byte[headLen];

        int readIdx = 0;
        int writeIdx = 0;

        while (!b_close_)
        {
            int recvLen = client_socket_.Receive(byteRecvTmp, 0, byteRecvTmp.Length, SocketFlags.None);
            if (0 == recvLen)
                continue;
            if (recvLen < 0)
                break;
            Array.Copy(byteRecvTmp, 0, Buffer, writeIdx, recvLen);
            writeIdx += recvLen;

            while (true)
            {
                int realBufferLen = writeIdx - readIdx;
                if (realBufferLen < headLen)
                    break;

                Array.Copy(Buffer, readIdx, hdBuf, 0, headLen);

                var obj = BytesToStuct(hdBuf, typeof(MSG_HEAD));
                if (obj == null)
                    continue;
                MSG_HEAD mh = (MSG_HEAD)obj;

                if (mh.pack_len > realBufferLen)
                    break;

                readIdx += (int)mh.pack_len;
                int bodyLen = (int)mh.pack_len - headLen;
                if (bodyLen > 0)
                {
                    byte[] bodyBuf = new byte[bodyLen];
                    
                    Array.Copy(Buffer, readIdx - bodyLen, bodyBuf, 0, bodyLen);
                    InterMsg interMsg;
                    interMsg.msgtype = (Msgpacket.MSG_TYPE)mh.msg_type;
                    interMsg.msg = parseMessage(bodyBuf, (Msgpacket.MSG_TYPE)mh.msg_type);
                    this.recv_que_.Enqueue(interMsg);
                }

                if (readIdx >= writeIdx)
                {
                    readIdx = 0;
                    writeIdx = 0;

                    Debug.Log("msg type:" + mh.msg_type + " msg len:" + mh.pack_len
                        + " writeIdx:" + writeIdx
                        + " readIdx:" + readIdx);
                }
            }
        }
    }

    public void connect(string ip, int port)
    {
        if (client_socket_ != null)
        {
            client_socket_.Close();
        }

        this.client_socket_ = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);

        client_socket_.Connect(ip, port);

        Msgpacket.MSG_LOGIN msg = new Msgpacket.MSG_LOGIN();
        msg.Id = 1;
        send_msg_inter(Msgpacket.MSG_TYPE.MsgLogin, msg);

        Msgpacket.MSG_GET_MAP msgGetMap = new Msgpacket.MSG_GET_MAP();
        send_msg_inter(Msgpacket.MSG_TYPE.MsgGetMap, msg);

        thread_send_.Start();
        thread_recv_.Start();
    }

    private void send_msg_inter(Msgpacket.MSG_TYPE msgtype, IMessage msg)
    {
        if (client_socket_ == null)
            return;

        byte[] datas = msg.ToByteArray();

        MSG_HEAD mh = new MSG_HEAD();
        mh.msg_type = (ushort)msgtype;
        mh.pack_len = (uint)(6 + datas.Length);

        byte[] head = StructToBytes(mh);
        Array.Resize(ref head, head.Length + datas.Length);
        Array.Copy(datas, 0, head, 6, datas.Length);

        client_socket_.Send(head);
    }

    public static byte[] StructToBytes(object structObj)
    {
        //得到结构体的大小
        int size = Marshal.SizeOf(structObj);
        //创建byte数组
        byte[] bytes = new byte[size];
        //分配结构体大小的内存空间
        IntPtr structPtr = Marshal.AllocHGlobal(size);
        //将结构体拷到分配好的内存空间
        Marshal.StructureToPtr(structObj, structPtr, false);
        //从内存空间拷到byte数组
        Marshal.Copy(structPtr, bytes, 0, size);
        //释放内存空间
        Marshal.FreeHGlobal(structPtr);
        //返回byte数组
        return bytes;
    }

    public static object BytesToStuct(byte[] bytes, Type type)
    {
        //得到结构体的大小
        int size = Marshal.SizeOf(type);
        //byte数组长度小于结构体的大小
        if (size > bytes.Length)
        {
            //返回空
            return null;
        }
        //分配结构体大小的内存空间
        IntPtr structPtr = Marshal.AllocHGlobal(size);
        //将byte数组拷到分配好的内存空间
        Marshal.Copy(bytes, 0, structPtr, size);
        //将内存空间转换为目标结构体
        object obj = Marshal.PtrToStructure(structPtr, type);
        //释放内存空间
        Marshal.FreeHGlobal(structPtr);
        //返回结构体
        return obj;
    }
}