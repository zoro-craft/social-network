'use client';

import { useToast } from '@/app/context/ToastContext';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { useImageInput } from '@/hooks/useImageInput';
import { apiRequest } from '@/utils/api';
import { useRouter } from 'next/navigation';
import { useEffect, useRef, useState } from 'react';

export default function CreatGroup() {
  const { auth, loading: authLoading } = useRequireAuth();
  const router = useRouter();
  const showToast = useToast();

  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const { file: imageFile, setFile: setImageFile } = useImageInput();
  const [imagePreview, setImagePreview] = useState('');
  const [loading, setLoading] = useState(false);
  const previewUrlRef = useRef(null);

  useEffect(() => () => {
    if (previewUrlRef.current) URL.revokeObjectURL(previewUrlRef.current);
  }, []);

  const handleImageChange = (event) => {
    if (previewUrlRef.current) URL.revokeObjectURL(previewUrlRef.current);

    const file = event.target.files?.[0] || null;
    if (!file) {
      previewUrlRef.current = null;
      setImageFile(null);
      setImagePreview('');
      event.target.value = '';
      return;
    }

    const preview = URL.createObjectURL(file);
    previewUrlRef.current = preview;
    setImageFile(file);
    setImagePreview(preview);
    event.target.value = '';
  };

  const clearImage = () => {
    if (previewUrlRef.current) URL.revokeObjectURL(previewUrlRef.current);
    previewUrlRef.current = null;
    setImageFile(null);
    setImagePreview('');
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!title.trim() || !description.trim()) {
      showToast('Title and description are required');
      return;
    }

    setLoading(true);

    try {
      const form = new FormData();
      form.append('titel', title.trim());
      form.append('description', description.trim());
      if (imageFile) form.append('image', imageFile);

      const res = await apiRequest('/groups', 'POST', form, true);

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
              maxLength={50}
              required
            />
          </div>

          <div className="glass-input-group">
            <textarea
              className="glass-textarea"
              placeholder="Describe your group..."
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              maxLength={200}
              rows={4}
              required
            />
          </div>

          <label className="glass-file-label">
            Group image
            <input
              className="glass-file-input"
              type="file"
              accept="image/jpeg,image/png,image/gif"
              onChange={handleImageChange}
            />
          </label>

          {imagePreview && (
            <div className="glass-group-image-preview">
              <img src={imagePreview} alt="Selected group image" />
              <button type="button" onClick={clearImage}>
                Remove image
              </button>
            </div>
          )}

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
