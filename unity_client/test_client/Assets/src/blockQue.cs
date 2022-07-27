
using System;
using System.Collections.Generic;
using System.Threading;

public class BlockQueue<T>
{
    private readonly Queue<T> _queue;
    private readonly int _capacity;
    private bool _closing;

    public int count
    {
        get { return _queue.Count; }
    }

    public BlockQueue(int capacity)
    {
        this._capacity = capacity;
        this._queue = new Queue<T>(capacity);
    }

    public void Enqueue(T item)
    {
        lock (_queue)
        {
            while (_queue.Count >= _capacity)
            {
                Monitor.Wait(_queue);
            }
            _queue.Enqueue(item);
            if (_queue.Count == 1)
            {
                Monitor.PulseAll(_queue);
            }
        }
    }

    public T Dequeue()
    {
        lock (_queue)
        {
            while (_queue.Count == 0)
            {
                Monitor.Wait(_queue);
            }
            T item = _queue.Dequeue();
            if (_queue.Count == _capacity - 1)
            {
                Monitor.PulseAll(_queue);
            }
            return item;
        }
    }


    public void Close()
    {
        lock (_queue)
        {
            _closing = true;
            Monitor.PulseAll(_queue);
        }
    }

    public bool TryDequeue(out T value)
    {
        lock (_queue)
        {
            while (_queue.Count == 0)
            {
                if (_closing)
                {
                    value = default(T);
                    return false;
                }
                Monitor.Wait(_queue);
            }
            value = _queue.Dequeue();
            if (_queue.Count == _capacity - 1)
            {
                Monitor.PulseAll(_queue);
            }
            return true;
        }
    }
}
