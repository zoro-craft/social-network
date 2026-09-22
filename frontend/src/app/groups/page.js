'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { apiRequest } from '@/utils/api';
import GroupCard from './GroupCard';
import { useToast } from '@/app/context/ToastContext';

export default function Groupspage() {
    const { auth, loading } = useRequireAuth();
    const router = useRouter();
    const showToast = useToast();

    const [allGroups, setAllGroups] = useState([]);
    const [groups, setGroups] = useState([]);

    useEffect(() => {
        if (!auth.authenticated) return;

        const fetchGroups = async () => {
            try {
                const res = await apiRequest('/groups');

                const all = res || [];

                // Default: show groups where user is NOT a member
                const notMember = all.filter(
                    (group) => group.is_member === false
                );

                setAllGroups(all);
                setGroups(notMember);
            } catch (error) {
                showToast(error.message || 'Failed to load groups');
                console.log(error);
                setAllGroups([]);
                setGroups([]);
            }
        };

        fetchGroups();
    }, [auth.authenticated]);

    // All Groups
    const showAllGroups = () => {
        const notMember = allGroups.filter(
            (group) => group.is_member === false
        );

        setGroups(notMember);
    };

    // My Groups
    const showMyGroups = () => {
        const member = allGroups.filter(
            (group) => group.is_member === true
        );

        setGroups(member);
    };

    if (loading) {
        return (
            <main className="glass-group-page">
                <div className="glass-group-wrapper">
                    <p className="glass-muted-text">
                        Loading groups...
                    </p>
                </div>
            </main>
        );
    }

    if (!auth.authenticated) return null;

    return (
        <main className="glass-group-page">
            <div className="glass-group-wrapper">
                <header className="glass-group-header">
                    <h1 className="glass-groups-title">
                        Groups
                    </h1>

                    <div className="glass-header-actions">

                        <button
                            className={`glass-btn-create glass-btn-primary`}
                            onClick={() => router.push('/groups/create')}
                        >
                            <span
                                className="glass-btn-icon"
                                aria-hidden="true"
                            >
                                +
                            </span>

                            Create Group
                        </button>

                        <button
                            className="glass-btn-browse"
                            onClick={showAllGroups}
                        >
                            All Groups
                        </button>

                        <button
                            className="glass-btn-browse"
                            onClick={showMyGroups}
                        >
                            My Groups
                        </button>

                    </div>
                </header>

                <section className="glass-groups-section">
                    <div className="glass-groups-grid">

                        {groups.length === 0 ? (
                            <p className="glass-empty-state">
                                No groups found.
                            </p>
                        ) : (
                            groups.map((group) => (
                                <GroupCard
                                    key={group.id}
                                    group={group}
                                />
                            ))
                        )}

                    </div>
                </section>
            </div>
        </main>
    );
}