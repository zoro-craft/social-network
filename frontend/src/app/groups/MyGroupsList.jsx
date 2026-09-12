'use client';

import { apiRequest } from "@/utils/api";
import { useState, useEffect } from "react";
import GroupCard from '@/app/groups/GroupCard';

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
        return <p className="glass-muted-text">Loading your groups...</p>;
    }

    if (error) {
        return <p className="glass-muted-text">{error}</p>;
    }

    if (groups.length === 0) {
        return <p className="glass-empty-state">You not  joined any groups yet.</p>;
    }

    return (
        <div className="glass-groups-grid">
            {groups.map((group) => (
                <GroupCard key={group.id} group={group} />
            ))}
        </div>
    );
}