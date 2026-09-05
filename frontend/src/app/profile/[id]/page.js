"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import "@/app/styles/profile.css";

const BACKEND_URL = "http://localhost:8080";

export default function ProfilePage() {
  const { id: profileId } = useParams();
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [isEditing, setIsEditing] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");

  useEffect(() => {
    fetch(`${BACKEND_URL}/profile/${profileId}`, { credentials: "include" })
      .then((res) => res.json())
      .then(setProfile)
      .catch((err) => console.error("Error fetching profile:", err))
      .finally(() => setLoading(false));
  }, [profileId]);

  const handleFollowToggle = async () => {
    try {
      const res = await fetch(`${BACKEND_URL}/follow/${profileId}`, { method: "POST", credentials: "include" });
      if (res.ok) {
        const data = await res.json();
        setProfile((prev) => ({ ...prev, relation: data.new_status }));
      }
    } catch (error) { console.error("Follow action failed:", error); }
  };

  const handleEditSubmit = async (e) => {
    e.preventDefault();
    setErrorMsg("");
    try {
      const res = await fetch(`${BACKEND_URL}/profile/edit`, {
        method: "POST",
        credentials: "include",
        body: new FormData(e.target),
      });
      if (res.ok) return window.location.reload();
      
      const errData = await res.json().catch(() => null);
      setErrorMsg([400, 413].includes(res.status) 
        ? "Image too large. Please upload a smaller one." 
        : errData?.error || "Failed to update profile.");
    } catch (error) {
      setErrorMsg("Network error. Please check your connection.");
    }
  };

  if (loading) return <div className="profile-wrapper">Loading...</div>;
  if (!profile) return <div className="profile-wrapper">User not found</div>;

  const { user, relation, followers_count, following_count } = profile;
  const avatarSrc = user.avatar ? `${BACKEND_URL}${user.avatar}` : "/default-avatar.png";
  const isPrivateHidden = user.is_private && relation !== "self" && relation !== "accepted";

  return (
    <div className="profile-wrapper">
      <div className="profile-layout">
        
        {/* Left Sidebar */}
        <div className="glass-card left-sidebar">
          <div className="avatar-container">
            <img src={avatarSrc} alt="Avatar" className="avatar-img" />
            <span className="online-dot"></span>
          </div>
          
          <h2 className="user-name">{user.first_name} {user.last_name}</h2>
          <p className="user-handle">@{user.nickname || user.first_name}</p>
          <span className="privacy-badge">{user.is_private ? "Private Profile" : "Public Profile"}</span>
          {user.about_me && <p className="about-me">{user.about_me}</p>}

          {/* Action Button */}
          {relation === "self" ? <button onClick={() => setIsEditing(true)} className="action-btn btn-edit">Edit Profile</button>
          : relation === "accepted" ? <button onClick={handleFollowToggle} className="action-btn btn-unfollow">Unfollow</button>
          : relation === "pending" ? <button onClick={handleFollowToggle} className="action-btn btn-pending">Cancel Request</button>
          : <button onClick={handleFollowToggle} className="action-btn btn-follow">Follow</button>}
        </div>

        {/* Right Area */}
        <div className="right-area">
          <div className="stats-grid">
            <div className="glass-card stat-item"><p className="stat-label">Followers</p><p className="stat-value">{followers_count || 0}</p></div>
            <div className="glass-card stat-item"><p className="stat-label">Following</p><p className="stat-value">{following_count || 0}</p></div>
            <div className="glass-card stat-item"><p className="stat-label">Posts</p><p className="stat-value">0</p></div>
          </div>

          <div className="glass-card">
            <h3 style={{ marginTop: 0, marginBottom: "1rem" }}>User Activity</h3>
            <div className="empty-state">
              {isPrivateHidden ? <p>🔒 Private account. Follow to see posts.</p> : "No posts yet."}
            </div>
          </div>
        </div>
      </div>

      {/* Edit Modal */}
      {isEditing && (
        <div className="modal-overlay">
          <div className="modal-content">
            <h2 style={{ marginTop: 0 }}>Edit Profile</h2>
            {errorMsg && <div className="error-alert">⚠️ {errorMsg}</div>}

            <form onSubmit={handleEditSubmit} className="edit-form">
              <div className="form-group privacy-toggle">
                <label>Private Profile</label>
                <input type="checkbox" name="is_private" defaultChecked={user.is_private} />
              </div>

              <input type="text" name="first_name" className="form-input" defaultValue={user.first_name} placeholder="First Name" required maxLength={30} />
              <input type="text" name="last_name" className="form-input" defaultValue={user.last_name} placeholder="Last Name" required maxLength={30} />
              <input type="email" name="email" className="form-input" defaultValue={user.email} placeholder="Email" required maxLength={100} />
              <input type="date" name="date_of_birth" className="form-input" defaultValue={user.date_of_birth?.split("T")[0] || ""} required />
              <input type="text" name="nickname" className="form-input" defaultValue={user.nickname} placeholder="Nickname (Optional)" maxLength={30} />
              <textarea name="about_me" className="form-input" defaultValue={user.about_me} placeholder="About Me (Optional)" rows="3" maxLength={200}></textarea>

              <div className="form-group">
                <label className="stat-label">Avatar / Image (JPEG, PNG, GIF)</label>
                <input type="file" name="avatar" accept="image/jpeg, image/png, image/gif" className="form-file" />
                {user.avatar && (
                  <label className="remove-avatar-label">
                    <input type="checkbox" name="remove_avatar" /> Remove current avatar
                  </label>
                )}
              </div>

              <div className="modal-actions">
                <button type="button" onClick={() => setIsEditing(false)} className="action-btn btn-unfollow" style={{ marginTop: 0 }}>Cancel</button>
                <button type="submit" className="action-btn btn-follow" style={{ marginTop: 0 }}>Save Changes</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}