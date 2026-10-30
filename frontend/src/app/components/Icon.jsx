// One place for every UI glyph. The set is stroke-based, `currentColor` driven
// and sized in `em`, so a control's own font-size keeps the icon in scale and
// every responsive font tweak below keeps working without touching the markup.

const THUMB = (
  <>
    <path d="M9.6 11.4V20H5.7a1.7 1.7 0 0 1-1.7-1.7v-5.2a1.7 1.7 0 0 1 1.7-1.7z" />
    <path d="M9.6 11.4 12.6 4.7a2.2 2.2 0 0 1 3.5 1.6V9.6h4.1a1.8 1.8 0 0 1 1.8 2.1l-1.2 6.4a2.1 2.1 0 0 1-1.7 1.9H9.6z" />
  </>
);

const SPEAKER = (
  <path d="M11 5.2 6.6 9H3.6a1 1 0 0 0-1 1v4a1 1 0 0 0 1 1h3L11 18.8z" />
);

const ICONS = {
  home: (
    <>
      <path d="M3 10.5 12 3.6l9 6.9" />
      <path d="M5.6 9.2v10.2a1.4 1.4 0 0 0 1.4 1.4h2.6v-6h4.8v6h2.6a1.4 1.4 0 0 0 1.4-1.4V9.2" />
    </>
  ),

  film: (
    <>
      <rect x="2.8" y="4.8" width="18.4" height="14.4" rx="2.4" />
      <path d="M7.4 4.8v14.4M16.6 4.8v14.4M2.8 12h18.4M2.8 8.4h4.6M16.6 8.4h4.6M2.8 15.6h4.6M16.6 15.6h4.6" />
    </>
  ),

  users: (
    <>
      <path d="M15.5 20.5v-1.7a3.6 3.6 0 0 0-3.6-3.6H6.6a3.6 3.6 0 0 0-3.6 3.6v1.7" />
      <circle cx="9.3" cy="7.6" r="3.6" />
      <path d="M21 20.5v-1.7a3.6 3.6 0 0 0-2.7-3.5M15.8 4.2a3.6 3.6 0 0 1 0 7" />
    </>
  ),

  user: (
    <>
      <circle cx="12" cy="8.4" r="3.9" />
      <path d="M4.8 20.5v-1.4a5.3 5.3 0 0 1 5.3-5.3h3.8a5.3 5.3 0 0 1 5.3 5.3v1.4" />
    </>
  ),

  search: (
    <>
      <circle cx="10.8" cy="10.8" r="6.8" />
      <path d="M20.5 20.5 15.9 15.9" />
    </>
  ),

  bell: (
    <>
      <path d="M6.2 8.6a5.8 5.8 0 1 1 11.6 0c0 6.6 2.6 8.4 2.6 8.4H3.6s2.6-1.8 2.6-8.4Z" />
      <path d="M10.2 20.4a2.1 2.1 0 0 0 3.6 0" />
    </>
  ),

  chat: <path d="M7.9 20A9 9 0 1 0 4 16.1L2 22Z" />,

  message: <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />,

  thumbUp: THUMB,
  thumbDown: (
    <g transform="rotate(180 12 12)">
      {THUMB}
    </g>
  ),

  eye: (
    <>
      <path d="M2 12s3.7-6.6 10-6.6S22 12 22 12s-3.7 6.6-10 6.6S2 12 2 12z" />
      <circle cx="12" cy="12" r="2.9" />
    </>
  ),

  heart: (
    <path d="M12 20.5C12 20.5 3.5 15.3 3.5 9.2 3.5 6.3 5.8 4 8.7 4c1.6 0 2.9.8 3.3 1.7C12.4 4.8 13.7 4 15.3 4c2.9 0 5.2 2.3 5.2 5.2 0 6.1-8.5 11.3-8.5 11.3z" />
  ),

  share: (
    <>
      <path d="M7.5 16.5 16.5 7.5" />
      <path d="M9.5 7.5h7v7" />
    </>
  ),

  bookmark: (
    <path d="M6.6 3.6h10.8a1 1 0 0 1 1 1v16.5a.6.6 0 0 1-.9.5L12 17.7l-5.5 3.9a.6.6 0 0 1-.9-.5V4.6a1 1 0 0 1 1-1z" />
  ),

  trash: (
    <>
      <path d="M4 6.8h16" />
      <path d="M9.2 6.8V5.1a1.3 1.3 0 0 1 1.3-1.3h3a1.3 1.3 0 0 1 1.3 1.3v1.7" />
      <path d="M6.4 6.8 7.3 19.9a1.4 1.4 0 0 0 1.4 1.3h6.6a1.4 1.4 0 0 0 1.4-1.3l.9-13.1" />
      <path d="M10.4 10.6v6.4M13.6 10.6v6.4" />
    </>
  ),

  close: <path d="M6.5 6.5 17.5 17.5M17.5 6.5 6.5 17.5" />,

  plus: <path d="M12 5.5v13M5.5 12h13" />,

  volume: (
    <>
      {SPEAKER}
      <path d="M15.2 9.4a3.6 3.6 0 0 1 0 5.2M17.9 6.8a7.4 7.4 0 0 1 0 10.4" />
    </>
  ),

  volumeOff: (
    <>
      {SPEAKER}
      <path d="M15.8 9.5 21 14.7M21 9.5l-5.2 5.2" />
    </>
  ),

  smile: (
    <>
      <circle cx="12" cy="12" r="8.6" />
      <path d="M8.6 14.2a4.4 4.4 0 0 0 6.8 0" />
      <path d="M9.2 9.7h.01M14.8 9.7h.01" />
    </>
  ),

  lock: (
    <>
      <path d="M7.4 10.4V7.6a4.6 4.6 0 0 1 9.2 0v2.8" />
      <path d="M4.8 10.4h14.4a1 1 0 0 1 1 1V19a1 1 0 0 1-1 1H4.8a1 1 0 0 1-1-1v-7.6a1 1 0 0 1 1-1z" />
    </>
  ),
};

export default function Icon({
  name,
  filled = false,
  size,
  strokeWidth = 1.8,
  className,
  style,
}) {
  const glyph = ICONS[name];

  if (!glyph) return null;

  return (
    <svg
      className={className}
      width={size ?? '1em'}
      height={size ?? '1em'}
      viewBox="0 0 24 24"
      fill={filled ? 'currentColor' : 'none'}
      stroke="currentColor"
      strokeWidth={strokeWidth}
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
      focusable="false"
      style={{ verticalAlign: 'middle', flexShrink: 0, ...style }}
    >
      {glyph}
    </svg>
  );
}
