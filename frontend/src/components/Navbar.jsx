'use client';

import { useState, useRef, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '../app/context/AuthContext';
import { API_BASE } from '@/utils/api';
import styles from '../app/styles/navbar.module.css';
import Avatar from './Avatar';

export default function Navbar() {
  const router = useRouter();
  const { auth, setAuth } = useAuth();
  const [open, setOpen] = useState(false);
  const wrapRef = useRef(null);

  const goToProfile = () => {
    setOpen(false);
    router.push(`/profile/${auth.id}`);
  };

  const goToSearch = () => {
    setOpen(false);
    router.push('/search');
  };

  const goToNotifications = () => {
    console.log('goToNotifications');
    setOpen(false);
    router.push('/notifications');
  };

  const gotoGrouplist = () => {
    setOpen(false);
    router.push('/groups/list');
  }

  useEffect(() => {
    function handleClickOutside(e) {
      if (wrapRef.current && !wrapRef.current.contains(e.target)) {
        setOpen(false);
      }
    }

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleLogout = async () => {
    await fetch(`${API_BASE}/logout`, {
      method: 'POST',
      credentials: 'include',
    });
    setAuth({ authenticated: false, nickname: '' });
    setOpen(false);
    router.push('/login');
  };

  return (
    <header className={styles.navbar}>
      <div className={styles.logo} onClick={() => router.push('/')}>
        01Forum
      </div>

      <div className={styles.navActions}>
        {auth.authenticated && (
          <>
            <button
              className={styles.searchButton}
              onClick={gotoGrouplist}
              aria-label="Groups"
            >
              Groups
            </button>
            <button
              className={styles.searchButton}
              onClick={goToSearch}
              aria-label="Search"
            >
              Search
            </button>
            <button
              className={styles.searchButton}
              onClick={goToNotifications}
              aria-label="Notifications"
            >
              🔔
            </button>
          </>
        )}
        
        <div className={styles.authButtons}>
          {auth.authenticated ? (
            <div className={styles.avatarWrap} ref={wrapRef}>
              <button
                className={styles.avatarButton}
                onClick={() => setOpen(!open)}
              >
                <Avatar src={auth.avatar} name={auth.first_name} size={28} />
                {auth.first_name}
              </button>

              {open && (
                <div className={styles.dropdown}>
                  <button className={styles.dropdownItem} onClick={goToProfile}>
                    Profile
                  </button>
                  <button
                    className={`${styles.dropdownItem} ${styles.danger}`}
                    onClick={handleLogout}
                  >
                    Logout
                  </button>
                </div>
              )}
            </div>
          ) : (
            <>
              <button
                className={`${styles.btn} ${styles.login}`}
                onClick={() => router.push('/login')}
              >
                Login
              </button>
              <button
                className={`${styles.btn} ${styles.register}`}
                onClick={() => router.push('/register')}
              >
                Register
              </button>
            </>
          )}
        </div>
      </div>
    </header>
  );
}
