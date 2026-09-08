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

      <div className={styles.descriptionBox}>
        {group.description ? (
          <p>{group.description}</p>
        ) : (
          <p className={styles.mutedText}>No description yet.</p>
        )}
      </div>

      {/* Hide the button if the user is the creator */}
      {!isCreator && (
        <button
          className={`${styles.joinBtn} ${
            joining ? styles.joinBtnSent : ''
          }`}
          onClick={() => onJoin(group.id)}
          disabled={joining}
        >
          {joining ? 'Requesting...' : 'Request to join'}
        </button>
      )}
    </div>
  );
}