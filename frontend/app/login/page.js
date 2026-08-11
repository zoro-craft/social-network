'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '../context/AuthContext';
import { API_BASE } from '../lib/api';
import { useToast } from '@/app/context/ToastContext';
import styles from '../styles/auth.module.css';

export default function LoginPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const { refreshAuth } = useAuth();
  const showToast = useToast();
  const [form, setForm] = useState({
    login: '',
    password: '',
  });

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      const body = new URLSearchParams(form);

      const res = await fetch(`${API_BASE}/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        credentials: 'include',
        body,
      });

      if (!res.ok) {
        const data = await res.json();
        showToast(data.message || 'Login failed');
        setLoading(false);
        return;
      }
      await refreshAuth();
      showToast('login successfully', false);
      router.push('/');
    } catch (err) {
      showToast('Could not reach the server');
      setLoading(false);
    }
  };

  return (
    <div className={styles.wrapper}>
      <div className={styles.card}>
        <h1 className={styles.title}>Welcome back</h1>
        <p className={styles.subtitle}>Log in to continue to 01Forum</p>

        <form className={styles.form} onSubmit={handleSubmit}>
          <input
            className={styles.input}
            name="login"
            type="text"
            placeholder="Email or nickname"
            value={form.login}
            onChange={handleChange}
            required
          />
          <input
            className={styles.input}
            name="password"
            type="password"
            placeholder="Password"
            value={form.password}
            onChange={handleChange}
            required
          />
          <button className={styles.button} type="submit" disabled={loading}>
            {loading ? 'Logging in...' : 'Login'}
          </button>
        </form>

        <p className={styles.footerText}>
          Don&apos;t have an account? <a href="/register">Register</a>
        </p>
      </div>
    </div>
  );
}
