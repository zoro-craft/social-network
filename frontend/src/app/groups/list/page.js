"use client";

import { useState, useEffect } from "react";
import { useRequireAuth } from "@/hooks/useRequireAuth";
import { fetchGroups } from "@/app/groups/list/fetchgroups";
import styles from '../../styles/groups.module.css';

export default function GroupListPage() {
    const [groups, setGroups] = useState([]);

    const { auth, loading } = useRequireAuth();

    useEffect(() => {
        if (loading || !auth.authenticated) {
            return;
        }

        async function getGroups() {
            try {
                const data = await fetchGroups();
                setGroups(data);
            } catch (error) {
                console.error("Error fetching groups:", error);
            }
        }

        getGroups();
    }, [loading, auth.authenticated]);

    if (loading) {
        return <div className={styles.loading}>Loading...</div>;
    }

    if (!auth.authenticated) {
        return null;
    }

    return (
        <main className={styles.page}>
            <div className={styles.container}>
                <div className={styles.header}>
                    <h1 className={styles.title}>Groups</h1>
                    <p className={styles.subtitle}>
                        Discover and join groups that interest you.
                    </p>
                </div>

                {groups.length === 0 ? (
                    <div className={styles.empty}>No groups available yet.</div>
                ) : (
                    <div className={styles.grid}>
                        {groups.map((group) => (
                            <div className={styles.groupCard} key={group.id}>
                                <h2 className={styles.groupTitle}>
                                    {group.title}
                                </h2>
                                {group.description && (
                                    <p className={styles.groupDescription}>
                                        {group.description}
                                    </p>
                                )}

                                <div className={styles.groupActions}>
                                    <button className={styles.groupButton}>
                                        View Group
                                    </button>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </main>
    );
}
