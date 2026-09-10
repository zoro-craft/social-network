'use client';

import { apiRequest } from "@/utils/api";
import { useState, useEffect } from "react";
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { useRouter } from "next/navigation";
import GroupCard from '@/app/groups/GroupCard';
import MyGroupsList from '@/app/groups/MyGroupsList';
import styles from '../styles/groups.module.css';

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
            <div className={styles.page}>
                <p className={styles.mutedText}>Loading groups...</p>
            </div>
        );
    }

    if (!auth.authenticated) return null;

    return (
        <main className={styles.page}>
            <div className={styles.header}>
                <h1 className={styles.pageTitle}>Groups</h1>

                <div className={styles.headerActions}>
                    <button
                        className={styles.createButton}
                        onClick={() => router.push('/groups/create')}
                    >
                        <span className={styles.buttonIcon} aria-hidden="true">
                            +
                        </span>
                        Create Group
                    </button>

                    <button
                        className={styles.browseButton}
                        onClick={() => setActiveTab('all')}
                    >
                        All Groups
                    </button>

                    <button
                        className={styles.browseButton}
                        onClick={() => setActiveTab('my')}
                    >
                        My Groups
                    </button>
                </div>
            </div>

            <section className={styles.groupsSection} id="my-groups">
                <h2 className={styles.sectionTitle}>
                    {activeTab === 'all' ? 'All Groups' : 'My Groups'}
                </h2>

                {activeTab === 'all' ? (
                    groups.length === 0 ? (
                        <p className={styles.emptyState}>
                            No groups yet. Create the first one!
                        </p>
                    ) : (
                        <div className={styles.container}>
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