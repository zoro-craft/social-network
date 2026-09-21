'use client';

import { useParams, useRouter } from 'next/navigation';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { useToast } from '@/app/context/ToastContext';
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
                showToast(error.message || 'Unable to load group');
                router.push("/groups")
            }
        };

        fetchGroupInfo();
    }, [auth.authenticated, id]);

    if (loading) {
        return (
            <main className="glass-group-page">
                <div className="glass-group-wrapper">
                    <p className="glass-muted-text">Loading group...</p>
                </div>
            </main>
        );
    }

    if (!auth.authenticated) return null;

    if (!group) return null;

    const goTo = (path) => router.push(`/groups/${id}${path}`);

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
                    <img
                        className="glass-group-icon glass-group-icon--lg"
                        src="/group-avatar.jpeg"
                        alt={`${group.titel} avatar`}
                    />

                    <div className="glass-header-info">
                        <h1 className="glass-group-title">{group.titel}</h1>
                        <p className="glass-group-description">{group.description}</p>

                    </div>
                </header>

                <section className="glass-group-actions" aria-label="Group actions">
                    <button
                        className="glass-action-btn"
                        onClick={() => goTo('/invite')}
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
                        onClick={() => goTo('/events/create')}
                    >
                        Create Event
                    </button>

                    <button
                        className="glass-action-btn"
                        onClick={() => goTo('/members')}
                    >
                        Group Members
                    </button>
                </section>

                <section className="glass-group-content">
                    <div className="glass-group-column">
                        <h2 className="glass-section-title">Posts</h2>
                        <GroupPosts groupId={Number(id)} />
                    </div>

                    <div className="glass-group-column">
                        <h2 className="glass-section-title">Events</h2>
                        <GroupEvent group_id={Number(id)} />
                    </div>
                </section>
            </div>
        </main>
    );
}