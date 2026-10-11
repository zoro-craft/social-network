'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useToast } from '@/app/context/ToastContext';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { useImageInput } from '@/hooks/useImageInput';
import { API_BASE } from '@/utils/api';

export default function CreateStoryPage() {
  const router = useRouter();
  const { auth, loading: authLoading } = useRequireAuth();
  const showToast = useToast();
  const { file: imageFile, handleChange: handleImageChange } = useImageInput();

  const [caption, setCaption] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!imageFile) {
      showToast('Pick an image for your story');
      return;
    }

    setLoading(true);

    try {
      const form = new FormData();
      form.append('caption', caption.trim());
      form.append('image', imageFile);

      const res = await fetch(`${API_BASE}/stories/create`, {
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
        showToast(data.error || 'Failed to create story');
        setLoading(false);
        return;
      }

      showToast('Story added!', false);
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
      <div className="glass-container glass-container-lg">
        <h2>New Story</h2>
        <p className="glass-subtitle">
          Your story disappears after 24 hours and is only visible to your followers.
        </p>

        <form onSubmit={handleSubmit} className="glass-form">
          <label className="glass-file-label">
            Story image
            <input
              className="glass-file-input"
              type="file"
              accept="image/jpeg,image/png,image/gif"
              onChange={handleImageChange}
            />
          </label>

          <div className="glass-input-group">
            <textarea
              className="glass-textarea"
              placeholder="Add a caption... (optional)"
              value={caption}
              onChange={(e) => setCaption(e.target.value)}
              maxLength={200}
              rows={3}
            />
          </div>

          <button className="glass-button" type="submit" disabled={loading}>
            {loading ? 'Sharing...' : 'Share to story'}
          </button>
        </form>
      </div>
    </div>
  );
}
