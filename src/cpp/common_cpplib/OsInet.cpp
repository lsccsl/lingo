
#include "OsInet.h"

#ifndef WIN32
#include <stdio.h>     
#include <sys/types.h>     
#include <sys/param.h>     

#include <sys/ioctl.h>     
#include <sys/socket.h>     
#include <net/if.h>     
#include <netinet/in.h>     
#include <net/if_arp.h>     
#include <arpa/inet.h>
#include <unistd.h>
#define MAXINTERFACES      16    
#else
#include <winsock2.h>
#endif


/**
 * @brief ��ȡ����ip
 */
int32 OsInet::GetLocalIP(uint32& ip)
{
#ifdef WIN32
	char hostname[100];
	gethostname(hostname, 100);
	struct hostent *hptr = NULL;
	if ((hptr = gethostbyname(hostname)) == NULL)
	{
		return -1;
	}
	struct in_addr **ppaddr;
	ppaddr = (struct in_addr**)hptr->h_addr_list;

	ip = ntohl((*ppaddr)->s_addr);

	return 0;
#else
	std::vector<uint32> ip_list;
	OsInet::GetLocalIPs(ip_list);
	uint32 local_ip = (127 << 24) | (0 << 16) | (0 << 8) | 1;

	for(uint32 i = 0; i < ip_list.size(); i ++)
	{
		if(local_ip == ip_list[i])
			continue;

		ip = ip_list[i];
		break;
	}
	return 0;
#endif
}

void OsInet::GetLocalIPs(std::vector<uint32>& ip_list)
{
#ifndef WIN32
	int fd, intrface, retn = 0;     
	struct ifreq buf[MAXINTERFACES];     
	struct arpreq arp;     
	struct ifconf ifc;     

	/*
	if (!(ioctl(fd, SIOCGIFHWADDR, (char *) &buf[intrface])))
	*/

	if((fd = socket(AF_INET,SOCK_DGRAM,0)) >= 0)
	{
		ifc.ifc_len = sizeof(buf);     
		ifc.ifc_buf = (caddr_t) buf;     
		if (!ioctl(fd, SIOCGIFCONF, (char *) &ifc))
		{
			intrface = ifc.ifc_len / sizeof(struct ifreq);   
			printf("interface num is intrface=%d\n\n\n", intrface);     
			while (intrface-- > 0)
			{
				printf("====================interface:%d====================\n",intrface);
				printf("net device %s\n", buf[intrface].ifr_name);

				/*Jugde whether the net card status is promisc    */     
				if (!(ioctl(fd, SIOCGIFFLAGS, (char *) &buf[intrface])))
				{
					if (buf[intrface].ifr_flags & IFF_PROMISC)
					{
						puts("the interface is PROMISC");     
						retn++;
					}
				}
				else
				{
					char str[256];     
					sprintf(str, "cpm: ioctl device %s",
						buf[intrface].ifr_name);     
					perror(str);
				}     

				/*Jugde whether the net card status is up                  */
				int isup = 1;
				if (buf[intrface].ifr_flags & IFF_UP)
				{
					isup = 1;
					puts("the interface status is UP");
				}
				else
				{
					isup = 0;
					puts("the interface status is DOWN");
				}     
				/*Get IP of the net card */     
				if (!(ioctl(fd, SIOCGIFADDR, (char *) &buf[intrface])))
				{
					puts("IP address is:");     
					puts(inet_ntoa(((struct sockaddr_in *)
						(&buf[intrface].ifr_addr))->sin_addr));     
					puts("");

					if(isup)
						ip_list.push_back(ntohl((((struct sockaddr_in *)(&buf[intrface].ifr_addr))->sin_addr).s_addr));

					//puts (buf[intrface].ifr_addr.sa_data);     
				}
				else
				{
					char str[256];     
					sprintf(str, "cpm: ioctl device %s",
						buf[intrface].ifr_name);     
					perror(str);
				}

				/*Get HW ADDRESS of the net card */     
				if (!(ioctl(fd, SIOCGIFHWADDR, (char *) &buf[intrface])))
				{
					puts("HW address is:");     
					printf("%02x:%02x:%02x:%02x:%02x:%02x\n",
						(unsigned char) buf[intrface].ifr_hwaddr.sa_data[0],   
						(unsigned char) buf[intrface].ifr_hwaddr.sa_data[1],
						(unsigned char) buf[intrface].ifr_hwaddr.sa_data[2],   
						(unsigned char) buf[intrface].ifr_hwaddr.sa_data[3],   
						(unsigned char) buf[intrface].ifr_hwaddr.sa_data[4],   
						(unsigned char) buf[intrface].ifr_hwaddr.sa_data[5]);
					puts("");     
					puts("");
				}
				else
				{
					char str[256];     
					sprintf(str, "cpm: ioctl device %s", buf[intrface].ifr_name);     
					perror(str);
				}
			}
		}
		else
			perror("cpm: ioctl");
	}
	else
		perror("cpm: socket");     

	close(fd);     
	return;
#endif
}
