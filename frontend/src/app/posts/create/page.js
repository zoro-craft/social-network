'use client';

import { useAuth } from '@/app/context/AuthContext';
import { useToast } from '@/app/context/ToastContext';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { API_BASE } from '@/utils/api';
import Avatar from '@/components/Avatar';
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

    const fetchFollowers = async () => {
      try {
        const res = await fetch(`${API_BASE}/profile/${auth.id}/followers`, {
          credentials: 'include',
        });

        if (!res.ok) {
          throw new Error('Failed to fetch followers');
        }

        const data = await res.json();
        setFollowers(data || []);
      } catch (error) {
        console.error('Error fetching followers:', error);
        setFollowers([]);
      }
    };

    fetchFollowers();
  }, [privacy, auth.id]);

  const toggleViewer = (id) => {
    setSelectedViewers((prev) =>
      prev.includes(id) ? prev.filter((v) => v !== id) : [...prev, id],
    );
  };

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    setImageFile(file || null);
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
      router.push('/posts');
    } catch {
      showToast('Could not reach the server');
      setLoading(false);
    }
  };

  if (authLoading || !auth.authenticated) {
    return null;
  }

  return (
    <div className="glass-page">
      <div className="glass-container" style={{ maxWidth: '600px' }}>
        <h2>Create Post</h2>
        <p style={{ textAlign: 'center', color: 'var(--text-muted)', marginBottom: '24px' }}>
          Share something with your network
        </p>

        <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
          <div className="glass-input-group">
            <input
              className="glass-input"
              type="text"
              placeholder="Title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
            />
          </div>

          <div className="glass-input-group">
            <textarea
              className="glass-textarea"
              placeholder="Write something..."
              value={content}
              onChange={(e) => setContent(e.target.value)}
              rows={4}
              required
            />
          </div>

          <div className="glass-input-group">
            <select
              className="glass-input"
              value={privacy}
              onChange={(e) => setPrivacy(e.target.value)}
              style={{ appearance: 'none' }}
            >
              <option value="public" style={{ background: '#0f172a' }}>Public</option>
              <option value="almost_private" style={{ background: '#0f172a' }}>Almost private (followers)</option>
              <option value="private" style={{ background: '#0f172a' }}>Private (chosen followers)</option>
            </select>
          </div>

          {privacy === 'private' && (
            <div className="glass-input-group" style={{ background: 'rgba(255, 255, 255, 0.02)', padding: '16px', borderRadius: '12px', border: '1px solid var(--border)' }}>
              <label style={{ display: 'block', marginBottom: '12px', color: 'var(--text-secondary)' }}>
                Choose who can see this post:
              </label>
              {followers.length === 0 && (
                <p style={{ color: 'var(--text-muted)', fontSize: '13px' }}>No followers yet</p>
              )}
              <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(120px, 1fr))', gap: '12px' }}>
                {followers.map((f) => (
                  <div
                    key={f.follower_id}
                    onClick={() => toggleViewer(f.follower_id)}
                    style={{
                      display: 'flex',
                      alignItems: 'center',
                      gap: '8px',
                      padding: '8px',
                      borderRadius: '8px',
                      cursor: 'pointer',
                      border: selectedViewers.includes(f.follower_id) ? '1px solid var(--accent-border)' : '1px solid var(--border)',
                      background: selectedViewers.includes(f.follower_id) ? 'var(--accent-soft)' : 'var(--bg-elevated)',
                      transition: 'all 0.2s ease',
                      fontSize: '13px',
                      color: selectedViewers.includes(f.follower_id) ? '#fff' : 'var(--text-secondary)',
                    }}
                  >
                    <Avatar src={f.avatar} name={f.first_name} size={32} />
                    <span style={{ overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{f.first_name}</span>
                  </div>
                ))}
              </div>
            </div>
          )}

          <label className="glass-file-label" style={{ marginTop: '8px' }}>
            Image
            <input
              className="glass-file-input"
              type="file"
              accept="image/jpeg,image/png,image/gif"
              onChange={handleImageChange}
            />
          </label>

          <button className="glass-button" type="submit" disabled={loading} style={{ marginTop: '12px' }}>
            {loading ? 'Posting...' : 'Post'}
          </button>
        </form>
      </div>
    </div>
  );
}
