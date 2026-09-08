import styles from '@/app/styles/groups.module.css';
import { useRequireAuth } from '@/hooks/useRequireAuth';

export default function GroupCard({ group, joining, onJoin }) {
  const { auth, loading } = useRequireAuth();

  // Use creator_id (or group.admin if that's how your database stores it)
  // Optional chaining prevents errors while auth is loading
  const isCreator = auth.id === (group.admin);

  return (
    <div className={styles.card}>
      <h2 className={styles.title}>{group.titel}</h2>

      <p className={styles.description}>{group.description}</p>

      {/* Hide the button if the user is the creator */}
      {!isCreator && (
        <button
          className={styles.Btn}
          onClick={() => onJoin(group.id)}
        >
        Request to join
        </button>
      )}
      {isCreator && (
        <button className={styles.Btn} 
        onClick={() => router.push(`/groups/${group.id}/invite`)} 
        > 
        Invite Member </button>
         )}
    </div>
  );
}