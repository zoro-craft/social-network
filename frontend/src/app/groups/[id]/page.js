'use client';

import { useParams, useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { API_BASE } from '@/utils/api';
import GroupPosts from '@/app/groups/GroupPosts';
import GroupEvent from '@/app/groups/Groupevents';
import styles from '@/app/styles/group-details.module.css';

export default function GroupDetailsPage() {
    const { auth, loading } = useRequireAuth();
    const router = useRouter();
    const { id } = useParams();
    const groupId = Number(id);

    const [group, setGroup] = useState(null);

    useEffect(() => {
        if (!auth.authenticated || !groupId) return;

        const fetchGroup = async () => {
            try {
                const res = await fetch(`${API_BASE}/groups`, {
                    credentials: 'include',
                });
                if (res.ok) {
                    const groups = await res.json();
                    const found = (groups || []).find((g) => g.id === groupId);
                    setGroup(found || null);
                }
            } catch {
                setGroup(null);
            }
        };

        fetchGroup();
    }, [auth.authenticated, groupId]);

    if (loading) {
        return (
            <main className={styles.page}>
                <p className={styles.mutedText}>Loading group...</p>
            </main>
        );
    }

    if (!auth.authenticated) return null;

    if (!group) {
        return (
            <main className={styles.page}>
                <p className={styles.mutedText}>Group not found</p>
            </main>
        );
    }

    return (
        <main className={styles.page}>
            <div className={styles.wrapper}>
                <div className={styles.toolbar}>
                    <button
                        className={styles.backButton}
                        onClick={() => router.push('/groups')}
                    >
                        <span aria-hidden="true">&larr;</span>
                        All Groups
                    </button>
                </div>

                <header className={styles.headerCard}>
                    <div className={styles.groupIcon} aria-hidden="true">
                        <span className={styles.personOne} />
                        <span className={styles.personTwo} />
                    </div>

                    <div className={styles.headerInfo}>
                        <h1 className={styles.title}>{group.titel}</h1>
                        <p className={styles.description}>{group.description}</p>
                    </div>
                </header>

                <section className={styles.actions} aria-label="Group actions">
                    <button
                        className={`${styles.actionButton} ${styles.inviteButton}`}
                        onClick={() => router.push(`/groups/${id}/invite`)}
                    >
                        Invite Members
                    </button>

                    <button
                        className={`${styles.actionButton} ${styles.primaryButton}`}
                        onClick={() => router.push(`/posts/create?group_id=${id}`)}
                    >
                        Create Post
                    </button>

                    <button
                        className={`${styles.actionButton} ${styles.primaryButton}`}
                        onClick={() => router.push(`/groups/${id}/events/create`)}
                    >
                        Create Event
                    </button>

                    <button
                        className={styles.actionButton}
                        onClick={() => router.push(`/groups/${id}/chat`)}
                    >
                        Go to Chat
                    </button>
                </section>

                <section style={{ marginTop: 32 }}>
                    <h2 className={styles.sectionTitle}>Posts</h2>
                    <GroupPosts groupId={groupId} />
                </section>

                <section style={{ marginTop: 32 }}>
                    <h2 className={styles.sectionTitle}>Events</h2>
                    <GroupEvent group_id={groupId} />
                </section>
            </div>
        </main>
    );
}