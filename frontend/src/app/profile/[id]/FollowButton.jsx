const FALLBACK = { text: "Follow", className: "btn-follow" };

function FollowButton({ relation, isPrivate, onToggle, disabled }) {
  const config = {
    self: { text: "Edit Profile", className: "btn-edit" },
    accepted: { text: "Unfollow", className: "btn-unfollow" },
    pending: { text: isPrivate ? "Cancel Request" : "Following", className: "btn-pending" },
    none: FALLBACK,
  }[relation] || FALLBACK;

  return (
    <button onClick={onToggle} disabled={disabled} className={`action-btn ${config.className}`}>
      {config.text}
    </button>
  );
}

export default FollowButton;
