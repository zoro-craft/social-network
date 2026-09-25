"use client";

import { useFollow } from "./useFollow";
import { useProfileModals } from "./ProfileModals";

const FALLBACK = { text: "Follow", className: "btn-follow" };

function FollowButton({ id, relation, isPrivate }) {
  const { openEditProfile, reload } = useProfileModals();
  const { toggleFollow, pending } = useFollow(id, reload);

  const config = {
    self: { text: "Edit Profile", className: "btn-edit" },
    accepted: { text: "Unfollow", className: "btn-unfollow" },
    pending: { text: isPrivate ? "Cancel Request" : "Following", className: "btn-pending" },
    none: FALLBACK,
  }[relation] || FALLBACK;

  const onToggle = relation === "self" ? openEditProfile : toggleFollow;

  return (
    <button onClick={onToggle} disabled={pending} className={`action-btn ${config.className}`}>
      {config.text}
    </button>
  );
}

export default FollowButton;
