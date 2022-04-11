/**
* @file mssqlwrapper.h
* @brief ms的db编程接口封装 select @@IDENTITY
* @author linsc
*/
#include "mssqlwrapper.h"
#include <string.h>
#include <memory.h>
#include <time.h>
#include <assert.h>
#include <strstream>
#include <sstream>
#ifdef WIN32
#define MSDBLIB
#define snprintf _snprintf
#endif
#include "sybdb.h"
#include "mylogex.h"
#include "myos_cpp.h"

struct mssql_handle_inter
{
	mssql_handle_inter():loginrec_(NULL),dbprocess_(NULL){}

	/**
	 * @brief 登陆信息
	 */
	LOGINREC * loginrec_;

	/**
	 * @brief 真实的数据库句柄
	 */
	DBPROCESS * dbprocess_;
};

/**
* @brief err call back
*/
static int __mssql_err_callback(DBPROCESS * dbproc, int severity, int dberr, int oserr, char *dberrstr, char *oserrstr)
{
	MYLOG_INFO(("__mssql_err_callback dbproc:%x severity:%d dberr:%d oserr:%d dberrstr:[%s] oserrstr:[%s]",
		dbproc, severity, dberr, oserr, dberrstr ? dberrstr : "", oserrstr ? oserrstr : ""));
	return 0;
}

static void tds_init()
{
	static int __binit_ = 0;

	if(!__binit_)
	{
		MYLOG_DEBUG(("call tds init"));

		dbinit();
		dberrhandle(__mssql_err_callback);
		__binit_ = 1;
	}
}

/**
* @brief constructor
*/
mssqlwrapper::mssqlwrapper(const int8 * ip, uint32 port,
	const int8 * user, const int8 * pwd,
	const int8 * dbname):db_ip_(ip ? ip : "0.0.0.0"),db_port_(port),db_user_(user ? user : ""),db_pwd_(pwd ? pwd : ""),db_name_(dbname ? dbname : ""),h_(NULL)
{
	MYLOG_DEBUG(("mssqlwrapper::mssqlwrapper [%s:%u] db_user_:%s db_pwd_:%s db_name_:%s",
		db_ip_.c_str(), db_port_, db_user_.c_str(), db_pwd_.c_str(), db_name_.c_str()));
	tds_init();
}

/**
* @brief constructor
*/
mssqlwrapper::~mssqlwrapper()
{
	MYLOG_DEBUG(("mssqlwrapper::~mssqlwrapper"));
	this->uninit();
}

/**
* @brief 初始化
*/
int mssqlwrapper::init()
{
	MYLOG_DEBUG(("mssqlwrapper::init"));
	this->uninit();

	this->h_ = new mssql_handle_inter;

	this->h_->loginrec_ = dblogin();
	DBSETLUSER(this->h_->loginrec_, this->db_user_.c_str());
	DBSETLPWD(this->h_->loginrec_, this->db_pwd_.c_str());
	//DBSETLCHARSET(this->h_->loginrec_, "UTF-8"); //--- 不乱码,但也看不清楚到底取出的是什么码
	DBSETLCHARSET(this->h_->loginrec_, "cp936");//ms推荐的编码方式
	//DBSETLCHARSET(this->h_->loginrec_, "UCS-2BE");   -- can't even conn
	//DBSETLCHARSET(this->h_->loginrec_, "UCS-2LE");  -- can't even conn
	//DBSETLCHARSET(this->h_->loginrec_, "ISO-8859-1"); -- 无效,乱码

	std::string srv = this->db_ip_;
	std::string port;
	myos::uint322string(this->db_port_, port);
	srv.append(":");
	srv.append(port);
	MYLOG_DEBUG(("srv:%s", srv.c_str()));

	this->h_->dbprocess_ = dbopen(this->h_->loginrec_, srv.c_str());

	if(this->h_->dbprocess_ == FAIL)
	{
		MYLOG_INFO(("connect MS SQL SERVER fail \r\n"));
		goto mssqlwrapper__init_err_;
	}
	MYLOG_DEBUG(("ConnectEMS conect MS SQL SERVER success\r\n"));

	if(dbuse(this->h_->dbprocess_, this->db_name_.c_str()) == FAIL)
	{
		MYLOG_INFO(("Open database name fail\r\n"));
		goto mssqlwrapper__init_err_;
	}

	MYLOG_INFO(("Open database name success\r\n"));

	return 0;

mssqlwrapper__init_err_:

	this->uninit();

	return -1;
}

