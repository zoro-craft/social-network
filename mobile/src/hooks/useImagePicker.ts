import { Platform } from 'react-native';
import * as ImagePicker from 'expo-image-picker';
import { File } from 'expo-file-system';

export type PickedImage = ImagePicker.ImagePickerAsset;

export async function pickImage(): Promise<
  { ok: true; asset: PickedImage } | { ok: false; error?: string }
> {
  const permission = await ImagePicker.requestMediaLibraryPermissionsAsync();
  if (!permission.granted) {
    return { ok: false, error: 'Photo library permission is required.' };
  }
  const result = await ImagePicker.launchImageLibraryAsync({
    mediaTypes: ['images'],
    allowsEditing: false,
    quality: 0.8,
  });
  if (!result.canceled && result.assets && result.assets.length > 0) {
    return { ok: true, asset: result.assets[0] };
  }
  return { ok: false };
}

export async function pickVideo(): Promise<
  { ok: true; asset: PickedImage } | { ok: false; error?: string }
> {
  const permission = await ImagePicker.requestMediaLibraryPermissionsAsync();
  if (!permission.granted) {
    return { ok: false, error: 'Video library permission is required.' };
  }
  const result = await ImagePicker.launchImageLibraryAsync({
    mediaTypes: ['videos'],
    allowsEditing: false,
    quality: 0.8,
    // Reels are short by nature; the backend caps the file size, not the length.
    videoMaxDuration: 90,
  });
  if (!result.canceled && result.assets && result.assets.length > 0) {
    return { ok: true, asset: result.assets[0] };
  }
  return { ok: false };
}

export function assetFileName(asset: PickedImage): string {
  return asset.fileName || asset.uri.split('/').pop() || `${Date.now()}.jpg`;
}

// Expo's fetch serialises FormData itself and only accepts a string, a real
// Blob, or an object exposing bytes(). The old React Native { uri, name, type }
// part shape matches none of those and fails with
// "Unsupported FormDataPart implementation", so wrap the picked asset in a
// FileSystem File, which provides bytes() and keeps the real filename.
//
// On web the browser FormData is used instead, where a Blob is what append
// expects; the picked asset is already a browser File there.
export function appendAsset(
  form: FormData,
  field: string,
  asset: PickedImage,
  fallbackMime = 'image/jpeg',
) {
  if (Platform.OS === 'web') {
    form.append(field, asset as unknown as Blob);
    return;
  }

  const fileName = assetFileName(asset);
  const file = new File(asset.uri);

  // Android's picker returns content:// URIs whose basename has no extension
  // (e.g. .../1000000123), and File exposes neither `name` nor `type` as
  // writable properties. Both matter to the multipart part: the backend
  // validates by filename extension, and the part's content-type header is read
  // from `type`. Hence defineProperty rather than plain assignment.
  const define = (key: string, value: string) => {
    try {
      Object.defineProperty(file, key, { value, configurable: true, enumerable: true });
    } catch {
      // A non-configurable property falls back to the URI basename.
    }
  };

  define('name', fileName);
  define('type', asset.mimeType || fallbackMime);

  form.append(field, file, fileName);
}
