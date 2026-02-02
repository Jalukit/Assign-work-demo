import { User } from '../lib/api'

type Props = {
  registrants: User[]
  isAdmin?: boolean
  onRemove?: (userId: string) => void
}

export default function RegistrantsList({ registrants, isAdmin, onRemove }: Props) {
  return (
    <div>
      <div style={{ fontWeight: 700, marginBottom: 8 }}>Registrants</div>
      {registrants.length === 0 ? (
        <div className="muted">ยังไม่มีคนลงทะเบียน</div>
      ) : (
        <ul style={{ margin: 0, paddingLeft: 18 }}>
          {registrants
            .sort((a, b) => a.email.localeCompare(b.email))
            .map((u) => (
              <li key={u.id} style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <span>
                  {u.email} <span className="muted">({u.role})</span>
                </span>
                {isAdmin && onRemove && (
                  <button
                    onClick={() => onRemove(u.id)}
                    style={{
                      width: 24,
                      height: 24,
                      padding: 0,
                      marginLeft: 'auto',
                      fontSize: 14,
                      background: 'rgba(255, 100, 100, 0.3)',
                      border: '1px solid rgba(255, 100, 100, 0.5)',
                      borderRadius: 4,
                      cursor: 'pointer',
                    }}
                    className='remove-button'
                    title={`Remove ${u.email}`}
                  >
                    ✕
                  </button>
                )}
              </li>
            ))}
        </ul>
      )}
    </div>
  )
}
