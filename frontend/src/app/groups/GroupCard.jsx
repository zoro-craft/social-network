import { useRouter } from "next/navigation";
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { useToast } from '@/app/context/ToastContext';
import { apiRequest } from "@/utils/api";

export default function GroupCard({ group }) {
    const { auth, loading } = useRequireAuth();
    const router = useRouter();
    const showToast = useToast();

    if (loading) return null;

    const isMember = group.is_member;
    const join = async () => {
        if (!auth.authenticated) return
        try {
            await apiRequest(`/groups/${group.id}/join`, 'POST')
        } catch (error) {
            showToast('Failed to request to join group');
            console.log('faild to request to join to groupe', error)
        }
    }
    const enterGroup = () => {
        router.push(`/groups/${group.id}`);
    };

    return (
        <div className="glass-group-card">
            <img
                className="glass-group-icon"
                src="/group-avatar.jpeg"
                alt={`${group.titel} avatar`}
            />

            <div className="glass-group-info">
                <h3 className="glass-group-title">
                    {group.titel}
                </h3>

                <p className="glass-group-description">
                    {group.description}
                </p>
            </div>

            {isMember ? (
                <button
                    className="glass-group-action-btn"
                    onClick={enterGroup}
                >
                    Enter Group
                </button>
            ) : (
                <button
                    className="glass-group-action-btn"
                    onClick={join}
                >
                    Request to join
                </button>
            )}

        </div>
    );
}