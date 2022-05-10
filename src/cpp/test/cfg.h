#pragma once

#include <string>

struct testcfg
{
	std::string remote_ip_;
	int remote_port_;

	std::string local_ip_;
};

extern testcfg __global_cfg_;