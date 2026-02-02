import { Link } from 'react-router-dom'
import { User } from '../lib/api'

export default function Sidebar({ user }: { user: User }) {
  const isAdmin = user.role === 'admin'

  return (
    <div className="sidebar">
      <h2 >Menu</h2>
      <ul>
        <li><Link to="/">Dashboard</Link></li>
        {isAdmin && <li><Link to="/assign-work">Assign Work</Link></li>}
        <li><Link to="/user-profiles">User Profiles</Link></li>
      </ul>
    </div>
  )
}