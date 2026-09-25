"use client";

import { useProfileModals } from "./ProfileModals";

function StatItem({ type, label, value }) {
  const { openModal } = useProfileModals();

  return (
    <div className="glass-card stat-item" onClick={() => openModal(type)}>
      <p>{label}</p>
      <p>{value}</p>
    </div>
  );
}

export default StatItem;
