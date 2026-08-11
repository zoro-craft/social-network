'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { API_BASE } from '../lib/api';
import { useToast } from '@/app/context/ToastContext';

import styles from '../styles/auth.module.css';

export default function RegisterPage() {
  const router = useRouter();
  const showToast = useToast();
  const [loading, setLoading] = useState(false);
  const [avatarFile, setAvatarFile] = useState(null);

  const [form, setForm] = useState({
    first_name: '',
    last_name: '',
    email: '',
    password: '',
    confirm_password: '',
    date_of_birth: '',
    nickname: '',
    about_me: '',
  });

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      let avatarPath = '';

      if (avatarFile) {
        const avatarForm = new FormData();
        avatarForm.append('avatar', avatarFile);

        const uploadRes = await fetch(`${API_BASE}/upload/avatar`, {
          method: 'POST',
          body: avatarForm,
        });

        const uploadData = await uploadRes.json();

        if (!uploadRes.ok) {
          showToast(uploadData.message || 'Avatar upload failed');
          setLoading(false);
          return;
        }

        avatarPath = uploadData.path;
      }

      const body = new URLSearchParams({ ...form, avatar: avatarPath });

      const res = await fetch(`${API_BASE}/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body,
      });

      const data = await res.json();

      if (!res.ok) {
        showToast(data.message || 'Registration failed');
        setLoading(false);
        return;
      }
      showToast('Account created successfully', false);
      router.push('/login');
    } catch (err) {
      showToast('Could not reach the server');
      setLoading(false);
    }
  };

  return (
    <div className={styles.wrapper}>
      <div className={styles.card}>
        <h1 className={styles.title}>Create account</h1>
        <p className={styles.subtitle}>Join 01Forum in a few seconds</p>

        <form className={styles.form} onSubmit={handleSubmit}>
          <div className={styles.row}>
            <input
              className={styles.input}
              name="first_name"
              placeholder="First name"
              value={form.first_name}
              onChange={handleChange}
              required
            />
            <input
              className={styles.input}
              name="last_name"
              placeholder="Last name"
              value={form.last_name}
              onChange={handleChange}
              required
            />
          </div>

          <input
            className={styles.input}
            name="email"
            type="email"
            placeholder="Email"
            value={form.email}
            onChange={handleChange}
            required
          />

          <div className={styles.row}>
            <input
              className={styles.input}
              name="password"
              type="password"
              placeholder="Password"
              value={form.password}
              onChange={handleChange}
              required
            />
            <input
              className={styles.input}
              name="confirm_password"
              type="password"
              placeholder="Confirm password"
              value={form.confirm_password}
              onChange={handleChange}
              required
            />
          </div>

          <input
            className={styles.input}
            name="date_of_birth"
            type="date"
            value={form.date_of_birth}
            onChange={handleChange}
            required
          />

          <input
            className={styles.input}
            name="nickname"
            placeholder="Nickname (optional)"
            value={form.nickname}
            onChange={handleChange}
          />

          <textarea
            className={styles.textarea}
            name="about_me"
            placeholder="About me (optional)"
            value={form.about_me}
            onChange={handleChange}
          />

          <label className={styles.fileLabel}>
            Avatar (optional)
            <input
              className={styles.fileInput}
              type="file"
              accept="image/jpeg,image/png,image/gif"
              onChange={(e) => setAvatarFile(e.target.files[0])}
            />
          </label>
          <button className={styles.button} type="submit" disabled={loading}>
            {loading ? 'Registering...' : 'Register'}
          </button>
        </form>

        <p className={styles.footerText}>
          Already have an account? <a href="/login">Login</a>
        </p>
      </div>
    </div>
  );
}
