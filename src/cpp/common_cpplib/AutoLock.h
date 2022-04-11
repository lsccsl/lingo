#pragma once

#include <pthread.h>

class CAutoLock
{
    pthread_mutex_t* m_pLock;
    bool m_bLock;

public:
    CAutoLock(pthread_mutex_t* pLock, bool bLock = true);
    ~CAutoLock(void);

    void Unlock(void);
};


