/**
 * @file runtime_view
 */
#include "runtime_view.h"

/**
 * @brief constructor
 */
rt_view::rt_view(const char * ipc_name):
	us_(ipc_name ? ipc_name : "___runtime_view_")
{
	this->hlsn_ = MyListernerConstruct(NULL, 1024);

	this->us_.SetToNoBlock();

	event_handle_t e = {0};
	e.input = rt_view::_handle_input;
	e.context_data = (unsigned long)this;

	MyListernerAddFD(this->hlsn_, us_.GetFd(), E_FD_READ, &e);

	MyListernerRun(this->hlsn_);
}

/**
 * @brief destructor
 */
rt_view::~rt_view()
{
	MyListernerDestruct(this->hlsn_);
}

/**
 * @brief ÊäÈë»Øµ÷
 */
int rt_view::_handle_input(unsigned long context_data, int fd)
{
	rt_view * rtv = (rt_view *)context_data;

	int ret = 0;
	do
	{
		char actemp[4] = {0};
		ret = rtv->us_.Read(actemp, 
			sizeof(actemp), 
			NULL, 
			0);

		if(ret > 0)
			rtv->command(*((int *)actemp));

	}while(ret > 0);

	return 0;
}

/**
 * @brief ËÍÃüÁî
 */
void rt_view::send_cmd(int cmd, const char * ipc_name)
{
	CMyUnixSocket us("temp");
	us.Write(&cmd, 
		sizeof(cmd), 
		ipc_name ? ipc_name : "___runtime_view_");
}
