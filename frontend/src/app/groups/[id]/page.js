'use client';

import { useParams, useRouter } from 'next/navigation';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { useToast } from '@/app/context/ToastContext';
import styles from '@/app/styles/group-details.module.css';
import { useEffect, useState } from 'react';
import { apiRequest } from '@/utils/api';
import GroupPosts from '@/app/groups/GroupPosts';
import GroupEvent from '@/app/groups/Groupevents';

export default function GroupDetailsPage() {
    const { auth, loading } = useRequireAuth();
    const router = useRouter();
    const { id } = useParams();
    const [group, setGroupInfo] = useState(null);
    const showToast = useToast();


    useEffect(() => {
        if (!auth.authenticated || !id) return;

        const fetchGroupInfo = async () => {
            try {
                const res = await apiRequest(`/groups/${id}`, "POST");
                setGroupInfo(res);

            } catch (error) {
                showToast('Unable to load group');
                router.push("/groups")
            }
        };

        fetchGroupInfo();
    }, [auth.authenticated, id]);

    if (loading) {
        return (
            <main className={styles.page}>
                <p className={styles.mutedText}>Loading group...</p>
            </main>
        );
    }

    if (!auth.authenticated) return null;

    if (!group) return null;

    const goTo = (path) => router.push(`/groups/${id}${path}`);

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
                    <img
                        className={styles.groupIcon}
                        src="/group-avatar.jpeg"
                        alt={`${group.titel} avatar`}
                    />

                    <div className={styles.headerInfo}>
                        <h1 className={styles.title}>{group.titel}</h1>
                        <p className={styles.description}>{group.description}</p>

                    </div>
                </header>

                <section className={styles.actions} aria-label="Group actions">
                    <button
                        className={`${styles.actionButton} ${styles.inviteButton}`}
                        onClick={() => goTo('/invite')}
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
                        onClick={() => goTo('/events/create')}
                    >
                        Create Event
                    </button>

                    <button
                        className={styles.actionButton}
                        onClick={() => goTo('/chat')}
                    >
                        Go to Chat
                    </button>
                </section>

                <section className={styles.content}>
                    <div className={styles.section}>
                        <h2 className={styles.sectionTitle}>Posts</h2>
                        <GroupPosts groupId={Number(id)} />
                    </div>

                    <div className={styles.section}>
                        <h2 className={styles.sectionTitle}>Events</h2>
                        <GroupEvent group_id={Number(id)} />
                    </div>
                </section>
            </div>
        </main>
    );
}