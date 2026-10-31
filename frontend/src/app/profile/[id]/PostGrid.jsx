import Image from "next/image";
import Link from "next/link";
import { API_BASE } from "@/utils/api";
import { isVideo } from "@/utils/media";

function PostGrid({ posts, isPrivateHidden }) {
  if (isPrivateHidden) return <p className="empty-state">🔒 Private account.</p>;
  if (!posts?.length) return <p className="empty-state">No posts yet.</p>;

  return (
    <div className="posts-grid">
      {posts.map((p) => (
        <Link key={p.id} href={`/posts/comment/${p.id}`} className="post-tile">
          {p.image && isVideo(p.image) ? (
            <video className="post-tile-img" src={`${API_BASE}${p.image}`} muted playsInline preload="metadata" />
          ) : p.image ? (
            <Image src={`${API_BASE}${p.image}`} className="post-tile-img" alt="" width={300} height={300} unoptimized />
          ) : (
            <p className="placeholder-text">{p.title}</p>
          )}
        </Link>
      ))}
    </div>
  );
}

export default PostGrid;
