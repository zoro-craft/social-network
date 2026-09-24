"use client";

import { useState } from "react";
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

  // Fetch data
  const { profile, error, reload } = useProfile(id);
  const { toggleFollow, pending } = useFollow(id, reload);

  // Modal states
  const [modalType, setModalType] = useState(null); // "followers" or "following"
  const [usersList, setUsersList] = useState([]);
  const [isEditing, setIsEditing] = useState(false);

  // Loading & Errors
  if (error) {
    return (
      <div className="profile-wrapper">
        <div className="glass-card">
          <h2>Profile unavailable</h2>
          <p>{error}</p>
          <Link href="/posts">Back to posts</Link>
        </div>
      </div>
    );
  }

  if (!profile) {
    return <div className="profile-wrapper">Loading...</div>;
  }

  // Extract data
  const { user, relation, followers_count, following_count, posts } = profile;
  
  // Check if profile is private 
  const isPrivateHidden = user.is_private && relation !== "self" && relation !== "accepted";

  // Open modal function
  const openModal = async (type) => {
    if (isPrivateHidden) return; // Stop if private
    
    setModalType(type); // Show modal
    try {
      const data = await apiRequest(`/profile/${id}/${type}`);
      setUsersList(data || []); // Save users
    } catch (err) {
      setModalType(null); // Close on error
      showToast(err.message);
    }
  };

  // Render UI
  return (
    <div className="profile-wrapper">
      <div className="profile-layout">
        
        {/* Left Sidebar: User Info */}
        <div className="glass-card left-sidebar">
          <Avatar src={user.avatar} name={user.first_name} size={112} />
          <h2>{user.first_name} {user.last_name}</h2>
          <p>@{user.nickname || user.first_name}</p>
          <p>{user.is_private ? "Private Profile" : "Public Profile"}</p>

          {/* Show if not private */}
          {!isPrivateHidden && (
            <>
              <p>Born: {user.date_of_birth?.split("T")[0]}</p>
              <p>{user.about_me}</p>
            </>
          )}

          <FollowButton 
            relation={relation} 
            isPrivate={user.is_private} 
            onToggle={relation === "self" ? () => setIsEditing(true) : toggleFollow} 
            disabled={pending} 
          />
        </div>

        {/* Right Area: Stats & Posts */}
        <div className="right-area">
          
          {/* Stats */}
          <div className="stats-grid">
            <div className="glass-card stat-item" onClick={() => openModal("followers")}>
              <p>Followers</p>
              <p>{followers_count || 0}</p>
            </div>

            <div className="glass-card stat-item" onClick={() => openModal("following")}>
              <p>Following</p>
              <p>{following_count || 0}</p>
            </div>

            <div className="glass-card stat-item">
              <p>Joined</p>
              <p>{user.created_at ? user.created_at.split("-")[0] : "Recently"}</p>
            </div>
          </div>

          {/* Posts */}
          <div className="glass-card">
            <h3>Posts</h3>
            <PostGrid posts={posts} isPrivateHidden={isPrivateHidden} />
          </div>

        </div>
      </div>

      {/* Modals */}
      
      {/* Edit Profile */}
      {isEditing && (
        <EditProfileModal 
          user={user} 
          onClose={() => setIsEditing(false)} 
          onSaved={async () => { 
            setIsEditing(false); 
            await reload(); 
            showToast("Profile updated", false); 
          }} 
        />
      )}

      {/* Users List */}
      {modalType && (
        <UsersModal 
          type={modalType} 
          users={usersList} 
          onClose={() => setModalType(null)} 
        />
      )}
    </div>
  );
}