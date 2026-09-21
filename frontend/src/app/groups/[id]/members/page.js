'use client';

import { apiRequest } from "@/utils/api";
import Avatar from "@/components/Avatar";
import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { useToast } from '@/app/context/ToastContext';

export default function GroupMembersPage() {
    const { auth, loading } = useRequireAuth();
    const { id } = useParams();
    const [members, setMembers] = useState([]);
    const showToast = useToast();
    const router = useRouter();

    useEffect(() => {
        if (!auth.authenticated || !id) return;
        const fetchMembers = async () => {
            try {
                const res = await apiRequest(`/groups/${id}/members`);
                setMembers(res || []);
            } catch {
                showToast("Unable to load members.");
                router.push("/groups");
            }
        };
        fetchMembers();
    }, [auth.authenticated, id]);

    if (loading || !auth.authenticated) {
        return (
            <div className="glass-feed-page">
                <p className="glass-muted-text">Loading ...</p>
            </div>
        );
    }

    return (
        <div className="glass-feed-page">
            <h1 className="glass-member-title">Group Members</h1>
            {members.length === 0 ? (
                <p className="glass-empty-state">No members found.</p>
            ) : (
                <div className="glass-member-grid">
                    {members.map((member) => (
                        <div key={member.id} className="glass-member-card">
                            <div className="glass-member-top">
                                <Avatar src={member.avatar} name={member.full_name} size={48} />
                                <h2 className="glass-member-name">{member.full_name}</h2>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}