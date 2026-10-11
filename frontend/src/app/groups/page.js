'use client';

import { useEffect, useMemo, useRef, useState } from 'react';
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
    const [search, setSearch] = useState('');
    const showToastRef = useRef(showToast);

    useEffect(() => {
        showToastRef.current = showToast;
    }, [showToast]);

    useEffect(() => {
        if (!auth.authenticated) return;

        const fetchGroups = async () => {
            try {
                const res = await apiRequest('/groups');
                const all = res || [];
                const notMember = all.filter((group) => group.is_member === false);

                setAllGroups(all);
                setGroups(notMember);
            } catch (error) {
                showToastRef.current(error.message || 'Failed to load groups');
                setAllGroups([]);
                setGroups([]);
            }
        };

        fetchGroups();
    }, [auth.authenticated]);

    const showAllGroups = () => {
        setGroups(allGroups.filter((group) => group.is_member === false));
    };

    const showMyGroups = () => {
        setGroups(allGroups.filter((group) => group.is_member === true));
    };

    const visibleGroups = useMemo(() => {
        const query = search.trim().toLowerCase();
        if (!query) return groups;
        return groups.filter((group) =>
            `${group.titel} ${group.description}`.toLowerCase().includes(query)
        );
    }, [groups, search]);

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
                            className="glass-btn-create glass-btn-primary"
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
                            Discover
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
                    <div className="glass-group-search-wrap">
                        <input
                            className="glass-group-search"
                            type="search"
                            placeholder="Search groups..."
                            value={search}
                            onChange={(event) => setSearch(event.target.value)}
                            aria-label="Search groups"
                        />
                    </div>

                    <div className="glass-groups-grid">
                        {visibleGroups.length === 0 ? (
                            <p className="glass-empty-state">
                                No groups found.
                            </p>
                        ) : (
                            visibleGroups.map((group) => (
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
