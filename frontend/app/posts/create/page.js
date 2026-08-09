'use client';

import { useAuth } from '@/app/context/AuthContext';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import styles from '../../styles/auth.module.css';

export default function CreatePostPage() {
  const router = useRouter();
  const { auth, loading: authLoading } = useAuth();

  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [privacy, setPrivacy] = useState('public');
  const [imageFile, setImageFile] = useState(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!authLoading && !auth.authenticated) {
      router.push('/login');
    }
  }, [authLoading, auth.authenticated, router]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    if (!title.trim() || !content.trim()) {
      setError('Title and content are required');
      return;
    }

    setLoading(true);

    try {
      const form = new FormData();
      form.append('title', title);
      form.append('content', content);
      form.append('privacy', privacy);
      if (imageFile) {
        form.append('image', imageFile);
      }

      const res = await fetch('http://localhost:8080/posts/create', {
        method: 'POST',
        credentials: 'include',
        body: form,
      });

      if (res.status === 401) {
        router.push('/login');
        return;
      }

      const data = await res.json();

      if (!res.ok) {
        setError(data.message || 'Failed to create post');
        setLoading(false);
        return;
      }

      router.push('/');
    } catch {
      setError('Could not reach the server');
      setLoading(false);
    }
  };

  if (authLoading || !auth.authenticated) {
    return null;
  }

  return (
    <div className={styles.wrapper}>
      <div className={styles.card}>
        <h1 className={styles.title}>Create Post</h1>

        <form className={styles.form} onSubmit={handleSubmit}>
          <input
            className={styles.input}
            type="text"
            placeholder="Title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
          />

          <textarea
            className={styles.textarea}
            placeholder="Write something..."
            value={content}
            onChange={(e) => setContent(e.target.value)}
            rows={4}
            required
          />

          <select
            className={styles.input}
            value={privacy}
            onChange={(e) => setPrivacy(e.target.value)}
          >
            <option value="public">Public</option>
            <option value="almost_private">Almost private</option>
            <option value="private">Private</option>
          </select>

          <label className={styles.fileLabel}>
            Image 
            <input
              className={styles.fileInput}
              type="file"
              accept="image/jpeg,image/png,image/gif"
              onChange={(e) => setImageFile(e.target.files[0])}
            />
          </label>

          {error && <p className={styles.error}>{error}</p>}

          <button className={styles.button} type="submit" disabled={loading}>
            {loading ? 'Posting...' : 'Post'}
          </button>
        </form>
      </div>
    </div>
  );
}
