"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import "@/app/styles/profile.css";

export default function ProfilePage() {
  const params = useParams();
  const profileId = params.id;

  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [isEditing, setIsEditing] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");

  const backendURL = "http://localhost:8080";

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const res = await fetch(`${backendURL}/profile/${profileId}`, {
          credentials: "include",
        });
        if (res.ok) {
          const data = await res.json();
          setProfile(data);
        }
      } catch (error) {
        console.error("Error fetching profile:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, [profileId]);

  const handleFollowToggle = async () => {
    try {
      const res = await fetch(`${backendURL}/follow/${profileId}`, {
        method: "POST",
        credentials: "include",
      });
      if (res.ok) {
        const data = await res.json();
        setProfile((prev) => ({ ...prev, relation: data.new_status }));
      }
    } catch (error) {
      console.error("Follow action failed:", error);
    }
  };

  const handleEditSubmit = async (e) => {
    e.preventDefault();
    setErrorMsg("");
    const formData = new FormData(e.target);
    try {
      const res = await fetch(`${backendURL}/profile/edit`, {
        method: "POST",
        credentials: "include",
        body: formData,
      });
      if (res.ok) {
        setIsEditing(false);
        window.location.reload();
      } else {
        if (res.status === 400 || res.status === 413) {
          setErrorMsg(
            "Image size is too large. Please upload a smaller image .",
          );
        } else {
          const errData = await res.json().catch(() => null);
          setErrorMsg(
            errData?.error || "Failed to update profile. Please try again.",
          );
        }
      }
    } catch (error) {
      console.error("Error updating profile:", error);
      setErrorMsg("Network error. Please check your connection.");
    }
  };

  if (loading) return <div className="profile-wrapper">Loading...</div>;
  if (!profile) return <div className="profile-wrapper">User not found</div>;

  const avatarSrc = profile.user.avatar
    ? `${backendURL}${profile.user.avatar}`
    : "/default-avatar.png";
  const dob = profile.user.date_of_birth
    ? profile.user.date_of_birth.split("T")[0]
    : "";

  const renderActionButton = () => {
    if (profile.relation === "self")
      return (
        <button
          onClick={() => setIsEditing(true)}
          className="action-btn btn-edit"
        >
          Edit Profile
        </button>
      );
    if (profile.relation === "accepted")
      return (
        <button
          onClick={handleFollowToggle}
          className="action-btn btn-unfollow"
        >
          Unfollow
        </button>
      );
    if (profile.relation === "pending")
      return (
        <button onClick={handleFollowToggle} className="action-btn btn-pending">
          Cancel Request
        </button>
      );
    return (
      <button onClick={handleFollowToggle} className="action-btn btn-follow">
        Follow
      </button>
    );
  };

  return (
    <div className="profile-wrapper">
      <div className="profile-layout">
        <div className="glass-card left-sidebar">
          <div className="avatar-container">
            <img src={avatarSrc} alt="Avatar" className="avatar-img" />
            <span className="online-dot"></span>
          </div>

          <h2 className="user-name">
            {profile.user.first_name} {profile.user.last_name}
          </h2>
          <p className="user-handle">
            @{profile.user.nickname || profile.user.first_name}
          </p>

          <span className="privacy-badge">
            {profile.user.is_private ? "Private Profile" : "Public Profile"}
          </span>

          {profile.user.about_me && (
            <p className="about-me">{profile.user.about_me}</p>
          )}
          {renderActionButton()}
        </div>

        <div className="right-area">
          <div className="stats-grid">
            <div className="glass-card stat-item">
              <p className="stat-label">Followers</p>
              <p className="stat-value">{profile.followers_count || 0}</p>
            </div>
            <div className="glass-card stat-item">
              <p className="stat-label">Following</p>
              <p className="stat-value">{profile.following_count || 0}</p>
            </div>
            <div className="glass-card stat-item">
              <p className="stat-label">Posts</p>
              <p className="stat-value">0</p>
            </div>
          </div>

          <div className="glass-card">
            <h3 style={{ marginTop: 0, marginBottom: "1rem" }}>
              User Activity
            </h3>
            {profile.user.is_private &&
            profile.relation !== "self" &&
            profile.relation !== "accepted" ? (
              <div className="empty-state">
                <p>🔒 This account is private.</p>
                <small>Follow to see their posts.</small>
              </div>
            ) : (
              <div className="empty-state">No posts yet.</div>
            )}
          </div>
        </div>
      </div>

      {isEditing && (
        <div className="modal-overlay">
          <div className="modal-content">
            <h2 style={{ marginTop: 0 }}>Edit Profile</h2>
            {errorMsg && (
              <div
                style={{
                  background: "rgba(239, 68, 68, 0.1)",
                  border: "1px solid #ef4444",
                  color: "#f87171",
                  padding: "0.8rem",
                  borderRadius: "0.5rem",
                  marginBottom: "1rem",
                  fontSize: "0.9rem",
                }}
              >
                ⚠️ {errorMsg}
              </div>
            )}

            <form onSubmit={handleEditSubmit} className="edit-form">
              <div className="form-group privacy-toggle">
                <label>Private Profile</label>
                <input
                  type="checkbox"
                  name="is_private"
                  defaultChecked={profile.user.is_private}
                />
              </div>

              <input
                type="text"
                name="first_name"
                className="form-input"
                defaultValue={profile.user.first_name}
                placeholder="First Name"
                required
                maxLength={30}
              />
              <input
                type="text"
                name="last_name"
                className="form-input"
                defaultValue={profile.user.last_name}
                placeholder="Last Name"
                required
                maxLength={30}
              />
              <input
                type="email"
                name="email"
                className="form-input"
                defaultValue={profile.user.email}
                placeholder="Email"
                required
                maxLength={100}
              />
              <input
                type="date"
                name="date_of_birth"
                className="form-input"
                defaultValue={dob}
                required
              />
              <input
                type="text"
                name="nickname"
                className="form-input"
                defaultValue={profile.user.nickname}
                placeholder="Nickname (Optional)"
                maxLength={30}
              />
              <textarea
                name="about_me"
                className="form-input"
                defaultValue={profile.user.about_me}
                placeholder="About Me (Optional)"
                rows="3"
                maxLength={200}
              ></textarea>

              <div className="form-group">
                <label className="stat-label">
                  Avatar / Image (JPEG, PNG, GIF)
                </label>
                <input
                  type="file"
                  name="avatar"
                  accept="image/jpeg, image/png, image/gif"
                  className="form-file"
                />
                {profile.user.avatar && (
                  <label
                    style={{
                      display: "flex",
                      alignItems: "center",
                      gap: "0.5rem",
                      marginTop: "0.5rem",
                      fontSize: "0.85rem",
                      color: "#f87171",
                    }}
                  >
                    <input type="checkbox" name="remove_avatar" />
                    Remove current avatar
                  </label>
                )}
              </div>

              <div className="modal-actions">
                <button
                  type="button"
                  onClick={() => setIsEditing(false)}
                  className="action-btn btn-unfollow"
                  style={{ marginTop: 0 }}
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="action-btn btn-follow"
                  style={{ marginTop: 0 }}
                >
                  Save Changes
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
