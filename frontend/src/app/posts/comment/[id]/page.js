'use client';

import { useToast } from '@/app/context/ToastContext';
import { useParams, useRouter } from 'next/navigation';
import { useEffect, useRef, useState } from 'react';
import { useComments } from './useComments';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { useImageInput } from '@/hooks/useImageInput';
import { API_BASE } from '@/utils/api';
import formatDate from '@/utils/format';
import Avatar from '@/components/Avatar';

export default function CommentPage() {
  const { id } = useParams();
  const { auth, loading: authLoading } = useRequireAuth();
  const router = useRouter();
  const showToast = useToast();

  const {
    post,
    comments,
    loading,
    loadingMore,
    hasMore,
    fetchPostAndComments,
    loadMoreComments,
    addComment,
  } = useComments(id);
  const [content, setContent] = useState('');
  const { file: imageFile, setFile: setImageFile, handleChange: handleImageChange } = useImageInput();
  const imageInputRef = useRef(null);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (!auth.authenticated) return;
    fetchPostAndComments();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [auth.authenticated]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!content.trim()) {
      showToast('Content is required');
      return;
    }

    setSubmitting(true);
    try {
      const res = await addComment(content, imageFile);

      if (res.status === 401) {
        router.push('/login');
        return;
      }

      if (!res.ok) {
        const data = await res.json();
        showToast(data.error || 'Failed to add comment');
        return;
      }

      setContent('');
      setImageFile(null);
      if (imageInputRef.current) imageInputRef.current.value = '';
      showToast('Comment added successfully', false);
    } catch {
      showToast('Could not reach the server');
    } finally {
      setSubmitting(false);
    }
  };

  if (authLoading || !auth.authenticated) {
    return null;
  }

  if (loading) {
    return (
      <div className="glass-feed-page">
        <p className="glass-muted-text">Loading...</p>
      </div>
    );
  }

  if (!post) {
    return (
      <div className="glass-feed-page">
        <p className="glass-muted-text">Post not found</p>
      </div>
    );
  }

  return (
    <div className="glass-feed-page">
      <div className="glass-post-card glass-post-container">
        {/* Original Post */}
        <div className="glass-original-post">
          <div className="glass-post-header">
            <span className="glass-post-author">
              Posted by {post.full_name}
            </span>
            <span>{formatDate(post.created_at)}</span>
          </div>
          <p className="glass-post-title">{post.title}</p>
          {post.image && (
            // eslint-disable-next-line @next/next/no-img-element
            <img
              className="glass-post-image"
              src={`${API_BASE}${post.image}`}
              alt=""
            />
          )}
          <p className="glass-post-content">{post.content}</p>
        </div>

        <h2 className="glass-comments-heading">Comments</h2>

        <form onSubmit={handleSubmit} className="glass-comment-form">
          <div className="glass-input-group">
            <input
              className="glass-input"
              type="text"
              placeholder="Write a comment..."
              value={content}
              onChange={(e) => setContent(e.target.value)}
              disabled={submitting}
            />
          </div>

          <div className="glass-comment-form-row">
            <label className="glass-file-label glass-file-label-inline">
              <span>Attach Image:</span>
              <input
                className="glass-file-input"
                type="file"
                accept="image/jpeg,image/png,image/gif"
                onChange={handleImageChange}
                disabled={submitting}
                ref={imageInputRef}
              />
            </label>
            <button
              className="glass-button glass-button-inline"
              type="submit"
              disabled={submitting}
            >
              {submitting ? 'Posting...' : 'Post'}
            </button>
          </div>
        </form>

        {comments.length === 0 && (
          <p className="glass-muted-text glass-empty-comments">No comments yet. Be the first!</p>
        )}

        <div className="glass-comments-list">
          {comments.map((c) => (
            <div key={c.id} className="glass-comment-card">
              <div
                onClick={() => router.push(`/profile/${c.user_id}`)}
                className="glass-comment-header"
              >
                <Avatar src={c.avatar} name={c.first_name} size={32} />
                <span className="glass-comment-author">{c.first_name}</span>
                <span className="glass-comment-date">
                  {formatDate(c.created_at)}
                </span>
              </div>

              <p
                className={`glass-comment-content ${c.image ? 'glass-comment-content-with-image' : ''}`}
              >
                {c.content}
              </p>

              {c.image && (
                // eslint-disable-next-line @next/next/no-img-element
                <img
                  className="glass-comment-image"
                  src={`${API_BASE}${c.image}`}
                  alt=""
                />
              )}
            </div>
          ))}
        </div>

        {!loading && hasMore && (
          <button
            className="glass-button glass-load-more"
            onClick={loadMoreComments}
            disabled={loadingMore}
          >
            {loadingMore ? 'Loading...' : 'Load older comments'}
          </button>
        )}
      </div>
    </div>
  );
}