/**
* @brief 反初始化
*/
void mssqlwrapper::uninit()
{
	MYLOG_DEBUG(("mssqlwrapper::uninit"));

	if(NULL == this->h_)
	{
		MYLOG_DEBUG(("db wrapper handle is null"));
		return;
	}

	if(this->h_->dbprocess_)
	{
		MYLOG_DEBUG(("close real db handle"));
		dbclose(this->h_->dbprocess_);

		this->h_->dbprocess_ = NULL;
	}

	if(this->h_->loginrec_)
	{
		MYLOG_DEBUG(("free login rec"));
		dbloginfree(this->h_->loginrec_);

		this->h_->loginrec_ = NULL;
	}

	delete this->h_;

	this->h_ = NULL;
}

/**
* @brief 执行sql语句,0:成功 -1:失败
*/
int32 mssqlwrapper::query(const int8 * pcsql)
{
	MYLOG_DEBUG(("mssqlwrapper::query [%s]", pcsql ? pcsql : ""));

	if(NULL == pcsql)
	{
		MYLOG_INFO(("pcsql is null"));
	}

	if(NULL == this->h_)
	{
		if(0 != this->init())
		{
			MYLOG_INFO(("init db conn fail ..."));
			return -1;
		}
	}

	MYLOG_DEBUG(("this->h_:%x", this->h_));
	assert(this->h_);
	assert(this->h_->dbprocess_);

	dbcmd(this->h_->dbprocess_, pcsql);
	if(dbsqlexec(this->h_->dbprocess_) == FAIL)
	{
		MYLOG_INFO(("exe [%s] err", pcsql));

		if(0 != this->init())
		{
			MYLOG_INFO(("init db err"));
			return -1;
		}

		dbcmd(this->h_->dbprocess_, pcsql);
		if(dbsqlexec(this->h_->dbprocess_) == FAIL)
		{
			MYLOG_INFO(("exe [%s] err", pcsql));
			return -1;
		}

	}

	return 0;
}

/**
* @brief 列信息
*/
struct tds_col
{ 
	std::vector<int8> name; 
	std::vector<int8> buffer; 
	uint32 type, size;
	int32 status;

	union{
		DBDATETIME dt_;
		DBDATETIME4 sdt_;
	}dt;
};

