'use client';

import { apiRequest } from "@/utils/api";
import Avatar from "@/components/Avatar";
import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { useToast } from '@/app/context/ToastContext';

export default function Userlistpage() {
    const { auth, loading } = useRequireAuth();
    const { id } = useParams();
    const [users, setUsers] = useState([]);
    const showToast = useToast();
    const router = useRouter();


    useEffect(() => {
        if (!auth.authenticated || !id) return;
        const Fetchusers = async () => {
            try {
                const res = await apiRequest(`/groups/${id}/users`);
                setUsers(res || []);
            } catch {
                showToast("Unable to load users.")
                 router.push("/groups")
            }
        };
        Fetchusers();
    }, [auth.authenticated, id]);

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
            if (response) {
                showToast('Invitation sent successfully' , false);
            }  

        } catch (error) {
            showToast('Failed to send invitation');
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

                        return (
                            <div key={user.id} className="glass-invite-card">
                                <div className="glass-invite-top">
                                    <Avatar src={user.avatar} name={user.full_name} size={48} />
                                    <h2 className="glass-invite-name">{user.full_name}</h2>
                                </div>
                                <button
                                    className="glass-invite-btn"
                                    onClick={() => inviteUser(user.id)}
                                >
                                    Invite
                                </button>
                            </div>
                        );
                    })}
                </div>
            )}
        </div>
    );
}