'use client';

import { useState, useRef, useEffect } from 'react';
import { usePathname, useRouter } from 'next/navigation';
import { useAuth } from '@/app/context/AuthContext';
import { apiRequest } from '@/utils/api';
import Avatar from './Avatar';
import Icon from './Icon';
import NotificationBell from './NotificationBell';
import ChatBell from './ChatBell';

// The primary destinations, defined once so the labels, icons and active-state
// rules can never drift apart between breakpoints.
const NAV_LINKS = [
  { href: '/', label: 'Feed', icon: 'home', also: ['/posts'] },
  { href: '/reels', label: 'Reels', icon: 'film' },
  { href: '/groups', label: 'Groups', icon: 'users' },
  { href: '/search', label: 'Search', icon: 'search' },
];

function isActive(pathname, link) {
  if (link.href === '/') {
    return (
      pathname === '/' ||
      (link.also || []).some((prefix) => pathname.startsWith(prefix))
    );
  }

  return pathname === link.href || pathname.startsWith(`${link.href}/`);
}

export default function Navbar() {
  const router = useRouter();
  const pathname = usePathname();
  const { auth, setAuth } = useAuth();
  const [open, setOpen] = useState(false);
  const wrapRef = useRef(null);

  const go = (href) => {
    setOpen(false);
    router.push(href);
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
      <button type="button" className="glass-logo" onClick={() => go('/')}>
        Social-Network
      </button>

      {auth.authenticated && (
        <nav className="glass-nav" aria-label="Primary">
          {NAV_LINKS.map((link) => {
            const active = isActive(pathname, link);

            return (
              <button
                key={link.href}
                type="button"
                className={`glass-nav-link ${active ? 'glass-nav-link-active' : ''}`}
                onClick={() => go(link.href)}
                aria-current={active ? 'page' : undefined}
              >
                <span className="glass-nav-link-icon" aria-hidden="true">
                  <Icon name={link.icon} />
                </span>
                <span className="glass-nav-link-label">{link.label}</span>
              </button>
            );
          })}
        </nav>
      )}

      <div className="glass-nav-actions">
        {auth.authenticated ? (
          <>
            <NotificationBell />
            <ChatBell />

            <div className="glass-avatar-wrap" ref={wrapRef}>
              <button
                className="glass-avatar-button"
                onClick={() => setOpen(!open)}
                aria-expanded={open}
                aria-haspopup="menu"
              >
                <Avatar src={auth.avatar} name={auth.first_name} size={28} />
                <span className="glass-avatar-name">{auth.first_name}</span>
              </button>

              {open && (
                <div className="glass-dropdown" role="menu">
                  <button
                    className="glass-dropdown-item"
                    onClick={() => go(`/profile/${auth.id}`)}
                  >
                    Profile
                  </button>
                  <button
                    className="glass-dropdown-item"
                    onClick={() => go('/notifications')}
                  >
                    Notifications
                  </button>
                  <button className="glass-dropdown-item" onClick={() => go('/chat')}>
                    Chat
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
          </>
        ) : (
          <>
            <button className="glass-btn" onClick={() => go('/login')}>
              Login
            </button>
            <button
              className="glass-btn glass-btn-register"
              onClick={() => go('/register')}
            >
              Register
            </button>
          </>
        )}
      </div>
    </header>
  );
}
