#include "AutoLock.h"


CAutoLock::CAutoLock(pthread_mutex_t* pLock, bool bLock)
    : m_pLock(pLock), m_bLock(bLock)
{
    if (bLock)
    {
        pthread_mutex_lock(m_pLock);
    }
}


CAutoLock::~CAutoLock(void)
{
    Unlock();
}


void CAutoLock::Unlock(void)
{
    if (m_bLock)
    {
        pthread_mutex_unlock(m_pLock);
        m_bLock = false;
    }
}
