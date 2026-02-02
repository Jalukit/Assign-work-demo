import { api, TaskView, User } from '../lib/api'
import { useEffect, useState } from 'react'
import RegistrantsList from '../components/RegistrantsList'

export default function AssignWork({ user }: { user: User }) {
  const [tasks, setTasks] = useState<TaskView[]>([])
  const [busy, setBusy] = useState(false)

  const isAdmin = user.role === 'admin'
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [capacity, setCapacity] = useState(3)
  const [msg, setMsg] = useState('')
  const [error, setError] = useState('')

   async function createTask() {
    setMsg('')
    try {
      await api.tasks.create(title, description, capacity)
      setTitle('')
      setDescription('')
      setCapacity(3)
      setMsg('สร้างงานสำเร็จ')
      await load()
    } catch (e: any) {
      setMsg(e.message || 'create task failed')
    }
  }

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

  return (
    <>
      <div className='container'>
        <div className="topbar">
          <div>
            <div style={{ fontSize: 20, fontWeight: 800 }}>Assign Work</div>
            <div className="muted">เลือกลงทะเบียนงานได้ และดูคนที่ลงทะเบียน</div>
          </div>
          <div className="badge">{user.email} • {user.role}</div>
        </div>

        {msg && <div className={msg.includes('สำเร็จ') ? 'success' : 'error'} style={{ marginBottom: 12 }}>{msg}</div>}
        {error && <div className="error" style={{ marginBottom: 12 }}>{error}</div>}

        {isAdmin && (
          <div className="card" style={{ marginBottom: 14 }}>
            <h3 style={{ marginTop: 0 }}>Admin: Create Task</h3>
            <div className="row">
              <div>
                <label>Title</label>
                <input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="เช่น งานช่วยจัดของ" />
              </div>
              <div>
                <label>Capacity</label>
                <input type="number" value={capacity} onChange={(e) => setCapacity(Number(e.target.value))} min={1} />
              </div>
            </div>
            <div style={{ marginTop: 12 }}>
              <label>Description</label>
              <textarea value={description} onChange={(e) => setDescription(e.target.value)} placeholder="รายละเอียดงาน (optional)" rows={3} />
            </div>
            <div style={{ marginTop: 12 }}>
              <button onClick={createTask} disabled={!title.trim()}>Create</button>
            </div>
          </div>
        )}

        {/* Task List with Registrants */}
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
          {tasks.map((t) => (
            <div className="card" key={t.id}>
              <p className="taskTitle">{t.title}</p>
              {t.description && <div className="muted">{t.description}</div>}
              <div className="taskMeta">
                <span className="badge">capacity: {t.capacity}</span>
                <span className="badge">registered: {t.registeredCount}</span>
                <span className="badge">remaining: {t.remaining}</span>
              </div>

              <div className="hr" />
              <RegistrantsList
                registrants={t.registrants}
                isAdmin={isAdmin}
                onRemove={async (userId) => {
                  try {
                    await api.tasks.adminRemove(t.id, userId)
                    setMsg('ลบสำเร็จ')
                    await load()
                  } catch (e: any) {
                    setMsg(e.message || 'remove failed')
                  }
                }}
              />
            </div>
          ))}
        </div>
      </div>
    </>
  )
}