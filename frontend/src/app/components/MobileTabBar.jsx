'use client';

import { usePathname, useRouter } from 'next/navigation';
import { useAuth } from '@/app/context/AuthContext';
import Icon from './Icon';

// Fixed bottom navigation shown only on phones (<768px); the sidebars are
// hidden there, so this is the only persistent way to switch sections.
export default function MobileTabBar() {
  const router = useRouter();
  const pathname = usePathname();
  const { auth } = useAuth();

  if (!auth.authenticated) return null;

  const tabs = [
    { href: '/', label: 'Feed', icon: 'home' },
    { href: '/reels', label: 'Reels', icon: 'film' },
    { href: '/groups', label: 'Groups', icon: 'users' },
    { href: '/search', label: 'Search', icon: 'search' },
    { href: `/profile/${auth.id}`, label: 'Profile', icon: 'user' },
  ];

  const isActive = (href) => {
    if (href === '/') {
      return pathname === '/' || pathname.startsWith('/posts');
    }

    return pathname === href || pathname.startsWith(`${href}/`);
  };

  return (
    <nav className="glass-tabbar" aria-label="Primary">
      {tabs.map((tab) => {
        const active = isActive(tab.href);

        return (
          <button
            key={tab.label}
            type="button"
            className={`glass-tabbar-link ${active ? 'glass-tabbar-link-active' : ''}`}
            onClick={() => router.push(tab.href)}
            aria-current={active ? 'page' : undefined}
          >
            <Icon name={tab.icon} size={20} />
            <span className="glass-tabbar-label">{tab.label}</span>
          </button>
        );
      })}
    </nav>
  );
}
