using System;
using System.Net;
using System.Net.Sockets;
using System.Threading;
using Google.Protobuf;
using System.Runtime.InteropServices;

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
    }

    public void send_msg(Msgpacket.MSG_TYPE msgtype, IMessage msg)
    {
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
}