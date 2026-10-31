import { API_BASE } from '@/utils/api';
import { isVideo } from '@/utils/media';

// PostMedia renders a post attachment as an image or a small video player.
export default function PostMedia({ src, className = 'glass-post-image' }) {
  if (!src) return null;

  const url = `${API_BASE}${src}`;
  if (isVideo(src)) {
    return <video className={`${className} glass-post-video`} src={url} controls playsInline preload="metadata" />;
  }

  // eslint-disable-next-line @next/next/no-img-element
  return <img className={className} src={url} alt="" />;
}
