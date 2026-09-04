import Image from 'next/image';
import { API_BASE } from '@/utils/api';

export default function Avatar({ src, name, size = 28 }) {
  const imgSrc = src ? `${API_BASE}${src}` : '/default-avatar.png';

  return (
    <Image
      src={imgSrc}
      alt={name || ''}
      width={size}
      height={size}
      unoptimized
      style={{ borderRadius: '50%', objectFit: 'cover' }}
    />
  );
}
