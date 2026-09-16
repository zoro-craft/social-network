const FALLBACK = { text: "Follow", className: "btn-follow" };

// Declarative label/style map keeps every relation state accounted for —
// unknown states fall back to "Follow" instead of rendering a dead button.
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
