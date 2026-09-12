import { useRouter } from "next/navigation";
import styles from '@/app/styles/groups.module.css';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { apiRequest } from "@/utils/api";

export default function GroupCard({ group }) {
    const { auth, loading } = useRequireAuth();
    const router = useRouter();

    if (loading) return null;

    const isMember = group.is_member;
    const join = async () => {
        if (!auth.authenticated) return
        try {
            await apiRequest(`/groups/${group.id}/join`, 'POST')
        } catch (error) {
            console.log('faild to request to join to groupe', error)
        }
    }
    const enterGroup = () => {
        router.push(`/groups/${group.id}`);
    };

    return (
        <div className={styles.card}>
            <img
                className={styles.groupIcon}
                src="/group-avatar.jpeg"
                alt={`${group.titel} avatar`}
            />

            <div className={styles.groupInfo}>
                <h3 className={styles.title}>
                    {group.titel}
                </h3>

                <p className={styles.description}>
                    {group.description}
                </p>
            </div>


            {isMember ? (
                <button
                    className={styles.actionButton}
                    onClick={enterGroup}
                >
                    Enter Group
                </button>
            ) : (
                <button
                    className={styles.actionButton}
                    onClick={join}
                >
                    Request to join
                </button>
            )}

        </div>
    );
}