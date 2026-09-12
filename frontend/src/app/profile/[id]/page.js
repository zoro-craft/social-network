"use client";
import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";

const BACKEND_URL = "http://localhost:8080";

export default function ProfilePage() {
  const { id } = useParams();
  const [profile, setProfile] = useState(null);
  const [isEditing, setIsEditing] = useState(false);
  const [modal, setModal] = useState({ type: null, users: [] });

  useEffect(() => {
    fetch(`${BACKEND_URL}/profile/${id}`, { credentials: "include" })
      .then((res) => res.json())
      .then(setProfile);
  }, [id]);

  const fetchApi = async (endpoint, method = "POST", body = null) => {
    const options = { method, credentials: "include" };
    if (body && method !== "GET") options.body = body; 
    
    const res = await fetch(`${BACKEND_URL}${endpoint}`, options);
    return res.ok ? res.json() : Promise.reject(await res.json());
  };

const handleFollow = () => {
    fetchApi(`/follow/${id}`).then((data) => {
      const newRelation = data.status || data.new_status || "none";
      let newCount = profile.followers_count || 0;
      
      if (profile.relation === "accepted" && newRelation !== "accepted") {
        newCount = Math.max(0, newCount - 1); 
      } else if (profile.relation !== "accepted" && newRelation === "accepted") {
        newCount += 1; 
      }

      setProfile({ 
        ...profile, 
        relation: newRelation, 
        followers_count: newCount 
      });
    }).catch(err => console.error("Error toggling follow:", err));
  };

  const handleEdit = (e) => {
    e.preventDefault();
    fetchApi("/profile/edit", "POST", new FormData(e.target))
      .then(() => window.location.reload())
      .catch((err) => alert(err.error || "Update failed"));
  };

 const openModal = (type) => {
    if (profile.user.is_private && profile.relation !== "self" && profile.relation !== "accepted") {
      return; 
    }
    
    fetchApi(`/profile/${id}/${type}`, "GET")
      .then((users) => {
        setModal({ type, users: users || [] });
      })
      .catch((err) => {
        console.error("Error fetching list:", err);
      });
  };

  if (!profile) return <div className="profile-wrapper">Loading...</div>;

  const { user, relation, followers_count, following_count, posts } = profile;
  const isPrivateHidden = user.is_private && relation !== "self" && relation !== "accepted";

  return (
    <div className="profile-wrapper">
      <div className="profile-layout">
        
        {/* Left Sidebar */}
        <div className="glass-card left-sidebar">
          <img src={user.avatar ? `${BACKEND_URL}${user.avatar}` : "/default-avatar.png"} className="avatar-img" alt="Avatar" />
          <h2 className="user-name">{user.first_name} {user.last_name}</h2>
          <p className="user-handle">@{user.nickname || user.first_name}</p>
          <span className="privacy-badge">{user.is_private ? "Private Profile" : "Public Profile"}</span>
          {!isPrivateHidden && user.date_of_birth && <p style={{ color: "#94a3b8", fontSize: "0.9rem", margin: "10px 0" }}>Born: {user.date_of_birth.split("T")[0]}</p>}
          {!isPrivateHidden && user.about_me && <p className="about-me">{user.about_me}</p>}

          {/* Action Button */}
          {relation === "self" ? <button onClick={() => setIsEditing(true)} className="action-btn btn-edit">Edit Profile</button>
          : relation === "accepted" ? <button onClick={handleFollow} className="action-btn btn-unfollow">Unfollow</button>
          : relation === "pending" ? <button onClick={handleFollow} className="action-btn btn-pending">Cancel Request</button>
          : <button onClick={handleFollow} className="action-btn btn-follow">Follow</button>}
        </div>

        {/* Right Area */}
        <div className="right-area">
          <div className="stats-grid">
            <div className="glass-card stat-item" onClick={() => openModal("followers")} style={{ cursor: isPrivateHidden ? "default" : "pointer" }}>
              <p className="stat-label">Followers</p> 
              <p className="stat-value">{followers_count || 0}</p>
            </div>
            <div className="glass-card stat-item" onClick={() => openModal("following")} style={{ cursor: isPrivateHidden ? "default" : "pointer" }}>
              <p className="stat-label">Following</p> 
              <p className="stat-value">{following_count || 0}</p>
            </div>
            <div className="glass-card stat-item">
              <p className="stat-label">Joined</p> 
              <p className="stat-value">{user.created_at ? user.created_at.split("-")[0] : "Recently"}</p>
            </div>
          </div>

          {/* Posts Grid */}
          <div className="glass-card">
            <h3 style={{ marginTop: 0 }}>Posts</h3>
            {isPrivateHidden ? <p className="empty-state">🔒 Private account.</p> : (
              <div className="posts-grid">
                {posts?.length ? posts.map((p) => (
                  <Link key={p.id} href={`/posts/comment/${p.id}`} className="post-tile">
                    {p.image ? (
                        <img src={`${BACKEND_URL}${p.image}`} className="post-tile-img" alt="" />
                    ) : (
                        <p className="placeholder-text" style={{padding: "10px", textAlign: "center"}}>{p.title}</p>
                    )}
                  </Link>
                )) : <p className="empty-state">No posts yet.</p>}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Edit Modal */}
      {isEditing && (
        <div className="modal-overlay" onClick={() => setIsEditing(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h2 style={{ marginTop: 0 }}>Edit Profile</h2>
            <form onSubmit={handleEdit} className="edit-form">
              <input name="first_name" defaultValue={user.first_name} className="form-input" required placeholder="First Name" />
              <input name="last_name" defaultValue={user.last_name} className="form-input" required placeholder="Last Name" />
              <input name="email" type="email" defaultValue={user.email} className="form-input" required placeholder="Email" />
              <input name="date_of_birth" type="date" defaultValue={user.date_of_birth?.split("T")[0]} className="form-input" required />
              <input name="nickname" defaultValue={user.nickname || ""} className="form-input" maxLength={30} placeholder="Nickname" />
              <textarea name="about_me" defaultValue={user.about_me} className="form-input" placeholder="About me" />
              <label className="privacy-toggle">
                <input type="checkbox" name="is_private" defaultChecked={!!user.is_private} />
                Private Profile
              </label>
              <input type="file" name="avatar" className="form-file" accept="image/*" />
              <button type="submit" className="action-btn btn-follow">Save</button>
            </form>
          </div>
        </div>
      )}

      {/* Users List Modal (Followers / Following) */}
      {modal.type && (
        <div className="modal-overlay" onClick={() => setModal({ type: null, users: [] })}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()} style={{ maxWidth: "400px" }}>
            <h2 style={{ marginTop: 0, textTransform: "capitalize" }}>{modal.type}</h2>
            <div style={{ display: "flex", flexDirection: "column", gap: "10px", marginTop: "15px", maxHeight: "300px", overflowY: "auto" }}>
              {modal.users?.length ? modal.users.map((u) => {
                const itemStyle = {textDecoration: "none", display: "flex", alignItems: "center", gap: "10px"};
                if (user.is_private && relation !== "self") {
                  return (
                    <div key={u.id} className="user-handle" style={itemStyle}>
                      <img src={u.avatar ? `${BACKEND_URL}${u.avatar}` : "/default-avatar.png"} style={{width: 40, height: 40, borderRadius: "50%"}} alt="" />
                      <div>
                        <div style={{ color: "white", fontWeight: "bold" }}>{u.first_name} {u.last_name}</div>
                        <div style={{ color: "#94a3b8", fontSize: "0.8rem" }}>@{u.nickname || u.first_name}</div>
                      </div>
                    </div>
                  );
                }
                return (
                  <Link href={`/profile/${u.id}`} key={u.id} className="user-handle" onClick={() => setModal({type: null})} style={itemStyle}>
                    <img src={u.avatar ? `${BACKEND_URL}${u.avatar}` : "/default-avatar.png"} style={{width: 40, height: 40, borderRadius: "50%"}} alt="" />
                    <div>
                      <div style={{ color: "white", fontWeight: "bold" }}>{u.first_name} {u.last_name}</div>
                      <div style={{ color: "#94a3b8", fontSize: "0.8rem" }}>@{u.nickname || u.first_name}</div>
                    </div>
                  </Link>
                );
              }) : <p style={{ color: "#94a3b8", textAlign: "center" }}>No {modal.type} found.</p>}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
