'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { API_BASE } from '../../lib/api';
import styles from '../../styles/profile.module.css';
import { useAuth } from '@/app/context/AuthContext';
import { useToast } from '@/app/context/ToastContext';
import { useVisibility } from './useVisibility';
import { useFollow } from './useFollow';
import Avatar from '../../components/Avatar';

export default function ProfilePage() {
  const router = useRouter();
  const { id } = useParams();
  const { auth, loading: authLoading } = useAuth();
  const showToast = useToast();
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [togglingVisibility, setTogglingVisibility] = useState(false);

  const { toggleVisibility } = useVisibility(setProfile);
  const { status, fetchStatus, follow, unfollow } = useFollow(id);

  useEffect(() => {
    if (!auth.authenticated || auth.id == id) return;
    fetchStatus();
  }, [auth.authenticated, id]);

  const handleFollowClick = async () => {
    if (status === 'none') {
      await follow();
    } else {
      await unfollow();
    }
  };

  useEffect(() => {
    if (!authLoading && !auth.authenticated) {
      router.push('/login');
      return;
    }

    if (authLoading) return;

    async function fetchProfile() {
      try {
        const res = await fetch(`${API_BASE}/profile/${id}`, {
          credentials: 'include',
        });
        const data = await res.json();
        setProfile(data);
      } catch {
        setProfile(null);
      } finally {
        setLoading(false);
      }
    }

    fetchProfile();
  }, [authLoading, auth.authenticated, router, id]);

  const handleToggleVisibility = async () => {
    setTogglingVisibility(true);
    try {
      const res = await toggleVisibility();

      if (!res.ok) {
        showToast('Failed to update visibility');
        return;
      }

      showToast(
        profile.is_public ? 'Profile set to private' : 'Profile set to public',
        false,
      );
    } catch {
      showToast('Could not reach the server');
    } finally {
      setTogglingVisibility(false);
    }
  };

  if (loading) return <p className={styles.status}>Loading...</p>;
  if (!profile) return <p className={styles.status}>User not found</p>;

  const isOwnProfile = auth.id === profile.id;

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Avatar src={profile.avatar} name={profile.first_name} size={150} />

        <div className={styles.headerInfo}>
          <div className={styles.headerTop}>
            <h1 className={styles.name}>
              {profile.first_name} {profile.last_name}
            </h1>

            {isOwnProfile ? (
              <button
                className={styles.visibilityButton}
                onClick={handleToggleVisibility}
                disabled={togglingVisibility}
              >
                {profile.is_public ? '🌐 Public' : '🔒 Private'} 
              </button>
            ) : (
              <button
                className={styles.visibilityButton}
                onClick={handleFollowClick}
              >
                {status === 'accepted' && 'Following'}
                {status === 'pending' && 'Requested'}
                {status === 'none' && 'Follow'}
              </button>
            )}
          </div>

          <div className={styles.stats}>
            <span>
              <strong>{profile.followers_count}</strong> followers
            </span>
            <span>
              <strong>{profile.following_count}</strong> following
            </span>
          </div>

          <div className={styles.bio}>
            {profile.nickname && (
              <span className={styles.nickname}>@{profile.nickname}</span>
            )}

            {profile.about_me && (
              <p className={styles.aboutMe}>{profile.about_me}</p>
            )}
          </div>
        </div>
      </div>

      <h2 className={styles.postsTitle}>Posts</h2>

      {profile.restricted ? (
        <p className={styles.status}>🔒 This account is private</p>
      ) : (
        <>
          {(!profile.posts || profile.posts.length === 0) && (
            <p className={styles.status}>No posts yet</p>
          )}

          <div className={styles.grid}>
            {(profile.posts || []).map((post) => (
              <div key={post.id} className={styles.postCard}>
                {post.image && (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img
                    className={styles.postImage}
                    src={`${API_BASE}${post.image}`}
                    alt=""
                  />
                )}
                <p className={styles.postTitle}>{post.title}</p>
                <p className={styles.postContent}>{post.content}</p>
              </div>
            ))}
          </div>
        </>
      )}
    </div>
  );
}
