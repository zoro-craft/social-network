'use client';

import { useState, useRef, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '../app/context/AuthContext';
import { apiRequest } from '@/utils/api';
import Avatar from './Avatar';
import NotificationBell from './NotificationBell';
import ChatBell from './ChatBell';

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

  const goToGroups = () => {
    setOpen(false)
    router.push('/groups')
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
    try {
      await apiRequest('/logout', 'POST');
    } catch {
      // Clear the session locally even if the backend call fails.
    }
    setAuth({ authenticated: false, nickname: '' });
    setOpen(false);
    router.push('/login');
  };

  return (
    <header className="glass-navbar">
      <div className="glass-logo" onClick={() => router.push('/')}>
       Social-Network
      </div>

      <div className="glass-nav-actions">
        {auth.authenticated && (
          <>
            <button
              className="glass-nav-search"
              onClick={goToSearch}
              aria-label="Search"
            >
              Search
            </button>
              <button
              className="glass-nav-search"
              onClick={goToGroups}
              aria-label="groups"
            >
              groups
            </button>
            <NotificationBell />
            <ChatBell />
          </>
        )}
        
        <div className="glass-auth-buttons">
          {auth.authenticated ? (
            <div className="glass-avatar-wrap" ref={wrapRef}>
              <button
                className="glass-avatar-button"
                onClick={() => setOpen(!open)}
              >
                <Avatar src={auth.avatar} name={auth.first_name} size={28} />
                {auth.first_name}
              </button>

              {open && (
                <div className="glass-dropdown">
                  <button className="glass-dropdown-item" onClick={goToProfile}>
                    Profile
                  </button>
                  <button
                    className="glass-dropdown-item glass-nav-danger"
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
                className="glass-btn"
                onClick={() => router.push('/login')}
              >
                Login
              </button>
              <button
                className="glass-btn glass-btn-register"
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
