'use client';

import { useState, useRef, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import debounce from 'lodash/debounce';
import { API_BASE } from '../lib/api';
import Avatar from '../components/Avatar';
import styles from '../styles/search.module.css';
import { useAuth } from '@/app/context/AuthContext';

export default function SearchPage() {
  const router = useRouter();
  const [query, setQuery] = useState('');
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const { auth, loading: authLoading } = useAuth();

  useEffect(() => {
    if (!authLoading && !auth.authenticated) {
      router.push('/login');
    }
  }, [authLoading, auth.authenticated, router]);

  const searchUsers = async (value) => {
    if (!value.trim()) {
      setResults([]);
      setLoading(false);
      return;
    }

    try {
      const res = await fetch(
        `${API_BASE}/users/search?q=${encodeURIComponent(value)}`,
        {
          credentials: 'include',
        },
      );
      if (res.status === 401) {
        router.push('/login');
        return;
      }
      const data = await res.json();
      setResults(data || []);
    } catch {
      setResults([]);
    } finally {
      setLoading(false);
    }
  };

  // debounce so we only call the API 300ms after the user stops typing
  const debouncedSearch = useRef(debounce(searchUsers, 300)).current;

  const handleChange = (e) => {
    const value = e.target.value;
    setQuery(value);
    setLoading(true);
    debouncedSearch(value);
  };

  return (
    <div className={styles.container}>
      <input
        className={styles.input}
        type="text"
        placeholder="Search by name..."
        value={query}
        onChange={handleChange}
        autoFocus
      />

      {loading && <p className={styles.status}>Searching...</p>}

      {!loading && query && results.length === 0 && (
        <p className={styles.status}>No users found</p>
      )}

      <div className={styles.results}>
        {results.map((user) => (
          <div
            key={user.id}
            className={styles.resultItem}
            onClick={() => router.push(`/profile/${user.id}`)}
          >
            <Avatar src={user.avatar} name={user.full_name} size={40} />
            <p className={styles.name}>{user.full_name}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
