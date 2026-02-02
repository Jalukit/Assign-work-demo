const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080'

export type User = { id: string; email: string; role: 'user' | 'admin' }

export type TaskView = {
  id: string
  title: string
  description: string
  capacity: number
  registeredCount: number
  remaining: number
  registrants: User[]
  createdAt: number
}

function getToken() {
  return localStorage.getItem('token') || ''
}

async function request<T>(path: string, opts: RequestInit = {}): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(opts.headers as any),
  }
  const token = getToken()
  if (token) headers.Authorization = `Bearer ${token}`

  const res = await fetch(`${API_BASE}${path}`, { ...opts, headers })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) {
    const msg = (data && (data.error || data.message)) || `HTTP ${res.status}`
    throw new Error(msg)
  }
  return data as T
}

export const api = {
  auth: {
    register: (email: string, password: string, role?: 'user' | 'admin') =>
      request<{ token: string; user: User }>('/api/auth/register', {
        method: 'POST',
        body: JSON.stringify({ email, password, role }),
      }),
    login: (email: string, password: string) =>
      request<{ token: string; user: User }>('/api/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password }),
      }),
    me: () => request<{ user: User }>('/api/me'),
  },
  tasks: {
    list: () => request<{ tasks: TaskView[] }>('/api/tasks'),
    create: (title: string, description: string, capacity: number) =>
      request<{ task: any }>('/api/tasks', {
        method: 'POST',
        body: JSON.stringify({ title, description, capacity }),
      }),
    register: (id: string) =>
      request<{ ok: boolean }>(`/api/tasks/${id}/register`, { method: 'POST' }),
    unregister: (id: string) =>
      request<{ ok: boolean }>(`/api/tasks/${id}/unregister`, { method: 'POST' }),
  },
}

export function saveAuth(token: string) {
  localStorage.setItem('token', token)
}

export function clearAuth() {
  localStorage.removeItem('token')
}
