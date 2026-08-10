'use client';

import { useState, useRef, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '../context/AuthContext';
import { API_BASE } from '../lib/api';
import Image from 'next/image';
import styles from '../styles/navbar.module.css';

export default function Navbar() {
  const router = useRouter();
  const { auth, setAuth } = useAuth();
  const [open, setOpen] = useState(false);
  const wrapRef = useRef(null);

  const goToProfile = () => {
    setOpen(false);
    router.push(`/profile/${auth.id}`);
  };

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
      <div className={styles.authButtons}>
        {auth.authenticated ? (
          <div className={styles.avatarWrap} ref={wrapRef}>
            <button
              className={styles.avatarButton}
              onClick={() => setOpen(!open)}
            >
              {auth.avatar ? (
                <Image
                  className={styles.avatarImg}
                  src={`${API_BASE}${auth.avatar}`}
                  alt=""
                  width={28}
                  height={28}
                  unoptimized
                />
              ) : (
                <span className={styles.avatarFallback}>
                  {auth.first_name?.charAt(0).toUpperCase()}
                </span>
              )}
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
    </header>
  );
}
