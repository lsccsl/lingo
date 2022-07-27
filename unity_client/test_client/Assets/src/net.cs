using System;
using System.Net;
using System.Net.Sockets;
using System.Threading;
using Google.Protobuf;
using System.Runtime.InteropServices;
using System.Collections.Generic;

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
        //sbyte��byte��short��ushort��int��uint��long��ulong
    }
    private Socket client_socket_ = null;

    private Thread thread_send_ = null;
    private Thread thread_recv_ = null;


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

    private void initMsgParse()
    {
        msg_parse_ = new Dictionary<Msgpacket.MSG_TYPE, Type>();
        msg_parse_.Add(Msgpacket.MSG_TYPE.MsgTestRes, typeof(Msgpacket.MSG_TEST));
        msg_parse_.Add(Msgpacket.MSG_TYPE.MsgLoginRes, typeof(Msgpacket.MSG_LOGIN_RES));
    }

    private Google.Protobuf.IMessage parseMessage(byte[] pbData, Msgpacket.MSG_TYPE msgtype)
    {
        Type msg_type;
        msg_parse_.TryGetValue(msgtype, out msg_type);
        Google.Protobuf.IMessage msgParse = (Google.Protobuf.IMessage)Activator.CreateInstance(msg_type);
        return msgParse.Descriptor.Parser.ParseFrom(pbData);
    }

    private void thread_send()
    {
        while (true)
        {
            var msg = send_que_.Dequeue();

            send_msg(msg.msgtype, msg.msg);
        }
    }

    private void thread_recv()
    {
        int headLen = Marshal.SizeOf(typeof(MSG_HEAD));

        byte[] byteRecvTmp = new byte[1024];

        byte[] Buffer = new byte[65535];
        
        byte[] hdBuf = new byte[headLen];

        int readIdx = 0;
        int writeIdx = 0;

        while (true)
        {
            int recvLen = client_socket_.Receive(byteRecvTmp, 0, byteRecvTmp.Length, SocketFlags.None);
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
                    continue;

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
        send_msg(Msgpacket.MSG_TYPE.MsgLogin, msg);

        thread_send_.Start();
        thread_recv_.Start();
    }

    public void send_msg(Msgpacket.MSG_TYPE msgtype, IMessage msg)
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
        //�õ��ṹ��Ĵ�С
        int size = Marshal.SizeOf(structObj);
        //����byte����
        byte[] bytes = new byte[size];
        //����ṹ���С���ڴ�ռ�
        IntPtr structPtr = Marshal.AllocHGlobal(size);
        //���ṹ�忽������õ��ڴ�ռ�
        Marshal.StructureToPtr(structObj, structPtr, false);
        //���ڴ�ռ俽��byte����
        Marshal.Copy(structPtr, bytes, 0, size);
        //�ͷ��ڴ�ռ�
        Marshal.FreeHGlobal(structPtr);
        //����byte����
        return bytes;
    }

    public static object BytesToStuct(byte[] bytes, Type type)
    {
        //�õ��ṹ��Ĵ�С
        int size = Marshal.SizeOf(type);
        //byte���鳤��С�ڽṹ��Ĵ�С
        if (size > bytes.Length)
        {
            //���ؿ�
            return null;
        }
        //����ṹ���С���ڴ�ռ�
        IntPtr structPtr = Marshal.AllocHGlobal(size);
        //��byte���鿽������õ��ڴ�ռ�
        Marshal.Copy(bytes, 0, structPtr, size);
        //���ڴ�ռ�ת��ΪĿ��ṹ��
        object obj = Marshal.PtrToStructure(structPtr, type);
        //�ͷ��ڴ�ռ�
        Marshal.FreeHGlobal(structPtr);
        //���ؽṹ��
        return obj;
    }
}