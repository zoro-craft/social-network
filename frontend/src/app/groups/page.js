'use client';

import { apiRequest } from "@/utils/api";
import { useState, useEffect } from "react";
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { useRouter } from "next/navigation";
import GroupCard from '@/components/GroupCard';
import styles from '../styles/groups.module.css';

export default function Groupspage() {
    const { auth, loading } = useRequireAuth();
    const [groups, setGroups] = useState([]);
    const [groupsLoading, setGroupsLoading] = useState(true);
    const router = useRouter();

    useEffect(() => {
        if (!auth.authenticated) return;
        const fetchGroups = async () => {
            try {
                const res = await apiRequest("/groups");
                setGroups(res || []);
            } catch {
                setGroups([]);
            } finally {
                setGroupsLoading(false);
            }
        };
        fetchGroups();
    }, [auth.authenticated]);

    if (loading || groupsLoading) {
        return (
            <div className="glass-feed-page">
                <p className={styles.mutedText}>Loading groups...</p>
            </div>
        );
    }

    if (!auth.authenticated) return null;

    return (
        <div className="glass-feed-page">
            <div className={styles.header}>
                <h1 className={styles.title}>Groups</h1>
                <button
                    className={`glass-button ${styles.createButton}`}
                    onClick={() => router.push('/groups/create')}
                >
                    + Create Group
                </button>
            </div>

            {groups.length === 0 && (
                <p className={styles.emptyState}>
                    No groups yet. Create the first one!
                </p>
            )}

            {groups.map((group) => (
                <GroupCard key={group.id} group={group} />
            ))}
        </div>
    );
}
