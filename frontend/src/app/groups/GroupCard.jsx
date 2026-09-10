import { useRouter } from "next/navigation";
import styles from '@/app/styles/groups.module.css';
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
        <div className={styles.card}>

            <div className={styles.groupIcon} aria-hidden="true">
                <span className={styles.personOne} />
                <span className={styles.personTwo} />
            </div>

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
                >
                    Request to join
                </button>
            )}

        </div>
    );
}