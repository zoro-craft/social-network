'use client';

import { useAuth } from '@/app/context/AuthContext';
import { useToast } from '@/app/context/ToastContext';
import { useParams, useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { useComments } from './useComments';
import { API_BASE } from '@/utils/api';
import formatDate from '@/utils/format';
import Avatar from '@/components/Avatar';
import styles from '../../../styles/comments.module.css';

export default function CommentPage() {
  const { id } = useParams();
  const { auth, loading: authLoading } = useAuth();
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
  const [imageFile, setImageFile] = useState(null);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (!authLoading && !auth.authenticated) {
      router.push('/login');
    }
  }, [authLoading, auth.authenticated, router]);

  useEffect(() => {
    if (!auth.authenticated) return;
    fetchPostAndComments();
  }, [auth.authenticated]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!content.trim()) return;

    setSubmitting(true);
    try {
      const res = await addComment(content, imageFile);

      if (res.status === 401) {
        router.push('/login');
        return;
      }

      if (!res.ok) {
        const data = await res.json();
        showToast(data.message || 'Failed to add comment');
        return;
      }

      setContent('');
    } catch {
      showToast('Could not reach the server');
    } finally {
      setSubmitting(false);
    }
  };
  const handleImageChange = (e) => {
    const file = e.target.files[0];
    setImageFile(file || null);
  };


  if (authLoading || !auth.authenticated) {
    return null;
  }

  if (loading) {
    return (
      <div className="glass-feed-page">
        <p style={{ color: 'var(--text-muted)' }}>Loading...</p>
      </div>
    );
  }

  if (!post) {
    return (
      <div className="glass-feed-page">
        <p style={{ color: 'var(--text-muted)' }}>Post not found</p>
      </div>
    );
  }

  return (
    <div className="glass-feed-page">
      <div className="glass-post-card" style={{ padding: '32px' }}>
        
        {/* Original Post */}
        <div style={{ marginBottom: '32px', paddingBottom: '24px', borderBottom: '1px solid var(--border-strong)' }}>
          <div className="glass-post-header">
            <span className="glass-post-author">Posted by {post.nickname}</span>
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

        <h2 style={{ fontSize: '20px', fontWeight: '600', marginBottom: '20px', color: '#fff' }}>Comments</h2>

        <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '16px', marginBottom: '40px' }}>
          <div className="glass-input-group">
            <input
              className="glass-input"
              type="text"
              placeholder="Write a comment..."
              value={content}
              onChange={(e) => setContent(e.target.value)}
            />
          </div>
          
          <div style={{ display: 'flex', gap: '12px', alignItems: 'center', justifyContent: 'space-between' }}>
            <label className="glass-file-label" style={{ width: 'auto', flexDirection: 'row', alignItems: 'center' }}>
              <span style={{ marginRight: '8px' }}>Attach Image:</span>
              <input
                className="glass-file-input"
                type="file"
                accept="image/jpeg,image/png,image/gif"
                onChange={handleImageChange}
              />
            </label>
            <button className="glass-button" type="submit" disabled={submitting} style={{ width: 'auto', padding: '10px 32px', marginTop: '0' }}>
              {submitting ? 'Posting...' : 'Post'}
            </button>
          </div>
        </form>

        {comments.length === 0 && (
          <p style={{ color: 'var(--text-muted)', textAlign: 'center', margin: '40px 0' }}>No comments yet. Be the first!</p>
        )}

        <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
          {comments.map((c) => (
            <div key={c.id} style={{ 
              background: 'rgba(255, 255, 255, 0.03)', 
              padding: '20px', 
              borderRadius: '16px', 
              border: '1px solid var(--border)' 
            }}>
              <div
                onClick={() => router.push(`/profile/${c.user_id}`)}
                style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '14px', cursor: 'pointer' }}
              >
                <Avatar src={c.avatar} name={c.first_name} size={32} />
                <span style={{ fontWeight: '600', color: 'var(--text-primary)', fontSize: '14px' }}>{c.first_name}</span>
                <span style={{ color: 'var(--text-muted)', fontSize: '12px', marginLeft: 'auto' }}>{formatDate(c.created_at)}</span>
              </div>

              <p style={{ fontSize: '14.5px', color: 'var(--text-primary)', marginBottom: c.image ? '16px' : '0', lineHeight: '1.6' }}>{c.content}</p>
              
              {c.image && (
                // eslint-disable-next-line @next/next/no-img-element
                <img
                  style={{ width: '100%', maxHeight: '350px', objectFit: 'cover', borderRadius: '12px', border: '1px solid var(--border)' }}
                  src={`${API_BASE}${c.image}`}
                  alt=""
                />
              )}
            </div>
          ))}
        </div>

        {!loading && hasMore && (
          <button
            className="glass-button"
            onClick={loadMoreComments}
            disabled={loadingMore}
            style={{ marginTop: '32px', background: 'rgba(255, 255, 255, 0.05)', color: 'var(--text-primary)' }}
          >
            {loadingMore ? 'Loading...' : 'Load older comments'}
          </button>
        )}
      </div>
    </div>
  );
}
