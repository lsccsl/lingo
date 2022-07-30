using System;
using Google.Protobuf;
using UnityEngine;

public class MapMgr
{
    int wid_;
    int hei_;
    int pitch_;
    byte[] map_data_;

    public int wid { get { return wid_; } }
    public int hei { get { return hei_; } }
    public int pitch { get { return pitch_; } }

    public void load_map(Msgpacket.MSG_GET_MAP_RES msg)
    {
        Debug.Log("load map");
        wid_ = msg.MapWid;
        hei_ = msg.MapHigh;
        pitch_ = msg.MapPitch;
        map_data_ = msg.MapData.ToByteArray();
    }


    public bool get_block(int x, int y)
    {
        if (x < 0 || x >= wid_)
            return true;
        if (y < 0 || y >= hei_)
            return true;

        int byte_idx = y * pitch_ + x / 8;
        int idx_bit = 7 - x % 8;
        byte pos_byte = map_data_[byte_idx];
        var pos_bit = pos_byte & (1 << idx_bit);

        return pos_bit == 0;
    }
}
