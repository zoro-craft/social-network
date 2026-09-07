import styles from '@/app/styles/groups.module.css';

export default function GroupCard({ group, joining, onJoin }) {
  return (
    <div className={styles.groupCard}>
      <h2 className={styles.groupTitle}>{group.titel}</h2>
      {group.description && (
        <p className={styles.groupDescription}>{group.description}</p>
      )}
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
    </div>
  );
}