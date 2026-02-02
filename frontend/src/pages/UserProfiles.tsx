import { User } from "../lib/api";


export default function UserProfiles({ user }: { user: User }) {
  console.log('UserProfiles render', user);
  return (
    <div className="container">
      <div className="topbar">
        <div>
          <div style={{ fontSize: 20, fontWeight: 800 }}>User Profiles</div>
          <div className="muted">Manage user profiles and roles</div>
        </div>
        <div className="badge">{user.email} • {user.role}</div>
      </div>
    </div>
  )
}