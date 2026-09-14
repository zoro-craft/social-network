"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import Avatar from "@/components/Avatar";
import { apiRequest } from "@/utils/api";
import { useToast } from "@/app/context/ToastContext";
import { useProfile } from "./useProfile";
import { useFollow } from "./useFollow";
import FollowButton from "./FollowButton";
import PostGrid from "./PostGrid";
import EditProfileModal from "./EditProfileModal";
import UsersModal from "./UsersModal";

export default function ProfilePage() {
  const { id } = useParams();
  const showToast = useToast();
  const { profile, error, reload } = useProfile(id);
  const { toggleFollow, pending: followPending } = useFollow(id, reload);
  const [modalType, setModalType] = useState(null);
  const [usersList, setUsersList] = useState([]);
  const [isEditing, setIsEditing] = useState(false);

  const openModal = async (type) => {
    if (!profile || (profile.user.is_private && profile.relation !== "self" && profile.relation !== "accepted")) return;
    setModalType(type);
    try {
      const data = await apiRequest(`/profile/${id}/${type}`);
      setUsersList(data || []);
    } catch (err) {
      setModalType(null);
      showToast(err.message);
    }
  };

  const closeModal = () => {
    setModalType(null);
    setUsersList([]);
  };

  if (error) {
    return (
      <div className="profile-wrapper">
        <div className="glass-card profile-error-card">
          <h2 style={{ marginTop: 0 }}>Profile unavailable</h2>
          <p>{error}</p>
          <Link href="/posts" className="action-btn btn-edit">Back to posts</Link>
        </div>
      </div>
    );
  }

  if (!profile) return <div className="profile-wrapper">Loading...</div>;

  const { user, relation, followers_count, following_count, posts } = profile;
  const isPrivateHidden = user.is_private && relation !== "self" && relation !== "accepted";

  return (
    <div className="profile-wrapper">
      <div className="profile-layout">
        <div className="glass-card left-sidebar">
          <Avatar src={user.avatar} name={user.first_name} size={112} />
          <h2 className="user-name">{user.first_name} {user.last_name}</h2>
          <p className="user-handle">@{user.nickname || user.first_name}</p>
          <span className="privacy-badge">{user.is_private ? "Private Profile" : "Public Profile"}</span>
          {!isPrivateHidden && user.date_of_birth && <p className="profile-dob">Born: {user.date_of_birth.split("T")[0]}</p>}
          {!isPrivateHidden && user.about_me && <p className="about-me">{user.about_me}</p>}

          <FollowButton relation={relation} isPrivate={user.is_private} onToggle={relation === "self" ? () => setIsEditing(true) : toggleFollow} disabled={followPending} />
        </div>

        <div className="right-area">
          <div className="stats-grid">
            {[
              { label: "Followers", value: followers_count || 0, type: "followers" },
              { label: "Following", value: following_count || 0, type: "following" },
              { label: "Joined", value: user.created_at ? user.created_at.split("-")[0] : "Recently" },
            ].map(({ label, value, type }) => (
              <div key={label} className="glass-card stat-item" onClick={type && !isPrivateHidden ? () => openModal(type) : undefined} style={{ cursor: type && !isPrivateHidden ? "pointer" : "default" }}>
                <p className="stat-label">{label}</p>
                <p className="stat-value">{value}</p>
              </div>
            ))}
          </div>

          <div className="glass-card">
            <h3 style={{ marginTop: 0 }}>Posts</h3>
            <PostGrid posts={posts} isPrivateHidden={isPrivateHidden} />
          </div>
        </div>
      </div>

      {isEditing && (
        <EditProfileModal user={user} onClose={() => setIsEditing(false)} onSaved={async () => { setIsEditing(false); await reload(); showToast("Profile updated", false); }} />
      )}

      {modalType && <UsersModal type={modalType} users={usersList} onClose={closeModal} />}
    </div>
  );
}
