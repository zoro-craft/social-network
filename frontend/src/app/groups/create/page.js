'use client';

import { useToast } from '@/app/context/ToastContext';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { apiRequest } from '@/utils/api';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

export default function CreatGroup() {
  const { auth, loading: authLoading } = useRequireAuth();
  const router = useRouter();
  const showToast = useToast();

  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!title.trim() || !description.trim()) {
      showToast('Title and description are required');
      return;
    }

    setLoading(true);

    try {
      const res = await apiRequest('/groups', 'POST', {
        titel: title.trim(),
        description: description.trim(),
      });

      showToast(res.message || 'Group created successfully!', false);
      router.push('/groups');
    } catch (err) {
      showToast(err.message || 'Failed to create group');
      setLoading(false);
    }
  };

  if (authLoading || !auth.authenticated) return null;

  return (
    <div className="glass-page">
      <div className="glass-container glass-container-lg">
        <h2>Create Group</h2>
        <p className="glass-subtitle">Start a new community</p>

        <form onSubmit={handleSubmit} className="glass-form">
          <div className="glass-input-group">
            <input
              className="glass-input"
              type="text"
              placeholder="Group title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
            />
          </div>

          <div className="glass-input-group">
            <textarea
              className="glass-textarea"
              placeholder="Describe your group..."
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={4}
              required
            />
          </div>

          <button
            className="glass-button"
            type="submit"
            disabled={loading}
          >
            {loading ? 'Creating...' : 'Create Group'}
          </button>
        </form>
      </div>
    </div>
  );
}
