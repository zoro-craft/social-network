'use client';

import { useAuth } from '@/app/context/AuthContext';
import { useToast } from '@/app/context/ToastContext';
import { useParams, useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { useComments } from './useComments';
import { API_BASE } from '../../../lib/api';
import formatDate from '../../../lib/format';
import Avatar from '../../../components/Avatar';
import styles from '../../../styles/comments.module.css';

export default function CommentPage() {
  const { id } = useParams();
  const { loading: authLoading, auth } = useAuth();
  const router = useRouter();
  const showToast = useToast();

  const { post, comments, loading, fetchPostAndComments, addComment } =
    useComments(id);
  const [content, setContent] = useState('');
  const [imageFile, setImageFile] = useState(null);
  const [previewUrl, setPreviewUrl] = useState(null);
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
      setImageFile(null);
      setPreviewUrl(null);
    } catch {
      showToast('Could not reach the server');
    } finally {
      setSubmitting(false);
    }
  };

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (previewUrl) URL.revokeObjectURL(previewUrl);
    if (!file) {
      setImageFile(null);
      setPreviewUrl(null);
      return;
    }
    setImageFile(file);
    setPreviewUrl(URL.createObjectURL(file));
  };

  const handleRemoveImage = () => {
    if (previewUrl) URL.revokeObjectURL(previewUrl);
    setImageFile(null);
    setPreviewUrl(null);
  };

  if (authLoading || !auth.authenticated) {
    return null;
  }

  return (
    <div className={styles.container}>
      {loading && <p className={styles.status}>Loading...</p>}
      {!loading && !post && <p className={styles.status}>Post not found</p>}

      {post && (
        <div className={styles.post}>
          <p className={styles.author}>Posted by {post.nickname}</p>
          <p className={styles.postTitle}>{post.title}</p>
          {post.image && (
            // eslint-disable-next-line @next/next/no-img-element
            <img
              className={styles.postImage}
              src={`${API_BASE}${post.image}`}
              alt=""
            />
          )}
          <p className={styles.content}>{post.content}</p>
          <p className={styles.date}>{formatDate(post.created_at)}</p>
        </div>
      )}

      <h2 className={styles.heading}>Comments</h2>

      <form className={styles.form} onSubmit={handleSubmit}>
        <input
          className={styles.input}
          type="text"
          placeholder="Write a comment..."
          value={content}
          onChange={(e) => setContent(e.target.value)}
        />
        <div className={styles.buttonRow}>
          <label
            className={`${styles.imageLabel}${previewUrl ? ` ${styles.imageActive}` : ''}`}
            title="Attach image"
          >
            📷
            <input
              className={styles.imageInput}
              type="file"
              accept="image/jpeg,image/png,image/gif"
              onChange={handleImageChange}
            />
          </label>
          <button className={styles.button} type="submit" disabled={submitting}>
            {submitting ? 'Posting...' : 'Post'}
          </button>
        </div>
        {previewUrl && (
          <div className={styles.preview}>
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img className={styles.previewImage} src={previewUrl} alt="" />
            <button
              type="button"
              className={styles.previewRemove}
              onClick={handleRemoveImage}
              title="Remove image"
            >
              ✕
            </button>
          </div>
        )}
      </form>

      {!loading && comments.length === 0 && (
        <p className={styles.status}>No comments yet</p>
      )}

      {comments.map((c) => (
        <div key={c.id} className={styles.comment}>
          <div
            className={styles.commentHeader}
            onClick={() => router.push(`/profile/${c.user_id}`)}
          >
            <Avatar src={c.avatar} name={c.nickname} size={28} />
            <p className={styles.commentAuthor}>{c.nickname}</p>
            <p className={styles.commentDate}>{formatDate(c.created_at)}</p>
          </div>

          <p className={styles.commentContent}>{c.content}</p>
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
  );
}
