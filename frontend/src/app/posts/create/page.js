'use client';

import { useToast } from '@/app/context/ToastContext';
import { useRouter, useSearchParams } from 'next/navigation';
import { useEffect, useState, Suspense } from 'react';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { useImageInput } from '@/hooks/useImageInput';
import { API_BASE } from '@/utils/api';
import Avatar from '@/components/Avatar';
import styles from '../../styles/posts.module.css';

function CreatePostForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const preselectedGroupId = searchParams.get('group_id');
  const { auth, loading: authLoading } = useRequireAuth();
  const showToast = useToast();

  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [privacy, setPrivacy] = useState('public');
  const { file: imageFile, handleChange: handleImageChange } = useImageInput();
  const [followers, setFollowers] = useState([]);
  const [selectedViewers, setSelectedViewers] = useState([]);
  const [loading, setLoading] = useState(false);

  const [groups, setGroups] = useState([]);
  const [selectedGroupId, setSelectedGroupId] = useState(preselectedGroupId || '');

  useEffect(() => {
    if (!auth.id) return;

    const fetchGroups = async () => {
      try {
        const res = await fetch(`${API_BASE}/groups/my-groups`, {
          credentials: 'include',
        });
        if (res.ok) {
          const data = await res.json();
          setGroups(data || []);
        }
      } catch {
        setGroups([]);
      }
    };

    fetchGroups();
  }, [auth.id]);

  useEffect(() => {
    if (privacy !== 'private' || !auth.id || selectedGroupId) return;

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
  }, [privacy, auth.id, selectedGroupId]);

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
    if (!selectedGroupId && privacy === 'private' && selectedViewers.length === 0) {
      showToast('Select at least one follower');
      return;
    }

    setLoading(true);

    try {
      const form = new FormData();
      form.append('title', title);
      form.append('content', content);
      form.append('privacy', selectedGroupId ? 'public' : privacy);
      if (selectedGroupId) form.append('group_id', selectedGroupId);
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
        showToast(data.error || 'Failed to create post');
        setLoading(false);
        return;
      }

      showToast('Post created successfully!', false);
      if (selectedGroupId) {
        router.push(`/groups/${selectedGroupId}`);
      } else {
        router.push('/posts');
      }
    } catch {
      showToast('Could not reach the server');
      setLoading(false);
    }
  };

  if (authLoading || !auth.authenticated) {
    return null;
  }

  const isGroupPost = !!selectedGroupId;

  return (
    <div className="glass-page">
      <div className={`glass-container ${styles.container}`}>
        <h2>Create Post</h2>
        <p className={styles.subtitle}>
          {isGroupPost ? 'Post to your group' : 'Share something with your network'}
        </p>

        <form onSubmit={handleSubmit} className={styles.form}>
          {groups.length > 0 && (
            <div className="glass-input-group">
              <select
                className={`glass-input ${styles.select}`}
                value={selectedGroupId}
                onChange={(e) => {
                  setSelectedGroupId(e.target.value);
                  if (e.target.value) {
                    setPrivacy('public');
                    setSelectedViewers([]);
                  }
                }}
              >
                <option value="" className={styles.option}>Post to your feed</option>
                {groups.map((g) => (
                  <option key={g.id} value={g.id} className={styles.option}>
                    Group: {g.titel}
                  </option>
                ))}
              </select>
            </div>
          )}

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

          {!isGroupPost && (
            <div className="glass-input-group">
              <select
                className={`glass-input ${styles.select}`}
                value={privacy}
                onChange={(e) => setPrivacy(e.target.value)}
              >
                <option value="public" className={styles.option}>Public</option>
                <option value="almost_private" className={styles.option}>Almost private (followers)</option>
                <option value="private" className={styles.option}>Private (chosen followers)</option>
              </select>
            </div>
          )}

          {isGroupPost && (
            <p className={styles.subtitle}>Group posts are always public</p>
          )}

          {privacy === 'private' && !isGroupPost && (
            <div className={`glass-input-group ${styles.viewerGroup}`}>
              <label className={styles.viewerLabel}>
                Choose who can see this post:
              </label>
              {followers.length === 0 && (
                <p className={styles.noFollowers}>No followers yet</p>
              )}
              <div className={styles.viewerGrid}>
                {followers.map((f) => (
                  <div
                    key={f.follower_id}
                    onClick={() => toggleViewer(f.follower_id)}
                    className={`${styles.viewerChip} ${
                      selectedViewers.includes(f.follower_id)
                        ? styles.viewerChipActive
                        : styles.viewerChipInactive
                    }`}
                  >
                    <Avatar src={f.avatar} name={f.first_name} size={32} />
                    <span className={styles.viewerName}>{f.first_name}</span>
                  </div>
                ))}
              </div>
            </div>
          )}

          <label className={`glass-file-label ${styles.fileLabel}`}>
            Image
            <input
              className="glass-file-input"
              type="file"
              accept="image/jpeg,image/png,image/gif"
              onChange={handleImageChange}
            />
          </label>

          <button className={`glass-button ${styles.submitButton}`} type="submit" disabled={loading}>
            {loading ? 'Posting...' : 'Post'}
          </button>
        </form>
      </div>
    </div>
  );
}

export default function CreatePostPage() {
  return (
    <Suspense fallback={<p>Loading...</p>}>
      <CreatePostForm />
    </Suspense>
  );
}