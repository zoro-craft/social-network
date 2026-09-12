'use client';

import { useParams, useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { API_BASE } from '@/utils/api';
import GroupPosts from '@/app/groups/GroupPosts';
import GroupEvent from '@/app/groups/Groupevents';

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
            <main className="glass-group-page">
                <p className="glass-muted-text">Loading group...</p>
            </main>
        );
    }

    if (!auth.authenticated) return null;

    if (!group) {
        return (
            <main className="glass-group-page">
                <p className="glass-muted-text">Group not found</p>
            </main>
        );
    }

    return (
        <main className="glass-group-page">
            <div className="glass-group-wrapper">
                <div className="glass-toolbar">
                    <button
                        className="glass-btn-back"
                        onClick={() => router.push('/groups')}
                    >
                        <span aria-hidden="true">&larr;</span>
                        All Groups
                    </button>
                </div>

                <header className="glass-header-card">
                    <div className="glass-group-icon glass-group-icon--lg" aria-hidden="true">
                        <span className="glass-person-one" />
                        <span className="glass-person-two" />
                    </div>

                    <div className="glass-header-info">
                        <h1 className="glass-group-title">{group.titel}</h1>
                        <p className="glass-group-description">{group.description}</p>
                    </div>
                </header>

                <section className="glass-group-actions" aria-label="Group actions">
                    <button
                        className="glass-action-btn glass-btn-invite"
                        onClick={() => router.push(`/groups/${id}/invite`)}
                    >
                        Invite Members
                    </button>

                    <button
                        className="glass-action-btn glass-btn-primary"
                        onClick={() => router.push(`/posts/create?group_id=${id}`)}
                    >
                        Create Post
                    </button>

                    <button
                        className="glass-action-btn glass-btn-primary"
                        onClick={() => router.push(`/groups/${id}/events/create`)}
                    >
                        Create Event
                    </button>

                    <button
                        className="glass-action-btn"
                        onClick={() => router.push(`/groups/${id}/chat`)}
                    >
                        Go to Chat
                    </button>
                </section>

                <section style={{ marginTop: 32 }}>
                    <h2 className="glass-section-title">Posts</h2>
                    <GroupPosts groupId={groupId} />
                </section>

                <section style={{ marginTop: 32 }}>
                    <h2 className="glass-section-title">Events</h2>
                    <GroupEvent group_id={groupId} />
                </section>
            </div>
        </main>
    );
}