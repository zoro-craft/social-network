import Image from 'next/image';
import { API_BASE } from '@/utils/api';

// `eager` marks an avatar that is above the fold (the sidebar profile card, for
// instance) so Next.js stops counting it as a lazy LCP candidate.
export default function Avatar({ src, name, size = 28, eager = false }) {
  const imgSrc = src ? `${API_BASE}${src}` : '/default-avatar.png';

  return (
    <Image
      src={imgSrc}
      alt={name || ''}
      width={size}
      height={size}
      loading={eager ? 'eager' : 'lazy'}
      unoptimized
      style={{ borderRadius: '50%', objectFit: 'cover' }}
    />
  );
}
