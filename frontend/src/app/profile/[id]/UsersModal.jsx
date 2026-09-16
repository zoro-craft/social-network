import Link from "next/link";
import Avatar from "@/components/Avatar";

function UsersModal({ type, users, onClose }) {
  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content users-modal" onClick={(e) => e.stopPropagation()}>
        <h2 style={{ marginTop: 0, textTransform: "capitalize" }}>{type}</h2>
        <div className="users-modal-list">
          {users.length ? users.map((u) => (
            <Link key={u.id} href={`/profile/${u.id}`} className="user-handle modal-user-row">
              <Avatar src={u.avatar} name={u.first_name} size={40} />
              <div>
                <div className="modal-user-name">{u.first_name} {u.last_name}</div>
                <div className="modal-user-handle">@{u.nickname || u.first_name}</div>
              </div>
            </Link>
          )) : <p className="users-modal-empty">No {type} found.</p>}
        </div>
      </div>
    </div>
  );
}

export default UsersModal;
