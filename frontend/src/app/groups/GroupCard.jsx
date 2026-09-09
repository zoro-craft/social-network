import { useRouter } from "next/navigation";
import styles from '@/app/styles/groups.module.css';
import { useRequireAuth } from '@/hooks/useRequireAuth';

export default function GroupCard({ group, joining, onJoin }) {
  const { auth, loading } = useRequireAuth();
  const router = useRouter();

  // Use creator_id (or group.admin if that's how your database stores it)
  // Optional chaining prevents errors while auth is loading
  const isCreator = auth.id === (group.admin);

  return (
    <div className={styles.card}>
      <div className={styles.groupIcon} aria-hidden="true">
        <span className={styles.personOne} />
        <span className={styles.personTwo} />
      </div>

      <div className={styles.groupInfo}>
        <h3 className={styles.title}>{group.titel}</h3>
        <p className={styles.description}>{group.description}</p>
      </div>

      {/* Hide the button if the user is the creator */}
      {!isCreator && (
        <button
          className={styles.actionButton}
          onClick={() => onJoin(group.id)}
        >
        Request to join
        </button>
      )}
      {isCreator && (
        <button
        className={styles.actionButton}
        onClick={() => router.push(`/groups/${group.id}/invite`)}
        >
        Invite Member </button>
         )}
    </div>
  );
}
