import { useNavigate } from 'react-router-dom'
import { useAuth } from '../../auth/hooks/useAuth'
import styles from './Header.module.css'

export const Header = () => {
  const navigate = useNavigate()
  const { user, logout } = useAuth()
  const handleLogout = async () => {
    try {
      logout()
      navigate('/')
    } catch (error) {
      console.error('Logout failed:', error)
    }
  }
  return (
    <header className={styles.header}>
      <p className={styles.title}>Chat-INFRA</p>
      {user && (
        <div>
          <p>{user.name}</p>
          <button className={styles.logout} onClick={handleLogout}>
            Logout
          </button>
        </div>
      )}
    </header>
  )
}