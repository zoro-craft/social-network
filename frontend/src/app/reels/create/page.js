'use client';

import { useEffect, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useToast } from '@/app/context/ToastContext';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { useImageInput } from '@/hooks/useImageInput';
import { API_BASE } from '@/utils/api';

// Matches utils.MaxVideoUploadSize on the backend, so oversized files fail fast
// instead of after a long upload.
const MAX_VIDEO_BYTES = 100 * 1024 * 1024;

export default function CreateReelPage() {
  const router = useRouter();
  const { auth, loading: authLoading } = useRequireAuth();
  const showToast = useToast();
  const { file: videoFile, setFile: setVideoFile } = useImageInput();

  const [caption, setCaption] = useState('');
  const [loading, setLoading] = useState(false);
  const [previewUrl, setPreviewUrl] = useState(null);

  const previewUrlRef = useRef(null);

  // Built in the change handler rather than an effect: object URLs must be revoked
  // when replaced, or the picked file stays in memory for the rest of the session.
  const handleVideoChange = (event) => {
    if (previewUrlRef.current) URL.revokeObjectURL(previewUrlRef.current);

    const file = event.target.files[0] || null;

    if (!file) {
      previewUrlRef.current = null;
      setVideoFile(null);
      setPreviewUrl(null);
      return;
    }

    const url = URL.createObjectURL(file);
    previewUrlRef.current = url;
    setVideoFile(file);
    setPreviewUrl(url);
  };

  useEffect(
    () => () => {
      if (previewUrlRef.current) URL.revokeObjectURL(previewUrlRef.current);
    },
    [],
  );

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!videoFile) {
      showToast('Pick a video for your reel');
      return;
    }

    if (videoFile.size > MAX_VIDEO_BYTES) {
      showToast('That video is too large (100MB max)');
      return;
    }

    setLoading(true);

    try {
      const form = new FormData();
      form.append('caption', caption.trim());
      form.append('video', videoFile);

      const res = await fetch(`${API_BASE}/reels/create`, {
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
        showToast(data.error || 'Failed to create reel');
        setLoading(false);
        return;
      }

      showToast('Reel posted!', false);
      router.push('/reels');
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
        <h2>New Reel</h2>
        <p className="glass-subtitle">
          Share a short video. MP4, MOV or WebM up to 100MB — reels stay in your profile.
        </p>

        <form onSubmit={handleSubmit} className="glass-form">
          <label className="glass-file-label">
            Reel video
            <input
              className="glass-file-input"
              type="file"
              accept="video/mp4,video/quicktime,video/webm,video/x-m4v"
              onChange={handleVideoChange}
            />
          </label>

          {previewUrl && (
            <video className="glass-reel-preview" src={previewUrl} controls playsInline />
          )}

          <div className="glass-input-group">
            <textarea
              className="glass-textarea"
              placeholder="Write a caption... (optional)"
              value={caption}
              onChange={(e) => setCaption(e.target.value)}
              maxLength={300}
              rows={3}
            />
          </div>

          <button className="glass-button" type="submit" disabled={loading}>
            {loading ? 'Posting...' : 'Share reel'}
          </button>
        </form>
      </div>
    </div>
  );
}
