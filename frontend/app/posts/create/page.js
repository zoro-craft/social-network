'use client';

import { useAuth } from '@/app/context/AuthContext';
import { useToast } from '@/app/context/ToastContext';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { API_BASE } from '../../lib/api';
import Avatar from '../../components/Avatar';
import styles from '../../styles/auth.module.css';

export default function CreatePostPage() {
  const router = useRouter();
  const { auth, loading: authLoading } = useAuth();
  const showToast = useToast();

  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [privacy, setPrivacy] = useState('public');
  const [imageFile, setImageFile] = useState(null);
  const [followers, setFollowers] = useState([]);
  const [selectedViewers, setSelectedViewers] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!authLoading && !auth.authenticated) {
      router.push('/login');
    }
  }, [authLoading, auth.authenticated, router]);

  useEffect(() => {
    if (privacy !== 'private' || !auth.id) return;

    fetch(`${API_BASE}/profile/${auth.id}/followers`, {
      credentials: 'include',
    })
      .then((res) => res.json())
      .then((data) => setFollowers(data || []));
  }, [privacy, auth.id]);

  const toggleViewer = (id) => {
    setSelectedViewers((prev) =>
      prev.includes(id) ? prev.filter((v) => v !== id) : [...prev, id],
    );
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!title.trim() || !content.trim()) {
      showToast('Title and content are required');
      return;
    }
    if (privacy === 'private' && selectedViewers.length === 0) {
      showToast('Select at least one follower');
      return;
    }

    setLoading(true);

    try {
      const form = new FormData();
      form.append('title', title);
      form.append('content', content);
      form.append('privacy', privacy);
      if (imageFile) form.append('image', imageFile);
      selectedViewers.forEach((id) => form.append('viewer_ids', id));

      const res = await fetch(`${API_BASE}/posts/create`, {
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
        showToast(data.message || 'Failed to create post');
        setLoading(false);
        return;
      }

      showToast('Post created successfully!', false);
      router.push('/');
    } catch {
      showToast('Could not reach the server');
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
            <option value="almost_private">Almost private (followers)</option>
            <option value="private">Private (chosen followers)</option>
          </select>

          {privacy === 'private' && (
            <div className={styles.viewerBox}>
              <p className={styles.viewerLabel}>
                Choose who can see this post:
              </p>
              {followers.length === 0 && (
                <p className={styles.viewerEmpty}>No followers yet</p>
              )}
              <div className={styles.viewerGrid}>
                {followers.map((f) => (
                  <div
                    key={f.follower_id}
                    className={`${styles.viewerItem} ${selectedViewers.includes(f.follower_id) ? styles.viewerSelected : ''}`}
                    onClick={() => toggleViewer(f.follower_id)}
                  >
                    <Avatar src={f.avatar} name={f.first_name} size={40} />
                    <span>{f.first_name}</span>
                  </div>
                ))}
              </div>
            </div>
          )}

          <label className={styles.fileLabel}>
            Image (optional)
            <input
              className={styles.fileInput}
              type="file"
              accept="image/jpeg,image/png,image/gif"
              onChange={(e) => setImageFile(e.target.files[0])}
            />
          </label>

          <button className={styles.button} type="submit" disabled={loading}>
            {loading ? 'Posting...' : 'Post'}
          </button>
        </form>
      </div>
    </div>
  );
}
