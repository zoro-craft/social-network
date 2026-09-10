'use client';

import { apiRequest } from "@/utils/api";
import { useState, useEffect } from "react";
import GroupCard from '@/app/groups/GroupCard';
import styles from '../styles/groups.module.css';

export default function MyGroupsList() {
    const [groups, setGroups] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchMyGroups = async () => {
            try {
                setLoading(true);
                const res = await apiRequest("/groups/my-groups");
                const myGroups = (res || []).map((group) => ({
                    ...group,
                    is_member: true,
                }));
                setGroups(myGroups || []);

            } catch (err) {
                setError("Failed to load your groups.");
                setGroups([]);
            } finally {
                setLoading(false);
            }
        };

        fetchMyGroups();
    }, []);

    if (loading) {
        return <p className={styles.mutedText}>Loading your groups...</p>;
    }

    if (error) {
        return <p className={styles.mutedText}>{error}</p>;
    }

    if (groups.length === 0) {
        return <p className={styles.emptyState}>You not  joined any groups yet.</p>;
    }

    return (
        <div className={styles.container}>
            {groups.map((group) => (
                <GroupCard key={group.id} group={group} />
            ))}
        </div>
    );
}