// isVideo reports whether a stored media path points at a playable video.
export function isVideo(src) {
  if (!src) return false;
  return /\.(mp4|webm|ogg)$/i.test(src.split('?')[0]);
}