/**
* @brief 获取结果集
*/
int32 mssqlwrapper::get_result(std::vector<std::vector<std::string> >& vresult)
{
	MYLOG_DEBUG(("mssqlwrapper::get_result"));

	vresult.clear();

	if(NULL == this->h_)
	{
		MYLOG_INFO(("h is null"));
		return -1;
	}

	if(NULL == this->h_->dbprocess_)
	{
		MYLOG_INFO(("db handle is null"));
		return -1;
	}

	DBINT erc = dbresults(this->h_->dbprocess_);
	if(erc == NO_MORE_RESULTS)
	{
		MYLOG_INFO(("no result"));
		return 0;
	}

	/* get colum type */
	if (erc != SUCCEED)
	{
		MYLOG_INFO(("dbresults failed"));
		return -1;
	}

	uint32 ncols = dbnumcols(this->h_->dbprocess_);

	std::vector<tds_col> v_col(ncols);

	if(v_col.size() != ncols)
	{
		MYLOG_INFO(("can't malloc col, not enough memory"));
		return -1;
	}

	{
		for(uint32 i = 0; i < ncols; i++)
		{
			//MYLOG_DEBUG(("get col %d", i));

			v_col[i].type = dbcoltype(this->h_->dbprocess_, i + 1);

			MYLOG_DEBUG(("get type %d", v_col[i].type));

			switch(v_col[i].type)
			{
			case SYBDATETIME:
				{
					//MYLOG_DEBUG(("SYBDATETIME"));
					v_col[i].size = dbcollen(this->h_->dbprocess_, i + 1);
					//MYLOG_DEBUG(("size:%d", v_col[i].size));
					erc = dbbind(this->h_->dbprocess_, (i + 1), DATETIMEBIND,
						sizeof(v_col[i].dt.dt_), (BYTE*)(&v_col[i].dt.dt_));
				}
				break;

			case SYBDATETIME4:
				{
					//MYLOG_DEBUG(("SYBDATETIME4"));
					v_col[i].size = dbcollen(this->h_->dbprocess_, i + 1);
					//MYLOG_DEBUG(("size:%d", v_col[i].size));
					erc = dbbind(this->h_->dbprocess_, (i + 1), SMALLDATETIMEBIND,
						sizeof(v_col[i].dt.sdt_), (BYTE*)(&v_col[i].dt.sdt_));
				}
				break;

			default:
				{
					if(SYBCHAR != v_col[i].type)
					{
						/* 非字符,亦转成字符类型 */
						//MYLOG_DEBUG(("not char"));
						v_col[i].size = dbwillconvert(v_col[i].type, SYBCHAR);

					}
					else
					{
						/* 本身是存的是字符类型 */
						//MYLOG_DEBUG(("char"));
						v_col[i].size = dbcollen(this->h_->dbprocess_, i + 1);
					}

					//MYLOG_DEBUG(("size:%d", v_col[i].size));
					v_col[i].buffer.resize(v_col[i].size + 3);
					if(v_col[i].buffer.size() < v_col[i].size)
					{
						MYLOG_INFO(("malloc err"));
						return -1;
					}

					memset(&v_col[i].buffer[0], 0, v_col[i].buffer.size());
					erc = dbbind(this->h_->dbprocess_, (i + 1), NTBSTRINGBIND,
						(DBINT)(v_col[i].buffer.size() - 1), (BYTE*)(&v_col[i].buffer[0]));
				}
				break;
			}

			//MYLOG_DEBUG(("erc:%d\r\n", erc));

			if (erc == FAIL)
			{
				MYLOG_INFO(("dbbind failed"));
				return -1;
			}
			erc = dbnullbind(this->h_->dbprocess_, (i + 1), &v_col[i].status);
			if (erc == FAIL)
			{
				MYLOG_INFO(("dbnullbind failed"));
				return -1;
			}
		}
	}

	int row_code;
	uint32 row_count = 0;

	//MYLOG_DEBUG(("bind end, begin row"));

	while((row_code = dbnextrow(this->h_->dbprocess_)) != NO_MORE_ROWS)
	{
		//MYLOG_DEBUG(("row_code:%d", row_code));
		switch (row_code)
		{
		case REG_ROW:
			{
				//MYLOG_DEBUG(("REG_ROW"));

				vresult.resize(row_count + 1);
				vresult[row_count].resize(v_col.size());

				for(uint32 j = 0; j < v_col.size(); j++)
				{
					//MYLOG_DEBUG(("row:%d col:%d status:%d", row_count, j, v_col[j].status));
					if(-1 == v_col[j].status)
					{
						//MYLOG_DEBUG(("col value is null"));
						vresult[row_count][j] = "";
					}
					else
					{
						//MYLOG_DEBUG(("type:%d", v_col[j].type));

						switch(v_col[j].type)
						{
						case SYBDATETIME:
							{
								//MYLOG_DEBUG(("SYBDATETIME %u %u", v_col[j].dt.dt_.dtdays, v_col[j].dt.dt_.dttime));
								//MYLOG_DUMP_BIN(&v_col[j].dt.dt_, sizeof(v_col[j].dt.dt_));

								DBDATEREC di = {0};
								dbdatecrack(this->h_->dbprocess_, &di, &(v_col[j].dt.dt_));
								tm t = {0};
								t.tm_year = di.year - 1900;
								t.tm_mon = di.month - 1;
								t.tm_mday = di.day;
								t.tm_hour = di.hour;
								t.tm_min = di.minute;
								t.tm_sec = di.second;

								uint64 ulltime = mktime((struct tm *)&t);

								vresult[row_count][j].resize(64);
								int32 ret = snprintf((char *)vresult[row_count][j].data(), vresult[row_count][j].size() - 1,
									"%llu", ulltime);
								vresult[row_count][j].resize(ret);
							}
							break;

						case SYBDATETIME4:
							{
								//MYLOG_DEBUG(("SYBDATETIME4 %u %u", v_col[j].dt.sdt_.days, v_col[j].dt.sdt_.minutes));
								//MYLOG_DUMP_BIN(&v_col[j].dt.sdt_, sizeof(v_col[j].dt.sdt_));

								DBDATETIME dt = {0};
								dt.dtdays = v_col[j].dt.sdt_.days;
								dt.dttime = v_col[j].dt.sdt_.minutes * 60 * 300;

								DBDATEREC di = {0};
								dbdatecrack(this->h_->dbprocess_, &di, &dt);

								//MYLOG_DEBUG(("year:%d month:%d day:%d dayofyear:%d weekday:%d hour:%d minute:%d second:%d millisecond:%d tzone:%d",
								//	di.year, di.month, di.day, di.dayofyear, di.weekday, di.hour, di.minute, di.second, di.millisecond, di.tzone));

								tm t = {0};
								t.tm_year = di.year - 1900;
								t.tm_mon = di.month - 1;
								t.tm_mday = di.day;
								t.tm_hour = di.hour;
								t.tm_min = di.minute;
								t.tm_sec = di.second;

								uint64 ulltime = mktime((struct tm *)&t);

								vresult[row_count][j].resize(64);
								int32 ret = snprintf((char *)vresult[row_count][j].data(), vresult[row_count][j].size() - 1,
									"%llu", ulltime);
								vresult[row_count][j].resize(ret);
							}
							break;

						default:
							//MYLOG_DEBUG(("col val:%s", &v_col[j].buffer[0]));
							vresult[row_count][j] = &v_col[j].buffer[0];
							break;
						}
					}
				}

				row_count ++;
			}
			break;

			case BUF_FULL:
				MYLOG_DEBUG(("BUF_FULL"));
				break;

			case FAIL:
				MYLOG_DEBUG(("dbresults failed"));
				return -1;
				break;

			default:
				MYLOG_DEBUG(("Data for computeid %d ignored\n", row_code));
				break;
		}
	}

	MYLOG_DEBUG(("fetch result:%d", vresult.size()));

	return 0;
}

