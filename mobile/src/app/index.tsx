import { useState, type ReactNode } from "react";
import {
  ActivityIndicator,
  Pressable,
  RefreshControl,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import { $api, API_URL, type SystemStatus } from "@/api/client";
import { formatEnergy } from "@/lib/energy";
import { useSettings, type EnergyUnit } from "@/stores/settings";

const colors = {
  bg: "#f6f6f4",
  surface: "#ffffff",
  divider: "#ececea",
  text: "#141413",
  text2: "#6f6f6b",
  up: "#1e8e3e",
  down: "#c5221f",
};

export default function HomeScreen() {
  const status = $api.useQuery("get", "/v1/system/status");
  const [refreshing, setRefreshing] = useState(false);

  async function onRefresh() {
    setRefreshing(true);
    try {
      await status.refetch();
    } finally {
      setRefreshing(false);
    }
  }

  return (
    <SafeAreaView style={styles.screen} edges={["top", "left", "right"]}>
      <ScrollView
        contentContainerStyle={styles.content}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
      >
        <Text style={styles.title}>Coolinary</Text>
        <Text style={styles.subtitle}>From your fridge to this week&apos;s meals.</Text>

        <Section title="System status">
          {status.isPending ? (
            <View style={styles.row}>
              <ActivityIndicator />
              <Text style={styles.muted}>Checking {API_URL || "(no API URL set)"}…</Text>
            </View>
          ) : null}
          {status.isError ? (
            <ErrorNotice error={status.error as unknown} onRetry={() => status.refetch()} />
          ) : null}
          {status.data ? <StatusDetails data={status.data} /> : null}
        </Section>

        <Section title="Energy unit">
          <EnergyUnitSwitch />
        </Section>
      </ScrollView>
    </SafeAreaView>
  );
}

function Section({ title, children }: { title: string; children: ReactNode }) {
  return (
    <View style={styles.section}>
      <Text style={styles.sectionTitle}>{title}</Text>
      <View style={styles.card}>{children}</View>
    </View>
  );
}

function isUp(value: string) {
  return value === "SERVING" || value.toLowerCase() === "ok";
}

function StatusRow({ label, value }: { label: string; value: string }) {
  const up = isUp(value);
  return (
    <View style={styles.statusRow}>
      <View style={[styles.dot, { backgroundColor: up ? colors.up : colors.down }]} />
      <Text style={styles.label}>{label}</Text>
      <Text style={[styles.value, { color: up ? colors.up : colors.down }]}>{value}</Text>
    </View>
  );
}

function StatusDetails({ data }: { data: SystemStatus }) {
  const services = Object.entries(data.services).sort(([a], [b]) => a.localeCompare(b));
  return (
    <View>
      <Text style={[styles.overall, { color: data.ok ? colors.up : colors.down }]}>
        {data.ok ? "All systems OK" : "Not OK: something is down"}
      </Text>

      <Text style={styles.groupLabel}>Gateway</Text>
      <StatusRow label="gateway" value={data.gateway} />

      <Text style={styles.groupLabel}>Services</Text>
      {services.map(([name, value]) => (
        <StatusRow key={name} label={name} value={value} />
      ))}

      <Text style={styles.groupLabel}>Stores</Text>
      <StatusRow label="redis cache" value={data.stores.redisCache} />
      <StatusRow label="redis state" value={data.stores.redisState} />

      <Text style={styles.footnote}>
        Checked {new Date(data.checkedAt).toLocaleTimeString()} · pull down to refresh
      </Text>
    </View>
  );
}

function describeError(error: unknown): string {
  if (error instanceof Error) return error.message;
  if (error && typeof error === "object" && "error" in error) return String(error.error);
  return JSON.stringify(error);
}

function ErrorNotice({ error, onRetry }: { error: unknown; onRetry: () => void }) {
  return (
    <View style={styles.errorBox}>
      <Text style={styles.errorTitle}>Can&apos;t reach the gateway</Text>
      <Text style={styles.muted}>Tried: {API_URL || "(EXPO_PUBLIC_API_URL is not set)"}</Text>
      <Text style={styles.muted}>Error: {describeError(error)}</Text>
      <Text style={styles.muted}>
        Check the phone is on the same Wi-Fi and EXPO_PUBLIC_API_URL in mobile/.env uses this
        machine&apos;s LAN IP.
      </Text>
      <Pressable style={styles.button} onPress={onRetry} accessibilityRole="button">
        <Text style={styles.buttonText}>Retry</Text>
      </Pressable>
    </View>
  );
}

const UNITS: EnergyUnit[] = ["kcal", "kJ"];

function EnergyUnitSwitch() {
  const energyUnit = useSettings((s) => s.energyUnit);
  const setEnergyUnit = useSettings((s) => s.setEnergyUnit);

  return (
    <View>
      <View style={styles.segmented} accessibilityRole="radiogroup">
        {UNITS.map((unit) => {
          const selected = unit === energyUnit;
          return (
            <Pressable
              key={unit}
              onPress={() => setEnergyUnit(unit)}
              accessibilityRole="radio"
              accessibilityState={{ selected }}
              style={[styles.segment, selected && styles.segmentSelected]}
            >
              <Text style={[styles.segmentText, selected && styles.segmentTextSelected]}>
                {unit}
              </Text>
            </Pressable>
          );
        })}
      </View>
      <Text style={styles.footnote}>
        Example daily target: {formatEnergy(2000, energyUnit)}. Saved on this device.
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  screen: { flex: 1, backgroundColor: colors.bg },
  content: { padding: 20, paddingBottom: 48 },
  title: { fontSize: 34, fontWeight: "700", color: colors.text },
  subtitle: { fontSize: 15, color: colors.text2, marginTop: 4 },
  section: { marginTop: 28 },
  sectionTitle: {
    fontSize: 13,
    fontWeight: "600",
    color: colors.text2,
    textTransform: "uppercase",
    letterSpacing: 0.5,
    marginBottom: 8,
  },
  card: { backgroundColor: colors.surface, borderRadius: 14, padding: 16 },
  row: { flexDirection: "row", alignItems: "center", gap: 10 },
  muted: { fontSize: 14, color: colors.text2, marginTop: 4 },
  overall: { fontSize: 17, fontWeight: "600" },
  groupLabel: { fontSize: 13, fontWeight: "600", color: colors.text2, marginTop: 16 },
  statusRow: {
    flexDirection: "row",
    alignItems: "center",
    paddingVertical: 8,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderBottomColor: colors.divider,
  },
  dot: { width: 8, height: 8, borderRadius: 4, marginRight: 10 },
  label: { flex: 1, fontSize: 15, color: colors.text },
  value: { fontSize: 13, fontWeight: "500" },
  footnote: { fontSize: 12, color: colors.text2, marginTop: 12 },
  errorBox: { gap: 2 },
  errorTitle: { fontSize: 17, fontWeight: "600", color: colors.down },
  button: {
    alignSelf: "flex-start",
    marginTop: 12,
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderRadius: 8,
    backgroundColor: colors.text,
  },
  buttonText: { color: colors.surface, fontWeight: "600" },
  segmented: {
    flexDirection: "row",
    backgroundColor: colors.bg,
    borderRadius: 10,
    padding: 3,
  },
  segment: { flex: 1, alignItems: "center", paddingVertical: 8, borderRadius: 8 },
  segmentSelected: { backgroundColor: colors.text },
  segmentText: { fontSize: 15, fontWeight: "600", color: colors.text2 },
  segmentTextSelected: { color: colors.surface },
});
