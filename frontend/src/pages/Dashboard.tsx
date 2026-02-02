import { useEffect, useMemo, useState } from 'react'
import { api, TaskView, User } from '../lib/api'

export default function Dashboard({ user }: { user: User }) {
  const [tasks, setTasks] = useState<TaskView[]>([])
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)

  // Admin create task
  
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [capacity, setCapacity] = useState(3)
  const [msg, setMsg] = useState('')

  async function load() {
    setBusy(true)
    setError('')
    try {
      const res = await api.tasks.list()
      setTasks(res.tasks)
    } catch (e: any) {
      setError(e.message || 'load failed')
    } finally {
      setBusy(false)
    }
  }

  useEffect(() => {
    load()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const myTaskIds = useMemo(() => {
    const set = new Set<string>()
    for (const t of tasks) {
      if (t.registrants.some((u) => u.id === user.id)) set.add(t.id)
    }
    return set
  }, [tasks, user.id])

  async function register(id: string) {
    setMsg('')
    try {
      await api.tasks.register(id)
      setMsg('ลงทะเบียนสำเร็จ')
      await load()
    } catch (e: any) {
      setMsg(e.message || 'register failed')
    }
  }

  async function unregister(id: string) {
    setMsg('')
    try {
      await api.tasks.unregister(id)
      setMsg('ยกเลิกสำเร็จ')
      await load()
    } catch (e: any) {
      setMsg(e.message || 'unregister failed')
    }
  }

  return (
    <div className="container">

      <div className="card" style={{ marginBottom: 14 }}>
        <div className="row" style={{ alignItems: 'end' }}>
          <div>
            <div className="muted">Task list</div>
            <div style={{ fontWeight: 700 }}>{busy ? 'Loading...' : `${tasks.length} tasks`}</div>
          </div>
          <div>
            <button onClick={load} disabled={busy}>Refresh</button>
          </div>
        </div>
      </div>

      <div className="tasks">
        {tasks.map((t) => {
          const isMine = myTaskIds.has(t.id)
          const full = t.remaining <= 0 || user.role === 'admin'
          return (
            <div className="card" key={t.id}>
              <p className="taskTitle">{t.title}</p>
              {t.description && <div className="muted">{t.description}</div>}
              <div className="taskMeta">
                <span className="badge">capacity: {t.capacity}</span>
                <span className="badge">registered: {t.registeredCount}</span>
                <span className="badge">remaining: {t.remaining}</span>
                {isMine && <span className="badge">you joined</span>}
              </div>

              <div className="row">
                <div>
                  {!isMine ? (
                    <button onClick={() => register(t.id)} disabled={full}>Join</button>
                  ) : (
                    <button onClick={() => unregister(t.id)}>Cancel</button>
                  )}
                </div>
              </div>

              <div className="hr" />
              <div style={{ fontWeight: 700, marginBottom: 8 }}>Registrants</div>
              {t.registrants.length === 0 ? (
                <div className="muted">ยังไม่มีคนลงทะเบียน</div>
              ) : (
                <ul style={{ margin: 0, paddingLeft: 18 }}>
                  {t.registrants
                    .sort((a, b) => a.email.localeCompare(b.email))
                    .map((u) => (
                      <li key={u.id}>
                        {u.email} <span className="muted">({u.role})</span>
                      </li>
                    ))}
                </ul>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}
