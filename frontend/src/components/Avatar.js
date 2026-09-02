import Image from 'next/image';
import { API_BASE } from '@/utils/api';

export default function Avatar({ src, name, size = 28 }) {
  if (src) {
    return (
      <Image
        src={`${API_BASE}${src}`}
        alt=""
        width={size}
        height={size}
        unoptimized
        style={{ borderRadius: '50%', objectFit: 'cover' }}
      />
    );
  }

  return (
    <span
      style={{
        width: size,
        height: size,
        borderRadius: '50%',
        background: 'linear-gradient(135deg, #a855f7, #7e22ce)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        color: '#fff',
        fontWeight: 700,
        fontSize: size * 0.45,
      }}
    >
      {name?.charAt(0).toUpperCase()}
    </span>
  );
}
