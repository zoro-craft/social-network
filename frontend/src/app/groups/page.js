'use client';

import { apiRequest } from "@/utils/api";
import { useState, useEffect } from "react";
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { useRouter } from "next/navigation";
import GroupCard from '@/app/groups/GroupCard';
import MyGroupsList from '@/app/groups/MyGroupsList';

export default function Groupspage() {
    const { auth, loading } = useRequireAuth();
    const [groups, setGroups] = useState([]);
    const [activeTab, setActiveTab] = useState('all');

    const router = useRouter();

    useEffect(() => {
        if (!auth.authenticated) return;

        const fetchGroups = async () => {
            try {
                const res = await apiRequest("/groups");
                setGroups(res || []);
            } catch {
                setGroups([]);
            }
        };

        fetchGroups();
    }, [auth.authenticated]);

    if (loading) {
        return (
            <div className="glass-group-page">
                <p className="glass-muted-text">Loading groups...</p>
            </div>
        );
    }

    if (!auth.authenticated) return null;

    return (
        <main className="glass-group-page">
            <div className="glass-group-header">
                <h1 className="glass-groups-title">Groups</h1>

                <div className="glass-header-actions">
                    <button
                        className="glass-btn-create"
                        onClick={() => router.push('/groups/create')}
                    >
                        <span className="glass-btn-icon" aria-hidden="true">
                            +
                        </span>
                        Create Group
                    </button>

                    <button
                        className="glass-btn-browse"
                        onClick={() => setActiveTab('all')}
                    >
                        All Groups
                    </button>

                    <button
                        className="glass-btn-browse"
                        onClick={() => setActiveTab('my')}
                    >
                        My Groups
                    </button>
                </div>
            </div>

            <section className="glass-groups-section" id="my-groups">
                <h2 className="glass-section-title">
                    {activeTab === 'all' ? 'All Groups' : 'My Groups'}
                </h2>

                {activeTab === 'all' ? (
                    groups.length === 0 ? (
                        <p className="glass-empty-state">
                            No groups yet. Create the first one!
                        </p>
                    ) : (
                        <div className="glass-groups-grid">
                            {groups.map((group) => (
                                <GroupCard
                                    key={group.id}
                                    group={group}
                                />
                            ))}
                        </div>
                    )
                ) : (
                    <MyGroupsList />
                )}
            </section>
        </main>
    );
}