/**
* @file mysqlwrapper.h
* @brief 
* @author linsc
*/
#include <strstream>
#include <sstream>
#include "mysql.h"
#include "mysqlwrapper.h"
#include "mylogex.h"
#include "myos_cpp.h"

/**
* @brief constructor
*/
MysqlWrapper::MysqlWrapper(const int8 * ip, uint32 port,
	const int8 * user, const int8 * pwd,
	const int8 * dbname):
		hdb_(NULL),db_ip_(ip ? ip : "0.0.0.0"),db_port_(port),db_user_(user ? user : ""),db_pwd_(pwd ? pwd : ""),db_name_(dbname ? dbname : "")
{
	MYLOG_DEBUG(("MysqlWrapper::MysqlWrapper db_name_:%s db_ip_:%s db_port_:%d user:%s pwd:%s",
		db_name_.c_str(), db_ip_.c_str(), db_port_, this->db_user_.c_str(), this->db_pwd_.c_str()));
}

/**
* @brief constructor
*/
MysqlWrapper::~MysqlWrapper()
{
	MYLOG_DEBUG(("MysqlWrapper::~MysqlWrapper"));
	this->uninit();
}

/**
* @brief 初始化
*/
int MysqlWrapper::init()
{
	MYLOG_DEBUG(("MysqlWrapper::init"));

	MYSQL * hreal = NULL;
	if(this->hdb_)
	{
		MYLOG_DEBUG(("already have connect"));

		if(0 == mysql_ping((MYSQL *)this->hdb_))
		{
			MYLOG_DEBUG(("ping ok"));
			return 0;
		}

		MYLOG_INFO(("lost connect to mysql, will reconnect ..."));

		this->uninit();
	}

	MYLOG_DEBUG(("connect to mysql"));

	hreal = mysql_init(NULL);
	if(NULL == hreal)
	{
		MYLOG_INFO(("init db handle err..."));
		return -1;
	}

	this->hdb_ = hreal;

	my_bool reconn = 1;
    mysql_options(hreal, MYSQL_OPT_RECONNECT, &reconn);
    mysql_options(hreal, MYSQL_SET_CHARSET_NAME, "utf8");

	hreal = mysql_real_connect(hreal,
		this->db_ip_.c_str(),
		this->db_user_.c_str(), this->db_pwd_.c_str(), this->db_name_.c_str(),
		this->db_port_, NULL, 0);

	if(NULL == hreal)
	{
		MYLOG_INFO(("connect to mysql err [%s:%u] user:%s pwd:%s ...",
			this->db_ip_.c_str(), this->db_port_, this->db_user_.c_str(), this->db_pwd_.c_str()));
		goto MysqlWrapper_init_err_;
	}

	if(0 != mysql_select_db(hreal, this->db_name_.c_str()))
	{
		MYLOG_INFO(("select db %s err ...", this->db_name_.c_str()));
		goto MysqlWrapper_init_err_;
	}

	MYLOG_DEBUG(("connect ok"));

	return 0;

MysqlWrapper_init_err_:

	this->uninit();
	return -1;
}

/**
* @brief 反初始化
*/
void MysqlWrapper::uninit()
{
	MYLOG_DEBUG(("MysqlWrapper::uninit"));

	if(this->hdb_)
	{
		mysql_close((MYSQL *)this->hdb_);
		this->hdb_ = NULL;
	}
}

/**
* @brief 执行sql语句
*/
int32 MysqlWrapper::query(const int8 * pcsql)
{
	//MYLOG_DEBUG(("MysqlWrapper::query pcsql:%s", pcsql ? pcsql : "null"));

	if(NULL == this->hdb_)
	{
		MYLOG_INFO(("mysql init err"));
		this->init();
		if(NULL == this->hdb_)
		{
			MYLOG_INFO(("connect to db err"));
			return -1;
		}
	}

	if(NULL == pcsql)
	{
		MYLOG_INFO(("err param"));
		return -1;
	}

	if(0 != mysql_query((MYSQL *)this->hdb_, pcsql))
	{
		MYLOG_INFO(("exe [%s] err [%d:%s] do it again", pcsql ? pcsql : "",
			mysql_errno((MYSQL *)this->hdb_), mysql_error((MYSQL *)this->hdb_)));

		this->init();
		if(NULL == this->hdb_)
		{
			MYLOG_INFO(("connect to db err"));
			return -1;
		}
		if(0 != mysql_query((MYSQL *)this->hdb_, pcsql))
		{
			MYLOG_INFO(("exe [%s] err [%d:%s] fail again ...", pcsql ? pcsql : "",
				mysql_errno((MYSQL *)this->hdb_), mysql_error((MYSQL *)this->hdb_)));
			return -1;
		}
	}

	//MYLOG_DEBUG(("exe end"));

	return 0;
}

