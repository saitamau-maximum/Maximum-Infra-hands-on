import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../../auth/hooks/useAuth'
import styles from './Header.module.css'
import apiClient from '../../utils/apiClient'
import { Icon } from '../../ui/icon'

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
      <Link to='/' className={styles.title}>Chat-INFRA</Link>
      {user && (
        <div className={styles.user}>
          <p className={styles.user_name}>{user.name}</p>
          <Icon src={`${apiClient.baseUrl}/api/user/icon/${user.id}`} alt='user icon' size={24}/>
          <button className={styles.logout_button} onClick={handleLogout}>
            Logout
          </button>
        </div>
      )}
    </header>
  )
}