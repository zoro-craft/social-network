'use client';

import { useParams, useRouter } from 'next/navigation';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import styles from '@/app/styles/group-details.module.css';

export default function GroupDetailsPage() {
    const { auth, loading } = useRequireAuth();
    const router = useRouter();
    const { id } = useParams();

    const group = {
        titel: 'Design Club',
        description: 'A place for designers to share ideas, feedback, and inspiration.',
    };

    if (loading) {
        return (
            <main className={styles.page}>
                <p className={styles.mutedText}>Loading group...</p>
            </main>
        );
    }

    if (!auth.authenticated) return null;

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
                    <div className={styles.groupIcon} aria-hidden="true">
                        <span className={styles.personOne} />
                        <span className={styles.personTwo} />
                    </div>

                    <div className={styles.headerInfo}>
                        <h1 className={styles.title}>{group.titel}</h1>
                        <p className={styles.description}>{group.description}</p>

                        <div className={styles.meta}>
                            <span className={styles.memberBadge}>
                                {group.member_count} members
                            </span>
                        </div>
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
                        onClick={() => goTo('/posts/create')}
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
                        <div className={styles.emptyCard}>
                            <p className={styles.emptyText}>No posts yet</p>
                        </div>
                    </div>

                    <div className={styles.section}>
                        <h2 className={styles.sectionTitle}>Events</h2>
                        <div className={styles.emptyCard}>
                            <p className={styles.emptyText}>No events yet</p>
                        </div>
                    </div>
                </section>
            </div>
        </main>
    );
}