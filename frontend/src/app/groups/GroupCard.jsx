import styles from '@/app/styles/groups.module.css';
import { useRequireAuth } from '@/hooks/useRequireAuth';

export default function GroupCard({ group, joining, onJoin }) {
  const { auth, loading } = useRequireAuth();
  
  // Use creator_id (or group.admin if that's how your database stores it)
  // Optional chaining prevents errors while auth is loading
  const isCreator = auth.id === ( group.admin);

  return (
    <div className={styles.groupCard}>
      <h2 className={styles.groupTitle}>{group.titel}</h2>
      {group.description && (
        <p className={styles.groupDescription}>{group.description}</p>
      )}

      {/* Hide the button section if the user is the creator */}
      {!isCreator && (
        <div className={styles.groupActions}>
          <button
            className={`${styles.joinButton} ${
              joining ? styles.joinButtonSent : ''
            }`}
            onClick={() => onJoin(group.id)}
            disabled={joining}
          >
            {joining ? 'Requesting...' : 'Request to join'}
          </button>
        </div>
      )}
    </div>
  );
}