import { useState } from "react";
import { apiRequest } from "@/utils/api";

function EditProfileModal({ user, onClose, onSaved }) {
  const [formError, setFormError] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setFormError(null);
    try {
      await apiRequest("/profile/edit", "POST", new FormData(e.target), true);
      onSaved();
    } catch (err) {
      setFormError(err.message);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <h2 style={{ marginTop: 0 }}>Edit Profile</h2>
        <form onSubmit={handleSubmit} className="edit-form">
          {formError && <p className="form-error">{formError}</p>}
          <input name="first_name" defaultValue={user.first_name} className="form-input" required placeholder="First Name" />
          <input name="last_name" defaultValue={user.last_name} className="form-input" required placeholder="Last Name" />
          <input name="email" type="email" defaultValue={user.email} className="form-input" required placeholder="Email" />
          <input name="date_of_birth" type="date" defaultValue={user.date_of_birth?.split("T")[0]} className="form-input" required />
          <input name="nickname" defaultValue={user.nickname || ""} className="form-input" maxLength={30} placeholder="Nickname" />
          <textarea name="about_me" defaultValue={user.about_me} className="form-input" placeholder="About me" />
          <label className="privacy-toggle">
            <input type="checkbox" name="is_private" defaultChecked={!!user.is_private} /> Private Profile
          </label>
          <input type="file" name="avatar" className="form-file" accept="image/*" />
          <button type="submit" className="action-btn btn-follow">Save</button>
        </form>
      </div>
    </div>
  );
}

export default EditProfileModal;
