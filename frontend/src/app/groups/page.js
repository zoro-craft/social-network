'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { apiRequest } from '@/utils/api';
import GroupCard from './GroupCard';
import styles from '@/app/styles/groups.module.css';
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
                showToast('Failed to load groups');
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
            <div className={styles.page}>
                <p className={styles.mutedText}>
                    Loading groups...
                </p>
            </div>
        );
    }

    if (!auth.authenticated) return null;

    return (
        <main className={styles.page}>
            <div className={styles.header}>
                <h1 className={styles.pageTitle}>
                    Groups
                </h1>

                <div className={styles.headerActions}>

                    <button
                        className={styles.createButton}
                        onClick={() => router.push('/groups/create')}
                    >
                        <span
                            className={styles.buttonIcon}
                            aria-hidden="true"
                        >
                            +
                        </span>

                        Create Group
                    </button>

                    <button
                        className={styles.browseButton}
                        onClick={showAllGroups}
                    >
                        All Groups
                    </button>

                    <button
                        className={styles.browseButton}
                        onClick={showMyGroups}
                    >
                        My Groups
                    </button>

                </div>
            </div>

            <section className={styles.groupsSection}>
                <div className={styles.container}>

                    {groups.length === 0 ? (
                        <p className={styles.emptyState}>
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
        </main>
    );
}