import { Navigate, Route, Routes, Link, useLocation } from 'react-router-dom'
import LoginPage from './pages/Login'
import RegisterPage from './pages/Register'
import Dashboard from './pages/Dashboard'
import AssignWork from './pages/AssignWork'
import Sidebar from './pages/Sidbar'
import { useAuth } from './lib/useAuth'
import UserProfiles from './pages/UserProfiles'

export default function App() {
  const { user, loading, refresh, logout } = useAuth()
  const loc = useLocation()

  if (loading) {
    return (
      <div className="container">
        <div className="card" style={{ maxWidth: 520, margin: '40px auto' }}>
          Loading...
        </div>
      </div>
    )
  }

  const authed = !!user

  return (
    <div className="app-layout">
      {authed && user && <Sidebar user={user} />}
      <div className="main-content">
        {authed && (
          <div className="container" style={{ paddingBottom: 0 }}>
            <div className="topbar">
              <div className="row" style={{ gap: 10 }}>
                <div className="logo">📋 AssignWork</div>
              </div>
              <button style={{ width: 140 }} onClick={logout}>Logout</button>
            </div>
          </div>
        )}

        <Routes>
        <Route
          path="/login"
          element={
            authed ? <Navigate to="/" /> : <LoginPage onDone={refresh} />
          }
        />
        <Route
          path="/register"
          element={
            authed ? <Navigate to="/" /> : <RegisterPage onDone={refresh} />
          }
        />
        <Route
          path="/"
          element={
            authed && user ? (
              <Dashboard user={user} />
            ) : (
              <Navigate to="/login" state={{ from: loc.pathname }} />
            )
          }
        />
        <Route
          path="/assign-work"
          element={
            authed && user && user.role === 'admin' ? (
              <AssignWork user={user} />
            ) : (
              <Navigate to="/" />
            )
          }
        />
        <Route
          path="/user-profiles"
          element={
            authed && user ? (
              <UserProfiles user={user} />
            ) : (
              <Navigate to="/" />
            )
          }
        />
        <Route path="*" element={<Navigate to="/" />} />
        </Routes>
      </div>
    </div>
  )
}
