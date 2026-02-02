import { useEffect, useState } from 'react'
import { api, clearAuth, User } from './api'

export function useAuth() {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  async function refresh() {
    setLoading(true)
    try {
      const me = await api.auth.me()
      setUser(me.user)
    } catch {
      setUser(null)
    } finally {
      setLoading(false)
    }
  }

  function logout() {
    clearAuth()
    setUser(null)
  }

  useEffect(() => {
    refresh()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return { user, loading, refresh, logout }
}
