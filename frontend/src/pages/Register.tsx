import { FormEvent, useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { api, saveAuth } from '../lib/api'

export default function RegisterPage({ onDone }: { onDone: () => void }) {
  const nav = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [role, setRole] = useState<'user' | 'admin'>('user')
  const [err, setErr] = useState('')
  const [busy, setBusy] = useState(false)

  async function submit(e: FormEvent) {
    e.preventDefault()
    setErr('')
    setBusy(true)
    try {
      const res = await api.auth.register(email, password, role)
      saveAuth(res.token)
      onDone()
      nav('/')
    } catch (e: any) {
      setErr(e.message || 'register failed')
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="container">
      <div className="card" style={{ maxWidth: 520, margin: '40px auto' }}>
        <h2 style={{ marginTop: 0 }}>Register</h2>
        <p className="muted">
          เดโม่นี้ให้เลือก role ได้เลย (admin จะสร้างงานได้) — ของจริงควรให้สร้าง admin เฉพาะในระบบหลังบ้าน
        </p>

        <form onSubmit={submit} className="row">
          <div>
            <label>Email</label>
            <input value={email} onChange={(e) => setEmail(e.target.value)} placeholder="you@example.com" />
          </div>
          <div>
            <label>Password</label>
            <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="ขั้นต่ำ 6 ตัว" />
          </div>
          <div>
            <label>Role</label>
            <select value={role} onChange={(e) => setRole(e.target.value as any)}>
              <option value="user">user</option>
              <option value="admin">admin</option>
            </select>
          </div>

          {err && <div className="error">{err}</div>}
          <div>
            <button disabled={busy}>{busy ? '...' : 'Create account'}</button>
          </div>
        </form>

        <div className="hr" />
        <div className="muted">
          มีบัญชีแล้ว? <Link to="/login">ไปล็อกอิน</Link>
        </div>
      </div>
    </div>
  )
}
