'use client';

import { useToast } from '@/app/context/ToastContext';
import { useParams, useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { useComments } from './useComments';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { API_BASE } from '@/utils/api';
import formatDate from '@/utils/format';
import Avatar from '@/components/Avatar';
import styles from '../../../styles/posts.module.css';

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
  const [imageFile, setImageFile] = useState(null);
  const [submitting, setSubmitting] = useState(false);

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
        <p className={styles.mutedText}>Loading...</p>
      </div>
    );
  }

  if (!post) {
    return (
      <div className="glass-feed-page">
        <p className={styles.mutedText}>Post not found</p>
      </div>
    );
  }

  return (
    <div className="glass-feed-page">
      <div className={`glass-post-card ${styles.postContainer}`}>
        
        {/* Original Post */}
        <div className={styles.originalPost}>
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

        <h2 className={styles.commentsHeading}>Comments</h2>

        <form onSubmit={handleSubmit} className={styles.commentForm}>
          <div className="glass-input-group">
            <input
              className="glass-input"
              type="text"
              placeholder="Write a comment..."
              value={content}
              onChange={(e) => setContent(e.target.value)}
            />
          </div>
          
          <div className={styles.commentFormRow}>
            <label className={`glass-file-label ${styles.fileLabelInline}`}>
              <span>Attach Image:</span>
              <input
                className="glass-file-input"
                type="file"
                accept="image/jpeg,image/png,image/gif"
                onChange={handleImageChange}
              />
            </label>
            <button className={`glass-button ${styles.submitButtonInline}`} type="submit" disabled={submitting}>
              {submitting ? 'Posting...' : 'Post'}
            </button>
          </div>
        </form>

        {comments.length === 0 && (
          <p className={styles.emptyComments}>No comments yet. Be the first!</p>
        )}

        <div className={styles.commentsList}>
          {comments.map((c) => (
            <div key={c.id} className={styles.commentCard}>
              <div
                onClick={() => router.push(`/profile/${c.user_id}`)}
                className={styles.commentHeader}
              >
                <Avatar src={c.avatar} name={c.first_name} size={32} />
                <span className={styles.commentAuthor}>{c.first_name}</span>
                <span className={styles.commentDate}>{formatDate(c.created_at)}</span>
              </div>

              <p className={`${styles.commentContent} ${c.image ? styles.commentContentWithImage : ''}`}>{c.content}</p>
              
              {c.image && (
                // eslint-disable-next-line @next/next/no-img-element
                <img
                  className={styles.commentImage}
                  src={`${API_BASE}${c.image}`}
                  alt=""
                />
              )}
            </div>
          ))}
        </div>

        {!loading && hasMore && (
          <button
            className={`glass-button ${styles.loadOlderButton}`}
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
