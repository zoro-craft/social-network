import { useRouter } from "next/navigation";
import { useState } from "react";
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { useToast } from '@/app/context/ToastContext';
import { API_BASE, apiRequest } from "@/utils/api";

export default function GroupCard({ group }) {
    const { auth, loading } = useRequireAuth();
    const router = useRouter();
    const showToast = useToast();
    const [requested, setRequested] = useState(Boolean(group.is_pending));
    const [busy, setBusy] = useState(false);

    if (loading) return null;

    const isMember = group.is_member;
    const image = group.image ? `${API_BASE}${group.image}` : '/group-avatar.jpeg';

    const join = async () => {
        if (!auth.authenticated || busy) return;

        setBusy(true);
        try {
            const res = await apiRequest(`/groups/${group.id}/join`, 'POST');
            setRequested(true);
            showToast(res.message || 'Join request sent', false);
        } catch (error) {
            showToast(error.message || 'Failed to request to join group');
        } finally {
            setBusy(false);
        }
    };

    const cancelRequest = async () => {
        if (!auth.authenticated || busy) return;

        setBusy(true);
        try {
            const res = await apiRequest(`/groups/${group.id}/join`, 'DELETE');
            setRequested(false);
            showToast(res.message || 'Join request cancelled', false);
        } catch (error) {
            showToast(error.message || 'Failed to cancel join request');
        } finally {
            setBusy(false);
        }
    };

    const enterGroup = () => {
        router.push(`/groups/${group.id}`);
    };

    return (
        <div className="glass-group-card">
            <img
                className="glass-group-icon"
                src={image}
                alt={`${group.titel} avatar`}
                onError={(event) => {
                    event.currentTarget.onerror = null;
                    event.currentTarget.src = '/group-avatar.jpeg';
                }}
            />

            <div className="glass-group-info">
                <h3 className="glass-group-title">
                    {group.titel}
                </h3>

                <p className="glass-group-description">
                    {group.description}
                </p>

                <div className="glass-group-meta">
                    {typeof group.member_count === 'number' && (
                        <span>
                            {group.member_count} member{group.member_count === 1 ? '' : 's'}
                        </span>
                    )}
                    {requested && <span className="glass-group-pending">Request pending</span>}
                </div>
            </div>

            {isMember ? (
                <button
                    className="glass-group-action-btn"
                    onClick={enterGroup}
                >
                    Enter Group
                </button>
            ) : requested ? (
                <button
                    className="glass-group-action-btn glass-group-action-btn--cancel"
                    onClick={cancelRequest}
                    disabled={busy}
                >
                    {busy ? 'Cancelling...' : 'Cancel request'}
                </button>
            ) : (
                <button
                    className="glass-group-action-btn"
                    onClick={join}
                    disabled={busy}
                >
                    {busy ? 'Sending...' : 'Request to join'}
                </button>
            )}

        </div>
    );
}
