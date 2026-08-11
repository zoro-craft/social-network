'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { API_BASE } from '../../lib/api';
import styles from '../../styles/profile.module.css';
import { useAuth } from '@/app/context/AuthContext';
import Avatar from '../../components/Avatar';

export default function ProfilePage() {
  const router = useRouter();
  const { id } = useParams();
  const { auth, loading: authLoading } = useAuth();
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);

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

  if (loading) return <p className={styles.status}>Loading...</p>;
  if (!profile) return <p className={styles.status}>User not found</p>;

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Avatar src={profile.avatar} name={profile.first_name} size={90} />

        <div>
          <h1 className={styles.name}>
            full name :{profile.first_name} {profile.last_name}
          </h1>
          {profile.nickname && (
            <p className={styles.nickname}>@{profile.nickname}</p>
          )}
          {profile.about_me && (
            <p className={styles.aboutMe}>{profile.about_me}</p>
          )}
        </div>
      </div>

      <h2 className={styles.postsTitle}>Posts</h2>

      {profile.posts.length === 0 && (
        <p className={styles.status}>No posts yet</p>
      )}

      <div className={styles.grid}>
        {profile.posts.map((post) => (
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
    </div>
  );
}