/**
* @brief 获取db版本信息  
*/
const int8 * mssqlwrapper::get_dbver()
{
	MYLOG_DEBUG(("mssqlwrapper::get_dbver"));

	static int8 ver_temp[256] = {0};

	std::vector<std::vector<std::string> > vresult;
	this->query("SELECT SERVERPROPERTY('productversion'), SERVERPROPERTY ('productlevel'), SERVERPROPERTY ('edition');");
	this->get_result(vresult);

	std::string temp;
	for(uint32 i = 0; i < vresult.size(); i ++)
	{
		for(uint32 j = 0; j < vresult[i].size(); j ++)
		{
			temp.append(vresult[i][j]);
			temp.append(" ");
		}
	}

	strncpy(ver_temp, temp.c_str(), sizeof(ver_temp) - 1);
	return ver_temp;
}

/**
* @brief 查询表是否存在
*/
int32 mssqlwrapper::table_exist(const int8 * table_name, int32& bexist)
{
	MYLOG_DEBUG(("mssqlwrapper::table_exist table_name:%s", table_name ? table_name : ""));

	if(NULL == table_name)
	{
		MYLOG_INFO(("table name is null"));
		return -1;
	}

	std::string temp = "select count(*) from sysobjects where name = '";
	temp.append(table_name);
	temp.append("'");
	MYLOG_DEBUG(("query table exist:%s", temp.c_str()));
	if(0 != this->query(temp.c_str()))
	{
		MYLOG_INFO(("query err:%s", temp.c_str()));
		return -1;
	}

	std::vector<std::vector<std::string> > vresult;
	this->get_result(vresult);

	if(vresult.size() < 1)
	{
		MYLOG_INFO(("db lib err ..."));
		return -1;
	}

	if(vresult[0].size() < 1)
	{
		MYLOG_INFO(("db lib err ..."));
		return -1;
	}

	MYLOG_DEBUG(("%s", vresult[0][0].c_str()));

	bexist = (vresult[0][0] != "0");
	return 0;
}

/**
* @brief 复制表结构 select *  into xxx from  xx where 1=0
*/
int32 mssqlwrapper::copy_table_struct(const int8 * table_src, const int8 * table_dst)
{
	MYLOG_DEBUG(("mssqlwrapper::copy_table_struct table_src:%s table_dst:%s", table_src ? table_dst : "", table_dst ? table_dst : ""));

	if(NULL == table_src || NULL == table_dst)
	{
		MYLOG_DEBUG(("invalid table name"));
		return -1;
	}

	std::ostringstream oss;
	oss << "select * into " << table_dst << " from " << table_src << " where 1=0" << std::ends;
	MYLOG_DEBUG(("%s", oss.str().c_str()));

	return this->query(oss.str().c_str());
}

/**
* @brief 获取最后一个插入的自增主键
*/
int32 mssqlwrapper::get_last_insert_id(uint64& last_id)
{
	MYLOG_DEBUG(("mssqlwrapper::get_last_insert_id"));

	if(0 != this->query("select @@identity"))
	{
		MYLOG_INFO(("query err"));
		return -1;
	}

	std::vector<std::vector<std::string> > vresult;
	if(0 != this->get_result(vresult))
	{
		MYLOG_INFO(("get result set err"));
		return -1;
	}

	if(vresult.size() != 1)
	{
		MYLOG_INFO(("result count is err"));
		return -1;
	}

	if(vresult[0].size() != 1)
	{
		MYLOG_INFO(("field count is err"));
		return -1;
	}

	myos::string2uint64(vresult[0][0], last_id);
	return 0;
}









