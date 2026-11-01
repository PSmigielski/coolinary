import * as SecureStore from "expo-secure-store";

/*
 * Encrypted on-device storage (iOS Keychain / Android Keystore).
 *
 * PRIVACY RULE: health data (age, height, weight, goals) must only ever be
 * stored here, on the device. It is never sent to the server and never put in
 * AsyncStorage, the query cache or logs. The only value derived from it that
 * leaves the device is the daily kcal target.
 *
 * Keys may contain only letters, digits, ".", "-" and "_". Keep values small
 * (some platforms cap them at ~2 KB).
 */

export type SecureKey = "health.profile";

export async function readSecure<T>(key: SecureKey): Promise<T | null> {
  const raw = await SecureStore.getItemAsync(key);
  return raw === null ? null : (JSON.parse(raw) as T);
}

export async function writeSecure<T>(key: SecureKey, value: T): Promise<void> {
  await SecureStore.setItemAsync(key, JSON.stringify(value));
}

export async function removeSecure(key: SecureKey): Promise<void> {
  await SecureStore.deleteItemAsync(key);
}
