import AsyncStorage from "@react-native-async-storage/async-storage";
import { createAsyncStoragePersister } from "@tanstack/query-async-storage-persister";
import { focusManager, QueryClient } from "@tanstack/react-query";
import type { PersistQueryClientProviderProps } from "@tanstack/react-query-persist-client";
import { AppState, Platform } from "react-native";

const DAY_MS = 24 * 60 * 60 * 1000;

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      staleTime: 30_000,
      // Must be >= persistOptions.maxAge, or restored entries are garbage-collected right away.
      gcTime: DAY_MS,
    },
  },
});

/** Server-state cache survives app reloads; stale data is shown, then refetched. */
export const persistOptions: PersistQueryClientProviderProps["persistOptions"] = {
  persister: createAsyncStoragePersister({
    storage: AsyncStorage,
    key: "coolinary.query-cache",
  }),
  maxAge: DAY_MS,
  // Bump when cached response shapes change incompatibly.
  buster: "v1",
};

// React Native has no window focus events: treat "app came to the foreground" as focus,
// so stale queries refetch when the user returns to the app.
if (Platform.OS !== "web") {
  focusManager.setEventListener((handleFocus) => {
    const subscription = AppState.addEventListener("change", (state) => {
      handleFocus(state === "active");
    });
    return () => subscription.remove();
  });
}