/**
* @brief 获取结果集
*/
int32 MysqlWrapper::get_result(std::vector<std::vector<std::string> >& vresult)
{
	//MYLOG_DEBUG(("MysqlWrapper::get_result"));

	if(NULL == this->hdb_)
	{
		MYLOG_INFO(("mysql init err"));
		return -1;
	}

	vresult.clear();

	MYSQL_RES * pmr = mysql_store_result((MYSQL *)this->hdb_);

	uint32 field_count = mysql_num_fields(pmr);
	//MYLOG_DEBUG(("field_count:%u", field_count));

	MYSQL_ROW row = NULL;
	uint32 row_count = 0;
	while (NULL != (row = mysql_fetch_row(pmr)))
	{
		vresult.resize(row_count + 1);
		vresult[row_count].resize(field_count);
		for (unsigned int i = 0; i < field_count; ++i)
		{
			vresult[row_count][i] = row[i] ? row[i] : "";
		}

		row_count ++;
	}

	mysql_free_result(pmr);

	return 0;
}

/**
* @brief 获取db时间 
*/
int32 MysqlWrapper::get_dbtime(uint64& time_db)
{
	//MYLOG_DEBUG(("MysqlWrapper::get_dbtime"));

	if(0 != this->query("select unix_timestamp()"))
	{
		MYLOG_INFO(("query db time err ..."));
		return -1;
	}

	std::vector<std::vector<std::string> > vresult;
	if(0 != this->get_result(vresult))
	{
		MYLOG_INFO(("get result err ..."));
		return -1;
	}

	if(1 != vresult.size())
	{
		MYLOG_INFO(("result count is wrong"));
		return -1;
	}

	myos::string2uint64(vresult[0][0], time_db);
	return 0;
}

/**
* @brief 查询表是否存在
*/
int32 MysqlWrapper::table_exist(const int8 * table_name, int32& bexist)
{
	MYLOG_DEBUG(("MysqlWrapper::table_exist table_name:%s", table_name ? table_name : ""));

	if(NULL == this->hdb_)
	{
		MYLOG_INFO(("mysql init err"));
		return -1;
	}

	if(NULL == table_name)
	{
		MYLOG_INFO(("table name is null"));
		return -1;
	}

	MYSQL_RES * pmr = mysql_list_tables((MYSQL *)this->hdb_, table_name);
	if(NULL == pmr)
	{
		MYLOG_DEBUG(("list table fail"));
		return -1;
	}

	bexist = mysql_num_rows(pmr) > 0;
	return 0;
}

/**
* @brief 复制表结构
*/
int32 MysqlWrapper::copy_table_struct(const int8 * table_src, const int8 * table_dst)
{
	MYLOG_DEBUG(("MysqlWrapper::copy_table_struct table_src:%s table_dst:%s", table_src ? table_dst : "", table_dst ? table_dst : ""));

	if(NULL == this->hdb_)
	{
		MYLOG_INFO(("mysql init err"));
		return -1;
	}

	if(NULL == table_src || NULL == table_dst)
	{
		MYLOG_DEBUG(("invalid table name"));
		return -1;
	}

	std::ostringstream oss;
	oss << "create table " << table_dst << " as select * from " << table_src << " where 1=0" << std::ends;
	MYLOG_DEBUG(("%s", oss.str().c_str()));

	return this->query(oss.str().c_str());
}

/**
* @brief 获取最后播入的id
*/
int32 MysqlWrapper::get_last_insert_id(uint64& last_id)
{
	MYLOG_DEBUG(("MysqlWrapper::get_last_insert_id"));

	last_id = mysql_insert_id((MYSQL *)this->hdb_);

	int32 err = mysql_errno((MYSQL *)this->hdb_);
	MYLOG_DEBUG(("err:%d", err));
	if(0 != err)
		return -1;

	return 0;
}

/**
* @brief 获取最后一次查询受影响的行数
*/
int32 MysqlWrapper::get_affected_rows(uint32& affected_rows)
{
	MYLOG_DEBUG(("MysqlWrapper::get_affected_rows"));

	my_ulonglong ret = mysql_affected_rows((MYSQL *)this->hdb_);

	if(-1 == ret)
	{
		return -1;
	}

	affected_rows = ret;
	return 0;
}













