'use client';

import { apiRequest } from "@/utils/api";
import Avatar from "@/components/Avatar";
import { useState, useEffect } from "react";
import { useParams } from "next/navigation";
import { useRequireAuth } from '@/hooks/useRequireAuth';

export default function Userlistpage() {
    const { auth, loading } = useRequireAuth();
    const { id } = useParams();
    const [users, setUsers] = useState([]);
    
    const [invitedUsers, setInvitedUsers] = useState(new Set());

    useEffect(() => {
        if (!auth.authenticated) return;
        const Fetchusers = async () => {
            try {
                const res = await apiRequest("/users");
                setUsers(res || []);
            } catch {
                setUsers([]);
            }
        };
        Fetchusers();
    }, [auth.authenticated]);

    if (loading) {
        return (
            <div className="glass-feed-page">
                <p className="glass-muted-text">Loading ...</p>
            </div>
        );
    }

    const inviteUser = async (userID) => {
    try {
        const response = await apiRequest(
            `/groups/${id}/invitations`,
            "POST",
            {
                group_id: Number(id),
                invited_user_id: userID,
            }
        );

        setInvitedUsers((prev) => {
            const newSet = new Set(prev);
            newSet.add(userID);
            return newSet;
        });

    } catch (error) {
        console.error("Invitation error:", error);
    }
};

    return (
        <div className="glass-feed-page">
            <h1 className="glass-invite-title">Invite Members</h1>
            {users.length === 0 ? (
                <p className="glass-empty-state">No users found.</p>
            ) : (
                <div className="glass-invite-grid">
                    {users.map((user) => {
                        const isInvited = invitedUsers.has(user.id);

                        return (
                            <div key={user.id} className="glass-invite-card">
                                <div className="glass-invite-top">
                                    <Avatar src={user.avatar} name={user.full_name} size={48} />
                                    <h2 className="glass-invite-name">{user.full_name}</h2>
                                </div>
                                <button
                                    className="glass-invite-btn"
                                    onClick={() => inviteUser(user.id)}
                                    disabled={isInvited}
                                >
                                    {isInvited ? "... wait the response" : "Invite"}
                                </button>
                            </div>
                        );
                    })}
                </div>
            )}
        </div>
    );
}