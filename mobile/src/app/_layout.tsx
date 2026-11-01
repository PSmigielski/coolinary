import { PersistQueryClientProvider } from "@tanstack/react-query-persist-client";
import { Stack } from "expo-router";
import { StatusBar } from "expo-status-bar";

import { persistOptions, queryClient } from "@/lib/query";

export default function RootLayout() {
  return (
    <PersistQueryClientProvider client={queryClient} persistOptions={persistOptions}>
      <StatusBar style="dark" />
      <Stack screenOptions={{ headerShown: false }} />
    </PersistQueryClientProvider>
  );
}
