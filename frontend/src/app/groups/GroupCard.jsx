import { useRouter } from "next/navigation";
import { useRequireAuth } from '@/hooks/useRequireAuth';

export default function GroupCard({ group }) {
    const { auth, loading } = useRequireAuth();
    const router = useRouter();

    if (loading) return null;

    const isMember = group.is_member;

    const enterGroup = () => {
        router.push(`/groups/${group.id}`);
    };

    return (
        <div className="glass-group-card">

            <div className="glass-group-icon" aria-hidden="true">
                <span className="glass-person-one" />
                <span className="glass-person-two" />
            </div>

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
                >
                    Request to join
                </button>
            )}

        </div>
    );
}