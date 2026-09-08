'use client';

import { apiRequest } from "@/utils/api";
import Avatar from "@/components/Avatar";
import { useState, useEffect } from "react";
import { useParams } from "next/navigation";
import { useRequireAuth } from '@/hooks/useRequireAuth';
import styles from '../../../styles/user.module.css';

export default function userlistpage() {
    const { auth, loading } = useRequireAuth();
    const { id } = useParams();
    const [users, setUsers] = useState([])
    useEffect(() => {
        if (!auth.authenticated) return
        const Fetchusers = async () => {
            try {
                const res = await apiRequest("/users")
                setUsers(res || [])
            } catch {
                setUsers([])
            }
        }
        Fetchusers();
    }, [auth.authenticated])

    if (loading) {
        return (
            <div className="glass-feed-page">
                <p className={styles.mutedText}>Loading ...</p>
            </div>
        );
    }

    const inviteUser = async (userID) => {
        try {
            await apiRequest(`/groups/${id}/invitations`, "POST", {
                group_id: Number(id),
                invited_user_id: userID,
            });
        } catch {
            // invitation failed
        }
    }

    return (
        <div className="glass-feed-page">
            <h1 className={styles.pageTitle}>Invite Members</h1>
            {users.length === 0 ? (
                <p className={styles.emptyState}>No users found.</p>
            ) : (
                <div className={styles.grid}>
                    {users.map((user) => (
                        <div key={user.id} className={styles.card}>
                            <div className={styles.cardTop}>
                                <Avatar src={user.avatar} name={user.full_name} size={48} />
                                <h2 className={styles.name}>{user.full_name}</h2>
                            </div>
                            <button
                                className={styles.inviteBtn}
                                onClick={() => inviteUser(user.id)}
                            >
                                Invite
                            </button>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
