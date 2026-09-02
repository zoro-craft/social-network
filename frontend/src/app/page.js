import styles from './styles/auth.module.css';
import Link from 'next/link';

export default function HomePage() {
  return (
    <div className={styles.wrapper}>
      <div className={styles.card} style={{ textAlign: 'center' }}>
        <h1 className={styles.title}>Welcome to Social Network</h1>
        <p className={styles.subtitle}>Authentication system is ready!</p>
        
        <div style={{ display: 'flex', gap: 12, justifyContent: 'center' }}>
          <Link href="/login">
            <button className={styles.button} style={{ marginTop: 0 }}>
              Go to Login
            </button>
          </Link>
          <Link href="/register">
            <button className={styles.button} style={{ marginTop: 0, background: 'var(--accent)' }}>
              Create Account
            </button>
          </Link>
        </div>
      </div>
    </div>
  );
}