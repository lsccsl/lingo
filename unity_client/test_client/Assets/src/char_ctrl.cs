using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class char_ctrl : MonoBehaviour
{
    public Vector3 target_pos_ = new Vector3(1,0,1);
    public float move_speed_ = 5.0f;

    public bool has_target_ = false;

    public float last_dis = float.MaxValue;

    List<Vector3> lst_target_ = new List<Vector3>();

    public Vector3 TargetPos
    {
        set { set_target_pos(value); }
        get { return target_pos_; }
    }

    // Start is called before the first frame update
    void Start()
    {
    }

    public void add_target_pos(Vector3 pos)
    {
        lst_target_.Add(pos);
    }

    // Update is called once per frame
    void Update()
    {
        move_to_target();
    }

    void move_to_target()
    {
        if (!has_target_)
        {
            if (lst_target_.Count > 0)
            {
                var lst = lst_target_.GetRange(0, 1);
                lst_target_.RemoveAt(0);
                foreach (var it_pos in lst)
                {
                    set_target_pos(it_pos);
                    break;
                }
            }
            return;
        }

        transform.Translate(new Vector3(0, 0, move_speed_ * Time.deltaTime));

        var pos = this.transform.position;
        float dis = Vector2.Distance(new Vector2(pos.x, pos.z), new Vector2(target_pos_.x, target_pos_.z));
        if (dis < 1 || last_dis < dis)
        {
            last_dis = float.MaxValue;
            has_target_ = false;
            return;
        }
        last_dis = dis;
    }

    public void set_target_pos(Vector3 pos)
    {
        has_target_ = true;
        target_pos_ = pos;

        Vector3 old_euler = this.transform.eulerAngles;
        this.transform.LookAt(target_pos_);
/*        Transform n = m_current_node.transform;

        Vector3 old_euler = this.transform.eulerAngles;
        this.transform.LookAt(n);
        Vector3 new_euler = this.transform.eulerAngles;

        float new_y_angle = Mathf.MoveTowardsAngle(old_euler.y, new_euler.y, m_rotate_speed * Time.deltaTime);

        this.transform.eulerAngles = new Vector3(0, new_y_angle, 0);
*/
    }
}
