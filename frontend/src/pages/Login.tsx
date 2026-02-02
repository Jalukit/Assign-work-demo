import { FormEvent, useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { api, saveAuth } from '../lib/api'

export default function LoginPage({ onDone }: { onDone: () => void }) {
  const nav = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [err, setErr] = useState('')
  const [busy, setBusy] = useState(false)

  async function submit(e: FormEvent) {
    e.preventDefault()
    setErr('')
    setBusy(true)
    try {
      const res = await api.auth.login(email, password)
      saveAuth(res.token)
      onDone()
      nav('/')
    } catch (e: any) {
      setErr(e.message || 'login failed')
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="container">
      <div className="card" style={{ maxWidth: 520, margin: '40px auto' }}>
        <h2 style={{ marginTop: 0, marginBottom: 20 }}>Login</h2>
        <form onSubmit={submit} className="row">
          <div>
            <label>Email</label>
            <input value={email} onChange={(e) => setEmail(e.target.value)} placeholder="you@example.com" />
          </div>
          <div>
            <label>Password</label>
            <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="ขั้นต่ำ 6 ตัว" />
          </div>
          {err && <div className="error">{err}</div>}
          <div>
            <button disabled={busy}>{busy ? '...' : 'Login'}</button>
          </div>
        </form>

        <div className="hr" />
        <div className="muted">
          ยังไม่มีบัญชี? <Link to="/register">ไปสมัคร</Link>
        </div>
      </div>
    </div>
  )
}
