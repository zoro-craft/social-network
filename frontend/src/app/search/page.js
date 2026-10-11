'use client';

import { useState, useRef, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import debounce from 'lodash/debounce';
import { API_BASE } from '@/utils/api';
import Avatar from '@/app/components/Avatar';
import { useRequireAuth } from '@/hooks/useRequireAuth';

export default function SearchPage() {
  const router = useRouter();
  const [query, setQuery] = useState('');
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const { auth, loading: authLoading } = useRequireAuth();

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
      if (!res.ok) {
        setResults([]);
        return;
      }
      const data = await res.json();
      setResults(Array.isArray(data) ? data : []);
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
    <div className="glass-search-page">
      <input
        className="glass-search-input"
        type="text"
        placeholder="Search by name..."
        value={query}
        onChange={handleChange}
        autoFocus
      />

      {loading && <p className="glass-search-status">Searching...</p>}

      {!loading && query && results.length === 0 && (
        <p className="glass-search-status">No users found</p>
      )}

      <div className="glass-search-results">
        {results.map((user) => (
          <div
            key={user.id}
            className="glass-search-item"
            onClick={() => router.push(`/profile/${user.id}`)}
          >
            <Avatar src={user.avatar} name={user.full_name} size={40} />
            <p className="glass-search-name">{user.full_name}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
